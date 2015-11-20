/**
 * C wrapper for the C++ PeerConnection code, to be go-compatible.
 */
#include "cpeerconnection.h"
#include "webrtc/base/thread.h"
#include "talk/app/webrtc/peerconnection.h"
// #include "talk/app/webrtc/peerconnectionfactory.h"
#include <iostream>

using namespace std;
using namespace webrtc;

CPeerConnection NewPeerConnection() {
  const MediaConstraintsInterface* constraints;
  PortAllocatorFactoryInterface *allocator;
  rtc::scoped_refptr<PeerConnectionFactoryInterface> pc_factory =
      webrtc::CreatePeerConnectionFactory(
          rtc::Thread::Current(),
          rtc::Thread::Current(),
          NULL, NULL, NULL);

  // prepare ICE servers
  // TODO: expose this
  PeerConnectionInterface::IceServers ice_servers;
  PeerConnectionInterface::IceServer ice_server;
  ice_server.uri = "stun:stun.l.google.com:19302";
  ice_servers.push_back(ice_server);

  rtc::scoped_refptr<webrtc::PeerConnectionInterface> pc; 
  pc = pc_factory->CreatePeerConnection(
    ice_servers,
    constraints,
    allocator,  // port allocator,
    NULL, // dtls
    NULL // pc observer
    );
  return (void *)pc;
}
