/**
 * C wrapper for the C++ PeerConnection code, to be go-compatible.
 */
#include "cpeerconnection.h"
// #include "webrtc/base/thread.h"
#include "webrtc/base/common.h"
#include "webrtc/base/common.h"
// #include "talk/app/webrtc/peerconnection.h"
// #include "talk/app/webrtc/peerconnectionfactory.h"
#include "talk/app/webrtc/peerconnectioninterface.h"
#include <iostream>
#include <unistd.h>

#define SUCCESS 0
#define FAILURE 1

using namespace std;
using namespace webrtc;

typedef rtc::scoped_refptr<webrtc::PeerConnectionInterface> PC;

// Peer acts as the glue between go and native code PeerConnectionInterface.
// However, it is not directly accessible through cgo.
class Peer
  : public CreateSessionDescriptionObserver,
    public PeerConnectionObserver {

 public:

  /*
   * Stub out all callbacks to become blocking, and return boolean success / fail.
   * Since the user wants to write go code, it'd be better to support goroutines
   * instead of callbacks.
   * This prevents the complication of casting Go function pointers and
   * then dealing with the risk of concurrently calling Go code from C from Go...
   * Which should be a much easier and safer for users of this library.
   * TODO(keroserene): Expand on this if there are more complicated callbacks.
   */
  Callback SuccessCallback = NULL;
  Callback FailureCallback = NULL;

  // CreateSessionDescriptionObserver implementation
  void OnSuccess(SessionDescriptionInterface* desc) {
    cout << "success" << endl;
    // if (this->SuccessCallback) {
      // this->SuccessCallback();
    // }
  }
  void OnFailure(const std::string& error) {
    cout << "failure" << endl;
    // if (this->FailureCallback) {
      // this->FailureCallback();
    // }
  }
  int AddRef() const {}
  int Release() const {}

  // PeerConnectionObserver Implementation
  void OnStateChange(PeerConnectionObserver::StateType state) {
    cout << "OnStateChange" << endl;
  }

  void OnAddStream(webrtc::MediaStreamInterface* stream) {
    cout << "OnAddStream" << endl;
  }

  void OnRemoveStream(webrtc::MediaStreamInterface* stream) {
    cout << "OnRemoveStream" << endl;
  }

  void OnRenegotiationNeeded() {
    cout << "OnRenegotiationNeeded" << endl;
  }

  void OnIceCandidate(const IceCandidateInterface* candidate) {
    cout << "OnIceCandidate" << endl;
  }

  void OnDataChannel(DataChannelInterface* data_channel) {
    cout << "OnDataChannel" << endl;
  }

  PeerConnectionInterface::RTCConfiguration *config;
  PeerConnectionInterface::RTCOfferAnswerOptions options;
  const MediaConstraintsInterface* constraints;

  PC pc_;  // This scoped_refptr must live in an object of some sort, or it will
           // be prematurely deallocated.

};

// TODO: Wrap as much as possible within the Peer class?
Peer *peer;


// Create and return a PeerConnection object.
PeerConnection NewPeerConnection() {

  rtc::scoped_refptr<PeerConnectionFactoryInterface> pc_factory;
  // TODO: Need to use a different constructor, later.
  pc_factory = CreatePeerConnectionFactory();
  if (!pc_factory.get()) {
    cout << "ERROR: Could not create PeerConnectionFactory" << endl;
    return NULL;
  }
  // PortAllocatorFactoryInterface *allocator;
  peer = new Peer();

  // TODO: prepare and expose IceServers for real.
  PeerConnectionInterface::IceServers ice_servers;
  PeerConnectionInterface::IceServer ice_server;
  ice_server.uri = "stun:stun.l.google.com:19302";
  ice_servers.push_back(ice_server);

  // Prepare RTC Configuration object. This is just the default one, for now.
  // TODO: A Go struct that can be passed and converted here.
  // TODO: Memory leak.
  peer->config = new PeerConnectionInterface::RTCConfiguration();
  peer->config->servers = ice_servers;
  // cout << "Preparing RTCConfiguration..." << peer->config << endl;
  // TODO: DTLS Certificates

  // Prepare a native PeerConnection object.
  peer->pc_ = pc_factory->CreatePeerConnection(
    *peer->config,
    peer->constraints,
    NULL, // port allocator
    NULL, // dtls
    peer
    );
  if (!peer->pc_.get()) {
    cout << "ERROR: Could not create PeerConnection." << endl;
    return NULL;
  }
  cout << "[C] Made a PeerConnection: " << peer->pc_ << endl;
  return (PeerConnection)peer;
}


// Blocking version of CreateOffer (or, will be soon)
// Returns 0 on success, -1 on failure.
int CreateOffer(PeerConnection pc) {
  PC cPC = ((Peer*)pc)->pc_;
  cout << "[C] CreateOffer" << peer << endl;

  // TODO: Provide an actual RTCOfferOptions as an argument.
  cPC->CreateOffer(peer, peer->options);

  // TODO: Up in PeerConnectionFactory, should probably use custom threads in
  // order for the callbacks to be *actually* registered correctly.
  cout << "[C] CreateOffer done!" << endl;
  return SUCCESS;
}

void CreateAnswer(PeerConnection pc, void* callback) {
}
