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

  typedef void* CGOPeer;
  typedef void* CGOsdp;  // Pointer to SessionDescriptionInterface*
  typedef void* CGODataChannel;
  typedef const char* CGOsdpString;

  CGOPeer CGOInitializePeer();
  // Below are "C methods" for the Peer class, which must be hidden from cgo.

  CGOPeer NewPeerConnection(CGOPeer cgoPeer);
  // PeerConnectionInterface::IceServers
  // void* GetIceServers(CGOPeePeerConnection pc);

  CGOsdp CGOCreateOffer(CGOPeer pc);
  CGOsdp CGOCreateAnswer(CGOPeer pc);

  CGOsdpString CGOSerializeSDP(CGOsdp sdp);
  int CGOSetLocalDescription(CGOPeer pc, CGOsdp sdp);
  int CGOSetRemoteDescription(CGOPeer pc, CGOsdp sdp);

  CGODataChannel CGOCreateDataChannel(CGOPeer pc, char *label, void *dict);

#ifdef __cplusplus
}
#endif

#endif
