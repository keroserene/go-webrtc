/**
 * C wrapper for the C++ PeerConnection code, to be go-compatible.
 */
#include "cpeerconnection.h"
#include "webrtc/base/common.h"
#include "talk/app/webrtc/peerconnectioninterface.h"
#include "talk/app/webrtc/test/fakeconstraints.h"
#include <iostream>
#include <unistd.h>
#include <future>
#include <string>


#define DTLS_SRTP "true"
#define SUCCESS 0
#define FAILURE -1
#define TIMEOUT_SECS 3

using namespace std;
using namespace webrtc;

// Smaller typedefs
typedef rtc::scoped_refptr<webrtc::PeerConnectionInterface> PC;
typedef SessionDescriptionInterface* SDP;
typedef rtc::scoped_refptr<DataChannelInterface> DataChannel;

// Peer acts as the glue between go and native code PeerConnectionInterface.
// However, it's not directly accessible from the Go side, which can only
// see what's exposed in the more pure extern "C" header file.
//
// The Go side may access this class through C.CGOPeer.
//
// This class also stubs libwebrtc's callback interface to be blocking,
// which allows the usage of goroutines, which is more idiomatic and easier
// for users of this library.
// The alternative would require casting Go function pointers, calling Go code
// from C code from Go code, which is less likely to be a good time.
//
// TODO(keroserene): More documentation...
// TODO: Better logging
class Peer
  : public CreateSessionDescriptionObserver,
    public PeerConnectionObserver {
 public:

  bool Initialize() {
    // Prepare everything.
    // Should be called before anything else happens.
    promiseSDP = promise<SDP>();
    // Due to the different threading model, in order for PeerConnectionFactory
    // to be able to post async messages without getting blocked, we need to use
    // external signalling and worker threads.
    signal_thread = new rtc::Thread();
    worker_thread = new rtc::Thread();
    signal_thread->Start();
    worker_thread->Start();
    pc_factory = CreatePeerConnectionFactory(
      signal_thread,
      worker_thread,
      NULL, NULL, NULL);
    if (!pc_factory.get()) {
      cout << "ERROR: Could not create PeerConnectionFactory" << endl;
      return false;
    }
    // PortAllocatorFactoryInterface *allocator;

    // TODO: Make actual media constraints, with an exposed Go interface.
    auto c = new FakeConstraints();
    c->AddOptional(MediaConstraintsInterface::kEnableDtlsSrtp, DTLS_SRTP);
    c->SetMandatoryReceiveAudio(false);
    c->SetMandatoryReceiveVideo(false);
    constraints = c;
    // cout << "[C] Peer initialized." << endl;
    return true;
  }

  void resetPromise() {
    promiseSDP = promise<SDP>();
  }

  //
  // CreateSessionDescriptionObserver implementation
  //
  void OnSuccess(SDP desc) {
    cout << "[C] SDP successfully created at " << desc << endl;
    promiseSDP.set_value(desc);
  }

  void OnFailure(const std::string& error) {
    cout << "[C] SDP Failure: " << error << endl;
    promiseSDP.set_value(NULL);
  }

  //
  // PeerConnectionObserver Implementation
  // TODO: cgo hooks
  //
  void OnStateChange(PeerConnectionObserver::StateType state) {
    cout << "[C] OnStateChange: " << state << endl;
  }

  void OnAddStream(webrtc::MediaStreamInterface* stream) {
    cout << "[C] OnAddStream: " << stream << endl;
  }

  void OnRemoveStream(webrtc::MediaStreamInterface* stream) {
    cout << "[C] OnRemoveStream: " << stream << endl;
  }

  void OnRenegotiationNeeded() {
    cout << "[C] OnRenegotiationNeeded" << endl;
  }

  void OnIceCandidate(const IceCandidateInterface* candidate) {
    cout << "[C] OnIceCandidate" << candidate << endl;
  }

  void OnDataChannel(DataChannelInterface* data_channel) {
    cout << "[C] OnDataChannel: " << data_channel << endl;
  }

  PeerConnectionInterface::RTCConfiguration *config;
  PeerConnectionInterface::RTCOfferAnswerOptions options;
  const MediaConstraintsInterface* constraints;

  PC pc_;

  // Passing SDPs through promises instead of callbacks, to allow the benefits
  // as described above.
  // However, this has the effect that CreateOffer and CreateAnswer must not be
  // concurrent, to themselves or each other (which isn't expected anyways) due
  // to the simplistic way in which futures are used here.
  promise<SDP> promiseSDP;

  rtc::scoped_refptr<PeerConnectionFactoryInterface> pc_factory;
  // TODO: prepare and expose IceServers for real.
  // PeerConnectionInterface::IceServers ice_servers;

 protected:
  rtc::Thread *signal_thread;
  rtc::Thread *worker_thread;

};  // class Peer

// Keep track of Peers in global scope to prevent deallocation, due to the
// required scoped_refptr from implementing the Observer interface.
vector<rtc::scoped_refptr<Peer>> localPeers;


// TODO: Make a better generalized class for every "Observer" later.
class PeerSDPObserver : public SetSessionDescriptionObserver {
 public:
  static PeerSDPObserver* Create() {
    return new rtc::RefCountedObject<PeerSDPObserver>();
  }
  virtual void OnSuccess() {
    // cout << "[C] SDP Set Success!" << endl;
    promiseSet.set_value(0);
  }
  virtual void OnFailure(const std::string& error) {
    cout << "[C ERROR] SessionDescription: " << error << endl;
    promiseSet.set_value(-1);
  }
  promise<int> promiseSet = promise<int>();

 protected:
  PeerSDPObserver() {}
  ~PeerSDPObserver() {}

};  // class PeerSDPObserver


//
// extern "C" Go-accessible functions.
//

// Create and return the Peer object, which provides initial native code
// glue for the PeerConnection constructor.
CGOPeer CGOInitializePeer() {
  rtc::scoped_refptr<Peer> localPeer = new rtc::RefCountedObject<Peer>();
  localPeer->Initialize();
  localPeers.push_back(localPeer);
  return localPeer;
}

// This helper converts RTCConfiguration struct from GO to C++.
PeerConnectionInterface::RTCConfiguration *castConfig_(
    CGORTCConfiguration *cgoConfig) {
  PeerConnectionInterface::RTCConfiguration* c =
      new PeerConnectionInterface::RTCConfiguration();

  // TODO: Parse Go ice server slice into C++ vector of IceServer structs.
  PeerConnectionInterface::IceServer *server = new
      PeerConnectionInterface::IceServer();
  server->uri = "stun:stun.l.google.com:19302";
  c->servers.push_back(*server);

  // Fragile cast from Go const "enum" to C++ Enum based on int ordering assumptions.
  // May need something better later.
  c->type = (PeerConnectionInterface::IceTransportsType)cgoConfig->IceTransportPolicy;
  c->bundle_policy = (PeerConnectionInterface::BundlePolicy)cgoConfig->BundlePolicy;
  return c;
}

// |Peer| method: create a native code PeerConnection object.
// Returns 0 on Success.
int CGOCreatePeerConnection(CGOPeer cgoPeer, CGORTCConfiguration *cgoConfig) {
  Peer *peer = (Peer*)cgoPeer;
  peer->config = castConfig_(cgoConfig);
  // cout << "RTCConfiguration: " << peer->config << endl;

  // Prepare a native PeerConnection object.
  peer->pc_ = peer->pc_factory->CreatePeerConnection(
    *peer->config,
    peer->constraints,
    NULL,  // port allocator
    NULL,  // TODO: DTLS
    peer
    );

  if (!peer->pc_.get()) {
    cout << "ERROR: Could not create PeerConnection." << endl;
    return FAILURE;
  }
  // cout << "[C] Made PeerConnection: " << peer->pc_ << endl;
  return SUCCESS;
}

bool SDPtimeout(future<SDP> *f, int seconds) {
  auto status = f->wait_for(chrono::seconds(TIMEOUT_SECS));
  return future_status::ready != status;
}

// PeerConnection::CreateOffer
// Blocks until libwebrtc succeeds in generating the SDP offer,
// @returns SDP (pointer), or NULL on timeeout.
CGOsdp CGOCreateOffer(CGOPeer cgoPeer) {
  // TODO: Provide an actual RTCOfferOptions as an argument.
  Peer* peer = (Peer*)cgoPeer;
  auto r = peer->promiseSDP.get_future();
  peer->pc_->CreateOffer(peer, peer->constraints);
  if (SDPtimeout(&r, TIMEOUT_SECS)) {
    cout << "[C] CreateOffer timed out after " << TIMEOUT_SECS << endl;
    peer->resetPromise();
    return NULL;
  }
  SDP sdp = r.get();  // blocking
  peer->resetPromise();
  return (CGOsdp)sdp;
}


// PeerConnection::CreateAnswer
// Blocks until libwebrtc succeeds in generating the SDP answer.
// @returns SDP, or NULL on timeout.
CGOsdp CGOCreateAnswer(CGOPeer cgoPeer) {
  Peer *peer = (Peer*)cgoPeer;
  cout << "[C] CreateAnswer" << peer << endl;
  auto r = peer->promiseSDP.get_future();
  peer->pc_->CreateAnswer(peer, peer->constraints);
  if (SDPtimeout(&r, TIMEOUT_SECS)) {
    cout << "[C] CreateAnswer timed out after " << TIMEOUT_SECS << endl;
    peer->resetPromise();
    return NULL;
  }
  SDP sdp = r.get();  // blocking
  peer->resetPromise();
  return (CGOsdp)sdp;
}


// Serialize SDP message to a string Go can use.
CGOsdpString CGOSerializeSDP(CGOsdp sdp) {
  auto s = new string();
  SDP cSDP = (SDP)sdp;
  cSDP->ToString(s);
  return (CGOsdpString)s->c_str();
}

int CGOSetLocalDescription(CGOPeer pc, CGOsdp sdp) {
  PC cPC = ((Peer*)pc)->pc_;
  auto obs = PeerSDPObserver::Create();
  auto r = obs->promiseSet.get_future();
  cPC->SetLocalDescription(obs, (SDP)sdp);
  return r.get();
}

int CGOSetRemoteDescription(CGOPeer pc, CGOsdp sdp) {
  PC cPC = ((Peer*)pc)->pc_;
  auto obs = PeerSDPObserver::Create();
  auto r = obs->promiseSet.get_future();
  cPC->SetRemoteDescription(obs, (SDP)sdp);
  return r.get();
}


CGODataChannel CGOCreateDataChannel(CGOPeer pc, char *label, void *dict) {
  PC cPC = ((Peer*)pc)->pc_;
  DataChannelInit *r = (DataChannelInit*)dict;
  // TODO: a real config struct, with correct fields
  DataChannelInit config;
  string *l = new string(label);
  auto channel = cPC->CreateDataChannel(*l, &config);
  cout << "Created data channel: " << channel << endl;
  return (CGODataChannel)channel;
}

