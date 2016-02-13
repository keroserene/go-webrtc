/**
 * C wrapper around the C++ webrtc::PeerConnectionInterface and related, which
 * allows compatibility with CGO's requirements so that everything may
 * ultimately be exposed in Go.
 */
#include "cpeerconnection.h"
#include "_cgo_export.h"  // Allow calling certain Go functions.

#include "talk/app/webrtc/peerconnectioninterface.h"
#include "talk/app/webrtc/test/fakeconstraints.h"
#include "talk/app/webrtc/jsepsessiondescription.h"
#include "talk/app/webrtc/webrtcsdp.h"

#include <iostream>
#include <future>

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
typedef rtc::scoped_refptr<DataChannelInterface> DataChannel;

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
    signalling_thread_ = new rtc::Thread();
    worker_thread_ = new rtc::Thread();
    signalling_thread_->SetName("CGO Signalling", NULL);
    worker_thread_->SetName("CGO Worker", NULL);
    signalling_thread_->Start();  // Must start before being passed to
    worker_thread_->Start();      // PeerConnectionFactory.

    pc_factory = CreatePeerConnectionFactory(
      worker_thread_,
      signalling_thread_,
      NULL, NULL, NULL);
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
    // TODO: This may need to be slightly more complicated...
    // https://w3c.github.io/webrtc-pc/#rtcpeerconnectionstate-enum
    cgoOnConnectionStateChange(goPeerConnection, new_state);
  }

  void OnIceGatheringChange(
      PeerConnectionInterface::IceGatheringState new_state) {
    cgoOnIceGatheringStateChange(goPeerConnection, new_state);
  }

  void OnDataChannel(DataChannelInterface* data_channel) {
    this->channels.push_back(data_channel);
    cgoOnDataChannel(goPeerConnection, data_channel);
  }

  // Note that Configuration is where ICE servers are specified.
  PeerConnectionInterface::RTCConfiguration *config;
  const MediaConstraintsInterface* constraints;

  PC pc_;                  // Pointer to webrtc::PeerConnectionInterface.
  void *goPeerConnection;  // Pointer to external Go PeerConnection struct,
                           // which is required to fire callbacks correctly.

  // Pass SDPs through promises instead of callbacks, to allow benefits as
  // described above. However, this means CreateOffer and CreateAnswer must
  // not be concurrent to themselves or each other (which isn't expected
  // anyways), due to the simplistic way futures are used here.
  promise<SDP> promiseSDP;

  rtc::scoped_refptr<PeerConnectionFactoryInterface> pc_factory;

  // Prevent deallocation of created DataChannels, since they are ref_ptr,
  // by keeping track of them in a vector.
  vector<DataChannel> channels;

 private:
  rtc::Thread *signalling_thread_;
  rtc::Thread *worker_thread_;

};  // class Peer

// Keep track of Peers in global scope to prevent deallocation, due to the
// required scoped_refptr from implementing the Observer interface.
vector<rtc::scoped_refptr<Peer>> localPeers;

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
CGO_Peer CGO_InitializePeer(void *goPc) {
  rtc::scoped_refptr<Peer> localPeer = new rtc::RefCountedObject<Peer>();
  localPeer->Initialize();
  localPeers.push_back(localPeer);
  localPeer->goPeerConnection = goPc;
  return localPeer;
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
    c->servers.push_back({
      "",  // TODO: Remove once webrtc deprecates the first uri field.
      urls,
      s.username,
      s.credential
    });
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
  peer->config = castConfig_(cgoConfig);
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
CGO_sdp CGO_CreateOffer(CGO_Peer cgoPeer) {
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
  return (CGO_sdp)sdp;
}

// PeerConnection::CreateAnswer
// Blocks until libwebrtc succeeds in generating the SDP answer.
// @returns SDP, or NULL on timeout.
CGO_sdp CGO_CreateAnswer(CGO_Peer cgoPeer) {
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
  return (CGO_sdp)sdp;
}

// Serialize SDP message to a string Go can use.
CGO_sdpString CGO_SerializeSDP(CGO_sdp sdp) {
  auto s = new string();
  SDP cSDP = (SDP)sdp;
  cSDP->ToString(s);
  char *cstr = (char *)malloc(s->length()+1 * sizeof(char));
  std::strcpy(cstr, s->c_str());
  return (CGO_sdpString)cstr;
}

// Given a fully serialized SDP string |msg|, return a CGO sdp object.
CGO_sdp CGO_DeserializeSDP(const char *type, const char *msg) {
  // TODO: Maybe use an enum instead of string for type.
  auto jsep_sdp = new JsepSessionDescription(type);
  SdpParseError err;
  auto msg_str = new string(msg);
  SdpDeserialize(*msg_str, jsep_sdp, &err);
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
  bool success = peer->pc_->SetConfiguration(*cConfig);
  if (success) {
    peer->config = cConfig;
    return SUCCESS;
  }
  return FAILURE;
}

// PeerConnection::CreateDataChannel
CGO_Channel CGO_CreateDataChannel(CGO_Peer cgoPeer, char *label, void *dict) {
  auto cPeer = (Peer*)cgoPeer;
  DataChannelInit *r = (DataChannelInit*)dict;
  // TODO: a real DataChannelInit config with correct fields.
  DataChannelInit config;
  string *l = new string(label);
  auto channel = cPeer->pc_->CreateDataChannel(*l, &config);
  if (NULL == channel) {
    CGO_DBG("Unable to create DataChannel.");
    return NULL;
  }
  cPeer->channels.push_back(channel);
  webrtc::DataChannelInterface* c = channel.get();
  return c;
}

// PeerConnection::Close
void CGO_Close(CGO_Peer peer) {
  auto cPeer = (Peer*)peer;
  cPeer->pc_->Close();
  CGO_DBG("Closed PeerConnection.");
}


//
// Test helpers which fake native callbacks.
//
void CGO_fakeIceCandidateError(CGO_Peer peer) {
  auto cPeer = (Peer*)peer;
  cPeer->OnIceConnectionChange(
      PeerConnectionInterface::IceConnectionState::kIceConnectionFailed);
}
