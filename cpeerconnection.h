#ifndef _C_PEERCONNECTION_H_
#define _C_PEERCONNECTION_H_

#define WEBRTC_POSIX 1

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void* CGO_Array;

  typedef void* CGO_Peer;
  typedef void* CGO_sdp;  // Pointer to SessionDescriptionInterface*
  typedef void* CGO_DataChannel;
  typedef const char* CGO_sdpString;

  typedef struct {
    char **urls;
    int   numUrls;

    char  *username;
    char  *credential;
  } CGO_IceServer;

  typedef struct {
    CGO_IceServer  *iceServers;
    int            numIceServers;

    int            iceTransportPolicy;
    int            bundlePolicy;
    // [BD] int      RtcpMuxPolicy;
    char           *peerIdentity;
    // [BD] CGO_Array Certificates;
    // [BD] int      IceCandidatePoolSize;
  } CGO_Configuration;

  CGO_Peer CGO_InitializePeer(void *pc);

  // Below are "C methods" for the Peer class, which must be hidden from cgo.

  int CGO_CreatePeerConnection(CGO_Peer, CGO_Configuration*);
  // PeerConnectionInterface::IceServers
  // void* GetIceServers(CGO_PeePeerConnection pc);

  CGO_sdp CGO_CreateOffer(CGO_Peer);
  CGO_sdp CGO_CreateAnswer(CGO_Peer);

  CGO_sdpString CGO_SerializeSDP(CGO_sdp);
  int CGO_SetLocalDescription(CGO_Peer, CGO_sdp);
  int CGO_SetRemoteDescription(CGO_Peer, CGO_sdp);

  int CGO_GetSignalingState(CGO_Peer);
  // int CGO_GetConfiguration(CGO_Peer);
  int CGO_SetConfiguration(CGO_Peer pc, CGO_Configuration*);

  CGO_DataChannel CGO_CreateDataChannel(CGO_Peer, char*, void*);

#ifdef __cplusplus
}
#endif

#endif
