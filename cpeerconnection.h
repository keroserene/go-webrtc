#ifndef _C_PEERCONNECTION_H_
#define _C_PEERCONNECTION_H_

#define WEBRTC_POSIX 1

#ifdef __cplusplus
extern "C" {
#endif

  typedef void (*Callback)();

  typedef void* PeerConnection;
  PeerConnection NewPeerConnection();

  // PeerConnectionInterface::IceServers
  void* GetIceServers(PeerConnection pc);

  // void CreateOffer(PeerConnection pc, Callback onsuccess, Callback onfailure);
  int CreateOffer(PeerConnection pc);
  void CreateAnswer(PeerConnection pc, void* callback);

#ifdef __cplusplus
}
#endif

#endif
