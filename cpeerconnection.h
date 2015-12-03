#ifndef _C_PEERCONNECTION_H_
#define _C_PEERCONNECTION_H_

#define WEBRTC_POSIX 1

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void (*Callback)();

  typedef void* PeerConnection;
  PeerConnection NewPeerConnection();

  // PeerConnectionInterface::IceServers
  void* GetIceServers(PeerConnection pc);

  int CreateOffer(PeerConnection pc);
  int CreateAnswer(PeerConnection pc);

#ifdef __cplusplus
}
#endif

#endif
