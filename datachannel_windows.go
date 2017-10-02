// +build windows
package webrtc

/*
#include "datachannel.h"
#include "string.h"
*/
import "C"
import (
	"unsafe"
)

//export dll_Channel_BufferedAmount
func dll_Channel_BufferedAmount(channel C.CGO_Channel) C.int {
	r1, _, _ := myfuncs["dll_Channel_BufferedAmount"].Call(uintptr(channel))
	return C.int(r1)
}

//export dll_Channel_ID
func dll_Channel_ID(channel C.CGO_Channel) C.int {
	r1, _, _ := myfuncs["dll_Channel_ID"].Call(uintptr(channel))
	return C.int(r1)
}

//export dll_Channel_Label
func dll_Channel_Label(channel C.CGO_Channel) *C.char {
	r1, _, _ := myfuncs["dll_Channel_Label"].Call(uintptr(channel))
	// The caller is going to call free using C.free, that will invoke gcc's free which
	// is NOT the same as MSVC's free
	res := C.strdup((*C.char)(unsafe.Pointer(r1)))
	myfuncs["dll_Free"].Call(r1)
	return res
}

//export dll_Channel_Protocol
func dll_Channel_Protocol(channel C.CGO_Channel) *C.char {
	r1, _, _ := myfuncs["dll_Channel_Protocol"].Call(uintptr(channel))
	// The caller is going to call free using C.free, that will invoke gcc's free which
	// is NOT the same as MSVC's free
	res := C.strdup((*C.char)(unsafe.Pointer(r1)))
	myfuncs["dll_Free"].Call(r1)
	return res
}

//export dll_Channel_RegisterObserver
func dll_Channel_RegisterObserver(o unsafe.Pointer, index C.int) unsafe.Pointer {
	r1, _, _ := myfuncs["dll_Channel_RegisterObserver"].Call(uintptr(o), uintptr(index))
	return unsafe.Pointer(r1)
}

//export dll_Channel_MaxRetransmitTime
func dll_Channel_MaxRetransmitTime(channel C.CGO_Channel) C.int {
	r1, _, _ := myfuncs["dll_Channel_MaxRetransmitTime"].Call(uintptr(channel))
	return C.int(r1)
}

//export dll_Channel_MaxRetransmits
func dll_Channel_MaxRetransmits(channel C.CGO_Channel) C.int {
	r1, _, _ := myfuncs["dll_Channel_MaxRetransmits"].Call(uintptr(channel))
	return C.int(r1)
}

//export dll_Channel_Negotiated
func dll_Channel_Negotiated(channel C.CGO_Channel) C.bool {
	r1, _, _ := myfuncs["dll_Channel_Negotiated"].Call(uintptr(channel))
	res := int8(r1)
	if res == 1 {
		return C.bool(true)
	} else {
		return C.bool(false)
	}
}

//export dll_getFakeDataChannel
func dll_getFakeDataChannel() unsafe.Pointer {
	r1, _, _ := myfuncs["dll_getFakeDataChannel"].Call()
	return unsafe.Pointer(r1)
}

//export dll_Channel_Ordered
func dll_Channel_Ordered(channel C.CGO_Channel) C.bool {
	r1, _, _ := myfuncs["dll_Channel_Ordered"].Call(uintptr(channel))
	res := int8(r1)
	if res == 1 {
		return C.bool(true)
	} else {
		return C.bool(false)
	}
}

//export dll_Channel_ReadyState
func dll_Channel_ReadyState(channel C.CGO_Channel) C.int {
	r1, _, _ := myfuncs["dll_Channel_ReadyState"].Call(uintptr(channel))
	return C.int(r1)
}

//export dll_Channel_Close
func dll_Channel_Close(channel C.CGO_Channel) {
	myfuncs["dll_Channel_Close"].Call(uintptr(channel))
}

//export dll_Channel_Send
func dll_Channel_Send(channel C.CGO_Channel, data unsafe.Pointer, size C.int, binary C.bool) {
	var intBool = 0
	if binary == C.bool(true) {
		intBool = 1
	}
	myfuncs["dll_Channel_Send"].Call(uintptr(channel), uintptr(data), uintptr(size), uintptr(intBool))
}

//export dll_fakeBufferAmount
func dll_fakeBufferAmount(channel C.CGO_Channel, amount C.int) {
	myfuncs["dll_fakeBufferAmount"].Call(uintptr(channel), uintptr(amount))
}

//export dll_fakeMessage
func dll_fakeMessage(channel C.CGO_Channel, data unsafe.Pointer, size C.int) {
	myfuncs["dll_fakeMessage"].Call(uintptr(channel), uintptr(data), uintptr(size))
}

//export dll_fakeStateChange
func dll_fakeStateChange(channel C.CGO_Channel, state C.int) {
	myfuncs["dll_fakeStateChange"].Call(uintptr(channel), uintptr(state))
}
