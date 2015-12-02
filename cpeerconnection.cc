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

#define SUCCESS 0
#define FAILURE 1

using namespace std;
using namespace webrtc;

typedef rtc::scoped_refptr<webrtc::PeerConnectionInterface> PC;
const MediaConstraintsInterface* constraints;

/*
 * Stub out all callbacks to become blocking, and return boolean success / fail.
 * Since the user wants to write go code, it'd be better to support goroutines
 * instead of callbacks.
 * This prevents the complication of casting Go function pointers and
 * then dealing with the risk of concurrently calling Go code from C from Go...
 * Which should be a much easier and safer for users of this library.
 * TODO(keroserene): Expand on this if there are more complicated callbacks.
 */
class Callbacks : CreateSessionDescriptionObserver {
 public:
  // void (*SuccessCallback)() = NULL;
  // void (*FailureCallback)() = NULL;
  Callback SuccessCallback = NULL;
  Callback FailureCallback = NULL;
  void OnSuccess(SessionDescriptionInterface* desc) {
    cout << "success" << endl;
    if (this->SuccessCallback) {
      this->SuccessCallback();
    }
  }
  void OnFailure(const std::string& error) {
    if (this->FailureCallback) {
      this->FailureCallback();
    }
  }
  int AddRef() const {}
  int Release() const {}
};

PeerConnection NewPeerConnection() {

  rtc::scoped_refptr<PeerConnectionFactoryInterface> pc_factory;
  pc_factory = CreatePeerConnectionFactory();

  PortAllocatorFactoryInterface *allocator;

  /*
  rtc::scoped_refptr<PeerConnectionFactoryInterface> pc_factory =
      webrtc::CreatePeerConnectionFactory(
          rtc::Thread::Current(),
          rtc::Thread::Current(),
          NULL, NULL, NULL);
  */

  // prepare ICE servers
  // TODO: expose this
  PeerConnectionInterface::IceServers ice_servers;
  PeerConnectionInterface::IceServer ice_server;
  ice_server.uri = "stun:stun.l.google.com:19302";
  ice_servers.push_back(ice_server);

  cout << ice_server.uri << endl;

  // rtc::scoped_refptr<webrtc::PeerConnectionInterface> pc;
  PC pc;
  pc = pc_factory->CreatePeerConnection(
    ice_servers,
    constraints,
    allocator,  // port allocator,
    NULL, // dtls
    NULL // pc observer
    );
  return (void *)pc;
  // return pc;
}

// void CreateOffer(PeerConnection pc, Callback onsuccess, Callback onfailure) {
/*
 * Blocking version of CreateOffer:
 * Returns 0 on success, -1 on failure.
 */
int CreateOffer(PeerConnection pc) {
  // rtc::scoped_refptr<webrtc::PeerConnectionInterface>pc
  PC *cPC = (PC*)pc;
  Callbacks *obs = new Callbacks();
  // (CreateSessionDescriptionObserver*)callback;
  cout << "[c] CreateOffer" << endl;
  // Constraints...
  // cPC->get()->CreateOffer((CreateSessionDescriptionObserver*)obs, NULL);
  cout << "[c] CreateOffer done" << endl;
  return SUCCESS;
}

void CreateAnswer(PeerConnection pc, void* callback) {
}

// PeerConnectionInterface::IceServers GetIceServers(PeerConnection pc) {
  // return pc.ice_servers;
// }
