#ifndef _C_PEERCONNECTION_H_
#define _C_PEERCONNECTION_H_

#define WEBRTC_POSIX 1

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void* CGOArray;

  typedef void* CGOPeer;
  typedef void* CGOsdp;  // Pointer to SessionDescriptionInterface*
  typedef void* CGODataChannel;
  typedef const char* CGOsdpString;

  typedef struct {
    CGOArray IceServers;
    int      IceTransportPolicy;
    int      BundlePolicy;
    int      RtcpMuxPolicy;
    char     *PeerIdentity;
    CGOArray Certificates;
    int      IceCandidatePoolSize;
  } CGORTCConfiguration;

  CGOPeer CGOInitializePeer();
  // Below are "C methods" for the Peer class, which must be hidden from cgo.

  int CGOCreatePeerConnection(CGOPeer, CGORTCConfiguration*);
  // PeerConnectionInterface::IceServers
  // void* GetIceServers(CGOPeePeerConnection pc);

  CGOsdp CGOCreateOffer(CGOPeer);
  CGOsdp CGOCreateAnswer(CGOPeer);

  CGOsdpString CGOSerializeSDP(CGOsdp);
  int CGOSetLocalDescription(CGOPeer, CGOsdp);
  int CGOSetRemoteDescription(CGOPeer, CGOsdp);

  CGODataChannel CGOCreateDataChannel(CGOPeer, char*, void*);

#ifdef __cplusplus
}
#endif

#endif
