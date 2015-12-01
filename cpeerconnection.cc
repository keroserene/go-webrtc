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

using namespace std;
using namespace webrtc;

typedef rtc::scoped_refptr<webrtc::PeerConnectionInterface> PC;
const MediaConstraintsInterface* constraints;

class Callbacks : CreateSessionDescriptionObserver {
 public:
  void OnSuccess() {
  }
  void OnFailure() {
  }
  int AddRef() {}
  int Release() {}
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

void CreateOffer(PeerConnection pc, void(*onsuccess), void(*onfailure)) {
  // rtc::scoped_refptr<webrtc::PeerConnectionInterface>pc
  cout << "[C] CreateOffer callback is " << onsuccess << onfailure << endl;
  PC *cPC = (PC*)pc;
  CreateSessionDescriptionObserver obs = new Callbacks();
  obs.OnSuccess = onsuccess;
  obs.OnFailure = onfailure;
  // (CreateSessionDescriptionObserver*)callback;
  cout << "[c] CreateOffer" << endl;
  // Constraints...
  cPC->get()->CreateOffer(obs, NULL);
  cout << "[c] CreateOffer done" << endl;
}

void CreateAnswer(PeerConnection pc, void* callback) {
}

// PeerConnectionInterface::IceServers GetIceServers(PeerConnection pc) {
  // return pc.ice_servers;
// }
