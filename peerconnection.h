#ifndef _C_PEERCONNECTION_H_
#define _C_PEERCONNECTION_H_

#define WEBRTC_POSIX 1

#include <stdint.h>

#include "mediastreamtrack.h"
#include "mediastream.h"
#include "rtpsender.h"

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void* CGO_Array;

  typedef void* CGO_Peer;
  typedef void* CGO_sdp;  // Pointer to SessionDescriptionInterface*
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

  typedef struct {
    const char *sdp_mid;
    int sdp_mline_index;
    const char *sdp;
  } CGO_IceCandidate;

  CGO_Peer CGO_InitializePeer(int pc);

  // Below are "C methods" for the Peer class, which must be hidden from cgo.

  int CGO_CreatePeerConnection(CGO_Peer, CGO_Configuration*);

  void CGO_PullAudio(CGO_Peer, uint16_t*, int, uint32_t);

  CGO_sdpString CGO_CreateOffer(CGO_Peer);
  CGO_sdpString CGO_CreateAnswer(CGO_Peer);

  CGO_sdpString CGO_SerializeSDP(CGO_sdp);
  CGO_sdp CGO_DeserializeSDP(const char *type, const char *msg);

  int CGO_SetLocalDescription(CGO_Peer, CGO_sdp);
  CGO_sdp CGO_GetLocalDescription(CGO_Peer);
  int CGO_SetRemoteDescription(CGO_Peer, CGO_sdp);
  CGO_sdp CGO_GetRemoteDescription(CGO_Peer);
  int CGO_AddIceCandidate(CGO_Peer cgoPeer, CGO_IceCandidate *cgoIC);

  int CGO_GetSignalingState(CGO_Peer);
  int CGO_IceConnectionState(CGO_Peer);
  int CGO_IceGatheringState(CGO_Peer);
  int CGO_SetConfiguration(CGO_Peer, CGO_Configuration*);

  CGO_RtpSender CGO_PeerConnection_AddTrack(CGO_Peer, CGO_MediaStreamTrack, CGO_MediaStream*, int);
  void CGO_PeerConnection_RemoveTrack(CGO_Peer, CGO_RtpSender);
  void* CGO_CreateDataChannel(CGO_Peer, char*, void*);

  // TODO: Move to mediastream.cc (when pc_factory is factored out of Peer).
  CGO_MediaStream CGO_NewMediaStream(CGO_Peer cgoPeer, const char* label);

  void CGO_Close(CGO_Peer);

  // Test helpers
  void CGO_fakeIceCandidateError(CGO_Peer peer);

#ifdef __cplusplus
}
#endif

#endif  // _C_PEERCONNECTION_H_
