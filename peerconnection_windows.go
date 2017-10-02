// +build windows
package webrtc

/*
#include "peerconnection.h"
#include "string.h"
*/
import "C"
import (
	"unsafe"
)

func init() {
	// mod, err := syscall.LoadDLL("webrtc.dll")
}

//export dll_AddIceCandidate
func dll_AddIceCandidate(cgoPeer C.CGO_Peer, cgoIC *C.CGO_IceCandidate) C.int {
	r1, _, _ := myfuncs["dll_AddIceCandidate"].Call(uintptr(cgoPeer), uintptr(unsafe.Pointer(cgoIC)))
	return C.int(r1)
}

//export dll_CreateAnswer
func dll_CreateAnswer(cgoPeer C.CGO_Peer) C.CGO_sdpString {
	r1, _, _ := myfuncs["dll_CreateAnswer"].Call(uintptr(cgoPeer))
	// The caller is going to call free using C.free, that will invoke gcc's free which
	// is NOT the same as MSVC's free
	res := C.strdup((*C.char)(unsafe.Pointer(r1)))
	myfuncs["dll_Free"].Call(r1)
	return C.CGO_sdpString(res)
}

//export dll_CreateDataChannel
func dll_CreateDataChannel(cgoPeer C.CGO_Peer, label *C.char, dict unsafe.Pointer) unsafe.Pointer {
	r1, _, _ := myfuncs["dll_CreateDataChannel"].Call(uintptr(cgoPeer), uintptr(unsafe.Pointer(label)), uintptr(unsafe.Pointer(dict)))
	return unsafe.Pointer(r1)
}

//export dll_CreateOffer
func dll_CreateOffer(cgoPeer C.CGO_Peer) C.CGO_sdpString {
	r1, _, _ := myfuncs["dll_CreateOffer"].Call(uintptr(cgoPeer))
	// The caller is going to call free using C.free, that will invoke gcc's free which
	// is NOT the same as MSVC's free
	res := C.strdup((*C.char)(unsafe.Pointer(r1)))
	myfuncs["dll_Free"].Call(r1)
	return C.CGO_sdpString(res)
}

//export dll_GetLocalDescription
func dll_GetLocalDescription(cgoPeer C.CGO_Peer) C.CGO_sdp {
	r1, _, _ := myfuncs["dll_GetLocalDescription"].Call(uintptr(cgoPeer))
	return C.CGO_sdp(r1)
}

//export dll_GetRemoteDescription
func dll_GetRemoteDescription(cgoPeer C.CGO_Peer) C.CGO_sdp {
	r1, _, _ := myfuncs["dll_GetRemoteDescription"].Call(uintptr(cgoPeer))
	return C.CGO_sdp(r1)
}

//export dll_GetSignalingState
func dll_GetSignalingState(cgoPeer C.CGO_Peer) C.int {
	r1, _, _ := myfuncs["dll_GetSignalingState"].Call(uintptr(cgoPeer))
	return C.int(r1)
}

//export dll_IceConnectionState
func dll_IceConnectionState(cgoPeer C.CGO_Peer) C.int {
	r1, _, _ := myfuncs["dll_IceConnectionState"].Call(uintptr(cgoPeer))
	return C.int(r1)
}

//export dll_IceGatheringState
func dll_IceGatheringState(cgoPeer C.CGO_Peer) C.int {
	r1, _, _ := myfuncs["dll_IceGatheringState"].Call(uintptr(cgoPeer))
	return C.int(r1)
}

//export dll_InitializePeer
func dll_InitializePeer(goPc C.int) C.CGO_Peer {
	r1, _, _ := myfuncs["dll_InitializePeer"].Call(uintptr(goPc))
	return C.CGO_Peer(r1)
}

//export dll_SetLocalDescription
func dll_SetLocalDescription(cgoPeer C.CGO_Peer, sdp C.CGO_sdp) C.int {
	r1, _, _ := myfuncs["dll_SetLocalDescription"].Call(uintptr(cgoPeer), uintptr(unsafe.Pointer(sdp)))
	return C.int(r1)
}

//export dll_SetRemoteDescription
func dll_SetRemoteDescription(cgoPeer C.CGO_Peer, sdp C.CGO_sdp) C.int {
	r1, _, _ := myfuncs["dll_SetRemoteDescription"].Call(uintptr(cgoPeer), uintptr(unsafe.Pointer(sdp)))
	return C.int(r1)
}

//export dll_CreatePeerConnection
func dll_CreatePeerConnection(cgoPeer C.CGO_Peer, cgoConfig *C.CGO_Configuration) C.int {
	r1, _, _ := myfuncs["dll_CreatePeerConnection"].Call(uintptr(cgoPeer), uintptr(unsafe.Pointer(cgoConfig)))
	return C.int(r1)
}

//export dll_SetConfiguration
func dll_SetConfiguration(cgoPeer C.CGO_Peer, cgoConfig *C.CGO_Configuration) C.int {
	r1, _, _ := myfuncs["dll_SetConfiguration"].Call(uintptr(cgoPeer), uintptr(unsafe.Pointer(cgoConfig)))
	return C.int(r1)
}

//export dll_Close
func dll_Close(peer C.CGO_Peer) {
	myfuncs["dll_Close"].Call(uintptr(peer))
}

//export dll_fakeIceCandidateError
func dll_fakeIceCandidateError(peer C.CGO_Peer) {
	myfuncs["dll_fakeIceCandidateError"].Call(uintptr(peer))
}

//export dll_DeserializeSDP
func dll_DeserializeSDP(sdpType *C.char, msg *C.char) C.CGO_sdp {
	r1, _, _ := myfuncs["dll_DeserializeSDP"].Call(uintptr(unsafe.Pointer(sdpType)), uintptr(unsafe.Pointer(msg)))
	return C.CGO_sdp(unsafe.Pointer(r1))
}

//export dll_SerializeSDP
func dll_SerializeSDP(sdp C.CGO_sdp) C.CGO_sdpString {
	r1, _, _ := myfuncs["dll_SerializeSDP"].Call(uintptr(sdp))
	// The caller is going to call free using C.free, that will invoke gcc's free which
	// is NOT the same as MSVC's free
	res := C.strdup((*C.char)(unsafe.Pointer(r1)))
	myfuncs["dll_Free"].Call(r1)
	return C.CGO_sdpString(res)
}
