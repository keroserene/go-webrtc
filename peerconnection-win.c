#include <_cgo_export.h> // Allow calling certain Go functions.

int CGO_AddIceCandidate(CGO_Peer cgoPeer, CGO_IceCandidate *cgoIC) {
  return dll_AddIceCandidate(cgoPeer, cgoIC);
}

CGO_sdpString CGO_CreateAnswer(CGO_Peer cgoPeer) {
  return dll_CreateAnswer(cgoPeer);
}

void* CGO_CreateDataChannel(CGO_Peer cgoPeer, char *label, void *dict) {
  return dll_CreateDataChannel(cgoPeer, label, dict);
}

CGO_sdpString CGO_CreateOffer(CGO_Peer cgoPeer) {
  return dll_CreateOffer(cgoPeer);
}

int CGO_CreatePeerConnection(CGO_Peer cgoPeer, CGO_Configuration *cgoConfig) {
  return dll_CreatePeerConnection(cgoPeer, cgoConfig);
}

CGO_sdp CGO_GetLocalDescription(CGO_Peer cgoPeer) {
  return dll_GetLocalDescription(cgoPeer);
}

CGO_sdp CGO_GetRemoteDescription(CGO_Peer cgoPeer) {
  return dll_GetRemoteDescription(cgoPeer);
}

int CGO_GetSignalingState(CGO_Peer cgoPeer) {
  return dll_GetSignalingState(cgoPeer);
}

int CGO_IceConnectionState(CGO_Peer cgoPeer) {
  return dll_IceConnectionState(cgoPeer);
}

int CGO_IceGatheringState(CGO_Peer cgoPeer) {
  return dll_IceGatheringState(cgoPeer);
}

CGO_Peer CGO_InitializePeer(int goPc) {
  return dll_InitializePeer(goPc);
}

int CGO_SetConfiguration(CGO_Peer cgoPeer, CGO_Configuration* cgoConfig) {
  return dll_SetConfiguration(cgoPeer, cgoConfig);
}

int CGO_SetLocalDescription(CGO_Peer cgoPeer, CGO_sdp sdp) {
  return dll_SetLocalDescription(cgoPeer, sdp);
}

int CGO_SetRemoteDescription(CGO_Peer cgoPeer, CGO_sdp sdp) {
  return dll_SetRemoteDescription(cgoPeer, sdp);
}

void CGO_Close(CGO_Peer peer) {
  return dll_Close(peer);
}

void CGO_fakeIceCandidateError(CGO_Peer peer) {
  return dll_fakeIceCandidateError(peer);
}

CGO_sdp CGO_DeserializeSDP(const char *type, const char *msg) {
  return dll_DeserializeSDP((char*)type, (char*)msg);
}

CGO_sdpString CGO_SerializeSDP(CGO_sdp sdp) {
  return dll_SerializeSDP(sdp);
}