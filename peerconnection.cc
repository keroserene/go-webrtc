/**
 * C wrapper around the C++ webrtc::PeerConnectionInterface and related, which
 * allows compatibility with CGO's requirements so that everything may
 * ultimately be exposed in Go.
 */
#include "peerconnection.h"
#include "datachannel.hpp"

#include <iostream>
#include <future>
#include <mutex>

#include "webrtc/api/test/fakeconstraints.h"
#include "webrtc/pc/test/fakeaudiocapturemodule.h"
#include "webrtc/api/jsepsessiondescription.h"
#include "webrtc/pc/webrtcsdp.h"

#define SUCCESS 0
#define FAILURE -1
#define TIMEOUT_SECS 3

#define CGO_DBG_ENABLED 0
#define CGO_DBG_MSG(os, msg) \
  (os) << endl << "[CGO] " << __func__ << "() - line " << __LINE__ << ": "\
       << msg << endl
#define CGO_DBG(msg) if (CGO_DBG_ENABLED) { CGO_DBG_MSG(cout, msg); }

using namespace std;
using namespace webrtc;

typedef rtc::scoped_refptr<webrtc::PeerConnectionInterface> PC;
typedef SessionDescriptionInterface* SDP;
typedef rtc::scoped_refptr<CGoDataChannelObserver> DCObserver;

// Peer acts as the glue between Go PeerConnection and the native
// webrtc::PeerConnectionInterface. However, it's not directly accessible
// through CGO, but indirectly through what's available in the more pure
// extern "C" header.
//
// The Go side may access this class through C.CGO_Peer.
class Peer
  : public PeerConnectionObserver,
    public CreateSessionDescriptionObserver {
 public:

  // Expected to be called before anything else happens for Peer.
  bool Initialize() {
    promiseSDP = promise<SDP>();

    // Due to the different threading model, in order for PeerConnectionFactory
    // to be able to post async messages without getting blocked, we need to use
    // external signalling and worker thread, accounted for in this class.
    signalling_thread_ = std::unique_ptr<rtc::Thread>(new rtc::Thread());
    worker_thread_ = std::unique_ptr<rtc::Thread>(new rtc::Thread());
    signalling_thread_->SetName("CGO Signalling", NULL);
    worker_thread_->SetName("CGO Worker", NULL);
    signalling_thread_->Start();  // Must start before being passed to
    worker_thread_->Start();      // PeerConnectionFactory.

    this->fake_audio_ = FakeAudioCaptureModule::Create();
    pc_factory = CreatePeerConnectionFactory(
      worker_thread_.get(),
      signalling_thread_.get(),
      this->fake_audio_, NULL, NULL);
    if (!pc_factory.get()) {
      CGO_DBG("Could not create PeerConnectionFactory");
      return false;
    }

    // Media constraints are hard-coded here and not exposed in Go, since
    // in this case, only DTLS/SCTP data channels are desired. If this ever
    // changes (eg. enabling the Media API) then this will need a Go interface.
    auto c = new FakeConstraints();
    c->AddOptional(MediaConstraintsInterface::kEnableDtlsSrtp, true);
    constraints = c;
    return true;
  }

  void resetPromise() {
    promiseSDP = promise<SDP>();
  }

  //
  // CreateSessionDescriptionObserver implementation
  //
  // These callbacks have been stubbed out using promises + futures, to be
  // blocking as far as Go is concerned, which allows the usage
  // of goroutines. This should be easier and more idiomatic for users.
  //
  void OnSuccess(SDP desc) {
    CGO_DBG("SDP successfully created.");
    promiseSDP.set_value(desc);
  }

  void OnFailure(const std::string& error) {
    CGO_DBG("SDP Failure: " + error);
    promiseSDP.set_value(NULL);
  }

  //
  // PeerConnectionObserver Implementation
  //
  void OnSignalingChange(PeerConnectionInterface::SignalingState state) {
    CGO_DBG("fired OnSignalingChange");
    cgoOnSignalingStateChange(goPeerConnection, state);
  }

  // This is required for the Media API.
  void OnAddStream(webrtc::MediaStreamInterface* stream) {
    CGO_DBG("unimplemented OnAddStream");
  }
  void OnRemoveStream(webrtc::MediaStreamInterface* stream) {
    CGO_DBG("unimplemented OnRemoveStream");
  }

  void OnRenegotiationNeeded() {
    CGO_DBG("fired OnRenegotiationNeeded");
    cgoOnNegotiationNeeded(goPeerConnection);
  }

  void OnIceCandidate(const IceCandidateInterface* ic) {
    std::string sdp;
    ic->ToString(&sdp);
    // Cast IceCandidate to Go-compatible C struct.
    CGO_IceCandidate cgoIC = {
      const_cast<char*>(ic->sdp_mid().c_str()),
      ic->sdp_mline_index(),
      const_cast<char*>(sdp.c_str())
    };
    cgoOnIceCandidate(goPeerConnection, cgoIC);
  }

  void OnIceConnectionChange(
      PeerConnectionInterface::IceConnectionState new_state) {
    if (PeerConnectionInterface::IceConnectionState::kIceConnectionFailed ==
        new_state) {
      cgoOnIceCandidateError(goPeerConnection);
      return;
    }
    cgoOnIceConnectionStateChange(goPeerConnection, new_state);
    // The ice connection state is sent to the go callback
    // which then translates it into RTCPeerConnectionState
    cgoOnConnectionStateChange(goPeerConnection, new_state);
  }

  void OnIceGatheringChange(
      PeerConnectionInterface::IceGatheringState new_state) {
    cgoOnIceGatheringStateChange(goPeerConnection, new_state);
  }

  void OnDataChannel(DataChannelInterface* channel) {
    DCObserver obs = new rtc::RefCountedObject<CGoDataChannelObserver>(channel);
    o_lock.lock();
    observers.push_back(obs);
    o_lock.unlock();
    auto o = obs.get();
    channel->RegisterObserver(o);
    cgoOnDataChannel(goPeerConnection, (void *)o);
  }

  void SetConfig(PeerConnectionInterface::RTCConfiguration *c) {
    if (config)
      delete config;
    config = c;
  }

  // Note that Configuration is where ICE servers are specified.
  PeerConnectionInterface::RTCConfiguration *config = NULL;
  const FakeConstraints* constraints = NULL;

  PC pc_;                  // Pointer to webrtc::PeerConnectionInterface.
  int goPeerConnection;    // Pointer to external Go PeerConnection struct,
                           // which is required to fire callbacks correctly.

  // Pass SDPs through promises instead of callbacks, to allow benefits as
  // described above. However, this means CreateOffer and CreateAnswer must
  // not be concurrent to themselves or each other (which isn't expected
  // anyways), due to the simplistic way futures are used here.
  promise<SDP> promiseSDP;

  rtc::scoped_refptr<PeerConnectionFactoryInterface> pc_factory;

  // Prevent deallocation of created DataChannels, since they are ref_ptr,
  // by keeping track of them in a vector.
  vector<DCObserver> observers;
  std::mutex o_lock;

 protected:
  ~Peer() {
    SetConfig(NULL);
    if (constraints)
      delete constraints;

    // NOTE: Clears these explicitly first since they use the threads
    observers.clear();
    pc_ = nullptr;
    pc_factory = nullptr;

    worker_thread_->Stop();
    signalling_thread_->Stop();

    worker_thread_ = nullptr;
    signalling_thread_ = nullptr;
  }

 private:
  std::unique_ptr<rtc::Thread> signalling_thread_;
  std::unique_ptr<rtc::Thread> worker_thread_;
  rtc::scoped_refptr<AudioDeviceModule> fake_audio_;
};  // class Peer

// Keep track of Peers in global scope to prevent deallocation, due to the
// required scoped_refptr from implementing the Observer interface.
vector<rtc::scoped_refptr<Peer>> localPeers;
std::mutex lp_lock;

class PeerSDPObserver : public SetSessionDescriptionObserver {
 public:
  static PeerSDPObserver* Create() {
    return new rtc::RefCountedObject<PeerSDPObserver>();
  }
  virtual void OnSuccess() {
    promiseSet.set_value(0);
  }
  virtual void OnFailure(const std::string& error) {
    CGO_DBG("SessionDescription: " + error);
    promiseSet.set_value(-1);
  }
  promise<int> promiseSet = promise<int>();

 protected:
  PeerSDPObserver() {}
  ~PeerSDPObserver() {}

};  // class PeerSDPObserver

//
// extern "C" Go-accessible functions:
//

// Create and return the Peer object, which provides initial native code
// glue for the PeerConnection constructor.
CGO_Peer CGO_InitializePeer(int goPc) {
  rtc::scoped_refptr<Peer> localPeer = new rtc::RefCountedObject<Peer>();
  localPeer->Initialize();
  lp_lock.lock();
  localPeers.push_back(localPeer);
  lp_lock.unlock();
  localPeer->goPeerConnection = goPc;
  return localPeer;
}

void CGO_DestroyPeer(CGO_Peer cgoPeer) {
  auto cPeer = (Peer*)cgoPeer;
  lp_lock.lock();
  localPeers.erase(
    std::remove_if(localPeers.begin(), localPeers.end(), [cPeer](rtc::scoped_refptr<Peer> localPeer){
      return localPeer.get() == cPeer;
    }),
    localPeers.end()
  );
  lp_lock.unlock();
}

// This helper converts RTCConfiguration struct from GO to C++.
PeerConnectionInterface::RTCConfiguration *castConfig_(
    CGO_Configuration *cgoConfig) {
  PeerConnectionInterface::RTCConfiguration* c =
      new PeerConnectionInterface::RTCConfiguration();

  // Pass in all IceServer structs for PeerConnectionInterface.
  vector<CGO_IceServer> servers( cgoConfig->iceServers,
      cgoConfig->iceServers + cgoConfig->numIceServers);
  for (auto s : servers) {
    // cgo only allows C arrays, but webrtc expects std::vectors
    vector<string> urls(s.urls, s.urls + s.numUrls);
    PeerConnectionInterface::IceServer is {};
    is.uri ="";  // TODO: Remove once webrtc deprecates the first uri field.
    is.urls = urls;
    is.username = s.username;
    is.password = s.credential;
    c->servers.push_back(is);
  }

  // Cast Go const "enums" to C++ Enums.
  c->type = (PeerConnectionInterface::IceTransportsType)
      cgoConfig->iceTransportPolicy;
  c->bundle_policy = (PeerConnectionInterface::
      BundlePolicy)cgoConfig->bundlePolicy;
  // TODO: [ED] extensions. Corresponding enum in configuration.go.
  // c->rtcp_mux_policy = (PeerConnectionInterface::
      // RtcpMuxPolicy)cgoConfig->RtcpMuxPolicy;
  return c;
}

// |Peer| method: create a native code PeerConnection object.
// Returns 0 on Success.
int CGO_CreatePeerConnection(CGO_Peer cgoPeer, CGO_Configuration *cgoConfig) {
  Peer *peer = (Peer*)cgoPeer;
  peer->SetConfig(castConfig_(cgoConfig));
  peer->pc_ = peer->pc_factory->CreatePeerConnection(
    *peer->config,
    peer->constraints,
    NULL,  // port allocator      (reasonable default already within)
    NULL,  // dtls identity store (reasonable default already within)
    peer   // "observer"
    );

  if (!peer->pc_.get()) {
    CGO_DBG("Could not create PeerConnection.");
    return FAILURE;
  }
  return SUCCESS;
}

bool SDPtimeout(future<SDP> *f, int seconds) {
  auto status = f->wait_for(chrono::seconds(TIMEOUT_SECS));
  return future_status::ready != status;
}

// PeerConnection::CreateOffer
// Blocks until libwebrtc succeeds in generating the SDP offer,
// @returns SDP (pointer), or NULL on timeeout.
CGO_sdpString CGO_CreateOffer(CGO_Peer cgoPeer) {
  Peer* peer = (Peer*)cgoPeer;
  auto r = peer->promiseSDP.get_future();
  peer->pc_->CreateOffer(peer, peer->constraints);
  if (SDPtimeout(&r, TIMEOUT_SECS)) {
    CGO_DBG("CreateOffer timed out after " + TIMEOUT_SECS);
    peer->resetPromise();
    return NULL;
  }
  SDP sdp = r.get();  // blocking
  peer->resetPromise();
  if (!sdp)
    return NULL;
  auto s = CGO_SerializeSDP(sdp);
  delete sdp;
  return s;
}

// PeerConnection::CreateAnswer
// Blocks until libwebrtc succeeds in generating the SDP answer.
// @returns SDP, or NULL on timeout.
CGO_sdpString CGO_CreateAnswer(CGO_Peer cgoPeer) {
  Peer *peer = (Peer*)cgoPeer;
  auto r = peer->promiseSDP.get_future();
  peer->pc_->CreateAnswer(peer, peer->constraints);
  if (SDPtimeout(&r, TIMEOUT_SECS)) {
    CGO_DBG("CreateAnswer timed out after " + TIMEOUT_SECS);
    peer->resetPromise();
    return NULL;
  }
  SDP sdp = r.get();  // blocking
  peer->resetPromise();
  if (!sdp)
    return NULL;
  auto s = CGO_SerializeSDP(sdp);
  delete sdp;
  return s;
}

// Serialize SDP message to a string Go can use.
CGO_sdpString CGO_SerializeSDP(CGO_sdp sdp) {
  SDP cSDP = (SDP)sdp;
  std::string s;
  cSDP->ToString(&s);
  return (CGO_sdpString)strdup(s.c_str());
}

// Given a fully serialized SDP string |msg|, return a CGO sdp object.
CGO_sdp CGO_DeserializeSDP(const char *type, const char *msg) {
  // TODO: Maybe use an enum instead of string for type.
  auto jsep_sdp = new JsepSessionDescription(type);
  SdpParseError err;
  std::string msg_str(msg);
  SdpDeserialize(msg_str, jsep_sdp, &err);
  return (CGO_sdp)jsep_sdp;
}

// PeerConnection::SetLocalDescription
int CGO_SetLocalDescription(CGO_Peer cgoPeer, CGO_sdp sdp) {
  PC cPC = ((Peer*)cgoPeer)->pc_;
  auto obs = PeerSDPObserver::Create();
  auto r = obs->promiseSet.get_future();
  cPC->SetLocalDescription(obs, (SDP)sdp);
  return r.get();
}

// PeerConnection::GetLocalDescription
CGO_sdp CGO_GetLocalDescription(CGO_Peer cgoPeer) {
  PC cPC = ((Peer*)cgoPeer)->pc_;
  return (CGO_sdp)cPC->local_description();
}

// PeerConnection::SetRemoteDescription
int CGO_SetRemoteDescription(CGO_Peer cgoPeer, CGO_sdp sdp) {
  PC cPC = ((Peer*)cgoPeer)->pc_;
  auto obs = PeerSDPObserver::Create();
  auto r = obs->promiseSet.get_future();
  cPC->SetRemoteDescription(obs, (SDP)sdp);
  return r.get();
}

// PeerConnection::GetRemoteDescription
CGO_sdp CGO_GetRemoteDescription(CGO_Peer cgoPeer) {
  PC cPC = ((Peer*)cgoPeer)->pc_;
  return (CGO_sdp)cPC->remote_description();
}

// PeerConnection::AddIceCandidate
int CGO_AddIceCandidate(CGO_Peer cgoPeer, CGO_IceCandidate *cgoIC) {
  PC cPC = ((Peer*)cgoPeer)->pc_;
  SdpParseError *error = nullptr;
  IceCandidateInterface *ic = webrtc::CreateIceCandidate(
    string(cgoIC->sdp_mid), cgoIC->sdp_mline_index, string(cgoIC->sdp), error);
  if (error || !ic) {
    CGO_DBG("SDP parse error");
    return FAILURE;
  }
  if (!cPC->AddIceCandidate(ic)) {
    CGO_DBG("Problem adding ICE candidate.");
    return FAILURE;
  }
  return SUCCESS;
}

// PeerConnection::signaling_state
int CGO_GetSignalingState(CGO_Peer cgoPeer) {
  PC cPC = ((Peer*)cgoPeer)->pc_;
  return cPC->signaling_state();
}

// PeerConnection::ice_connection_state (and more)
int CGO_IceConnectionState(CGO_Peer cgoPeer) {
  PC cPC = ((Peer*)cgoPeer)->pc_;
  return cPC->ice_connection_state();
}

// PeerConnection::ice_gathering_state
int CGO_IceGatheringState(CGO_Peer cgoPeer) {
  PC cPC = ((Peer*)cgoPeer)->pc_;
  return cPC->ice_gathering_state();
}

// PeerConnection::SetConfiguration
int CGO_SetConfiguration(CGO_Peer cgoPeer, CGO_Configuration* cgoConfig) {
  Peer *peer = (Peer*)cgoPeer;
  auto cConfig = castConfig_(cgoConfig);
  webrtc::RTCError *error = new webrtc::RTCError();
  bool success = peer->pc_->SetConfiguration(*cConfig, error);
  if (success) {
    peer->SetConfig(cConfig);
    return SUCCESS;
  }
  return (int) error->type();
}

// PeerConnection::CreateDataChannel
void* CGO_CreateDataChannel(CGO_Peer cgoPeer, char *label, CGO_DataChannelInit dict) {
  DataChannelInit config;
  config.ordered = dict.ordered;
  config.maxRetransmits = dict.maxRetransmits;
  config.maxRetransmitTime = dict.maxPacketLifeTime;
  config.negotiated = dict.negotiated;
  config.id = dict.id;

  auto cPeer = (Peer*)cgoPeer;
  std::string l(label);
  auto channel = cPeer->pc_->CreateDataChannel(l, &config);
  if (NULL == channel) {
    CGO_DBG("Unable to create DataChannel.");
    return NULL;
  }
  DCObserver obs = new rtc::RefCountedObject<CGoDataChannelObserver>(channel);
  cPeer->o_lock.lock();
  cPeer->observers.push_back(obs);
  cPeer->o_lock.unlock();
  auto o = obs.get();
  channel->RegisterObserver(o);
  return (void *)o;
}

// PeerConnection::Close
void CGO_Close(CGO_Peer peer) {
  auto cPeer = (Peer*)peer;
  cPeer->pc_->Close();
  CGO_DBG("Closed PeerConnection.");
}

void CGO_DeleteDataChannel(CGO_Peer cgoPeer, void* l) {
  auto cPeer = (Peer*)cgoPeer;
  auto o = (CGoDataChannelObserver*)l;
  auto observers = &cPeer->observers;
  cPeer->o_lock.lock();
  observers->erase(
    std::remove_if(observers->begin(), observers->end(), [o](DCObserver obs){
      return obs.get() == o;
    }),
    observers->end()
  );
  cPeer->o_lock.unlock();
}

//
// Test helpers which fake native callbacks.
//
void CGO_fakeIceCandidateError(CGO_Peer peer) {
  auto cPeer = (Peer*)peer;
  cPeer->OnIceConnectionChange(
      PeerConnectionInterface::IceConnectionState::kIceConnectionFailed);
}
