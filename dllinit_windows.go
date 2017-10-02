//+build windows
package webrtc

/*
#include "datachannel.h"
#include "peerconnection.h"
#include "string.h"
*/
import "C"
import (
	"sync"
	"syscall"
	"unsafe"
)

var once sync.Once
var dllHandle *syscall.DLL
var myfuncs map[string]*syscall.Proc

func LoadMyDllOnce() {
	once.Do(loadMyDll)
}

func loadMyDll() {
	var err error
	dllHandle, err = syscall.LoadDLL("webrtc.dll")
	if err != nil {
		panic(err)
	}
	myfuncs = make(map[string]*syscall.Proc)
	funcs := []string{
		// Enums
		"dll_DataStateConnecting",
		"dll_DataStateOpen",
		"dll_DataStateClosing",
		"dll_DataStateClosed",
		"dll_IceTransportPolicyNone",
		"dll_IceTransportPolicyRelay",
		"dll_IceTransportPolicyNoHost",
		"dll_IceTransportPolicyAll",
		"dll_BundlePolicyBalanced",
		"dll_BundlePolicyMaxCompat",
		"dll_BundlePolicyMaxBundle",
		"dll_SignalingStateStable",
		"dll_SignalingStateHaveLocalOffer",
		"dll_SignalingStateHaveLocalPrAnswer",
		"dll_SignalingStateHaveRemoteOffer",
		"dll_SignalingStateHaveRemotePrAnswer",
		"dll_SignalingStateClosed",
		"dll_IceConnectionStateNew",
		"dll_IceConnectionStateChecking",
		"dll_IceConnectionStateConnected",
		"dll_IceConnectionStateCompleted",
		"dll_IceConnectionStateFailed",
		"dll_IceConnectionStateDisconnected",
		"dll_IceConnectionStateClosed",
		"dll_IceGatheringStateNew",
		"dll_IceGatheringStateGathering",
		"dll_IceGatheringStateComplete",

		// DataChannel
		"dll_Channel_ReadyState",
		"dll_Channel_RegisterObserver",
		"dll_getFakeDataChannel",
		"dll_Channel_Label",
		"dll_Channel_Ordered",
		"dll_Channel_Negotiated",
		"dll_Channel_Protocol",
		"dll_Channel_BufferedAmount",
		"dll_Channel_MaxRetransmitTime",
		"dll_Channel_MaxRetransmits",
		"dll_Channel_ID",
		"dll_Channel_Send",
		"dll_Channel_Close",
		"dll_fakeMessage",
		"dll_fakeStateChange",
		"dll_fakeBufferAmount",
		"SetCallbackChannelOnMessage",
		"SetCallbackChannelOnStateChange",
		"SetCallbackChannelOnBufferedAmountChange",
		"SetCallbackOnIceCandidate",
		"SetCallbackOnIceCandidateError",
		"SetCallbackOnSignalingStateChange",
		"SetCallbackOnNegotiationNeeded",
		"SetCallbackOnIceConnectionStateChange",
		"SetCallbackOnConnectionStateChange",
		"SetCallbackOnIceGatheringStateChange",
		"SetCallbackOnDataChannel",

		// PeerConnection
		"dll_InitializePeer",
		"dll_SetConfiguration",
		"dll_Close",
		"dll_CreatePeerConnection",
		"dll_CreateOffer",
		"dll_CreateAnswer",
		"dll_SetLocalDescription",
		"dll_SetRemoteDescription",
		"dll_GetLocalDescription",
		"dll_GetRemoteDescription",
		"dll_GetSignalingState",
		"dll_IceConnectionState",
		"dll_IceGatheringState",
		"dll_AddIceCandidate",
		"dll_DeserializeSDP",
		"dll_SerializeSDP",
		"dll_CreateDataChannel",
		"dll_fakeIceCandidateError",

		// Misc
		"dll_Free",
	}
	for _, v := range funcs {
		proc, err := dllHandle.FindProc(v)
		if err != nil {
			panic(err)
		}
		myfuncs[v] = proc
	}

	// Initialize the enums
	_cgoIceTransportPolicyNone = dll_LoadInt("dll_IceTransportPolicyNone")
	_cgoIceTransportPolicyRelay = dll_LoadInt("dll_IceTransportPolicyRelay")
	_cgoIceTransportPolicyNoHost = dll_LoadInt("dll_IceTransportPolicyNoHost")
	_cgoIceTransportPolicyAll = dll_LoadInt("dll_IceTransportPolicyAll")
	_cgoBundlePolicyBalanced = dll_LoadInt("dll_BundlePolicyBalanced")
	_cgoBundlePolicyMaxCompat = dll_LoadInt("dll_BundlePolicyMaxCompat")
	_cgoBundlePolicyMaxBundle = dll_LoadInt("dll_BundlePolicyMaxBundle")
	_cgoSignalingStateStable = dll_LoadInt("dll_SignalingStateStable")
	_cgoSignalingStateHaveLocalOffer = dll_LoadInt("dll_SignalingStateHaveLocalOffer")
	_cgoSignalingStateHaveLocalPrAnswer = dll_LoadInt("dll_SignalingStateHaveLocalPrAnswer")
	_cgoSignalingStateHaveRemoteOffer = dll_LoadInt("dll_SignalingStateHaveRemoteOffer")
	_cgoSignalingStateHaveRemotePrAnswer = dll_LoadInt("dll_SignalingStateHaveRemotePrAnswer")
	_cgoSignalingStateClosed = dll_LoadInt("dll_SignalingStateClosed")

	_cgoIceConnectionStateNew = dll_LoadInt("dll_IceConnectionStateNew")
	_cgoIceConnectionStateChecking = dll_LoadInt("dll_IceConnectionStateChecking")
	_cgoIceConnectionStateConnected = dll_LoadInt("dll_IceConnectionStateConnected")
	_cgoIceConnectionStateCompleted = dll_LoadInt("dll_IceConnectionStateCompleted")
	_cgoIceConnectionStateFailed = dll_LoadInt("dll_IceConnectionStateFailed")
	_cgoIceConnectionStateDisconnected = dll_LoadInt("dll_IceConnectionStateDisconnected")
	_cgoIceConnectionStateClosed = dll_LoadInt("dll_IceConnectionStateClosed")
	_cgoIceGatheringStateNew = dll_LoadInt("dll_IceGatheringStateNew")
	_cgoIceGatheringStateGathering = dll_LoadInt("dll_IceGatheringStateGathering")
	_cgoIceGatheringStateComplete = dll_LoadInt("dll_IceGatheringStateComplete")

	C.CGO_DataChannelInit()

	// Setup the callbacks
	cb := syscall.NewCallback(func(goChannel int, data unsafe.Pointer, size int) uintptr {
		cgoChannelOnMessage(goChannel, data, size)
		return 0
	})
	myfuncs["SetCallbackChannelOnMessage"].Call(cb)
	cb = syscall.NewCallback(func(goChannel int) uintptr {
		cgoChannelOnStateChange(goChannel)
		return 0
	})
	myfuncs["SetCallbackChannelOnStateChange"].Call(cb)
	cb = syscall.NewCallback(func(goChannel int, amount int) uintptr {
		cgoChannelOnBufferedAmountChange(goChannel, amount)
		return 0
	})
	myfuncs["SetCallbackChannelOnBufferedAmountChange"].Call(cb)

	// We cannot pass a struct to NewCallback(), so we must use a pointer
	cb = syscall.NewCallback(func(p int, cIC *C.CGO_IceCandidate) uintptr {
		cgoOnIceCandidate(p, *cIC)
		return 0
	})
	myfuncs["SetCallbackOnIceCandidate"].Call(cb)

	cb = syscall.NewCallback(func(p int) uintptr {
		cgoOnIceCandidateError(p)
		return 0
	})
	myfuncs["SetCallbackOnIceCandidateError"].Call(cb)

	cb = syscall.NewCallback(func(p int, s SignalingState) uintptr {
		cgoOnSignalingStateChange(p, s)
		return 0
	})
	myfuncs["SetCallbackOnSignalingStateChange"].Call(cb)

	cb = syscall.NewCallback(func(p int) uintptr {
		cgoOnNegotiationNeeded(p)
		return 0
	})
	myfuncs["SetCallbackOnNegotiationNeeded"].Call(cb)

	cb = syscall.NewCallback(func(p int, s IceConnectionState) uintptr {
		cgoOnIceConnectionStateChange(p, s)
		return 0
	})
	myfuncs["SetCallbackOnIceConnectionStateChange"].Call(cb)

	cb = syscall.NewCallback(func(p int, s IceConnectionState) uintptr {
		cgoOnConnectionStateChange(p, s)
		return 0
	})
	myfuncs["SetCallbackOnConnectionStateChange"].Call(cb)

	cb = syscall.NewCallback(func(p int, s IceGatheringState) uintptr {
		cgoOnIceGatheringStateChange(p, s)
		return 0
	})
	myfuncs["SetCallbackOnIceGatheringStateChange"].Call(cb)

	cb = syscall.NewCallback(func(p int, o unsafe.Pointer) uintptr {
		cgoOnDataChannel(p, o)
		return 0
	})
	myfuncs["SetCallbackOnDataChannel"].Call(cb)

}

func init() {
	LoadMyDllOnce()
}

//*********** Enums are all initialized here
//export dll_DataStateConnecting
func dll_DataStateConnecting() C.int {
	r1, _, _ := myfuncs["dll_DataStateConnecting"].Call()
	// Variables (consts) need to be set explicitly here
	_cgoDataStateConnecting = int(r1)
	return C.int(r1)
}

//export dll_DataStateOpen
func dll_DataStateOpen() C.int {
	r1, _, _ := myfuncs["dll_DataStateOpen"].Call()
	// Variables (consts) need to be set explicitly here
	_cgoDataStateOpen = int(r1)
	return C.int(r1)
}

//export dll_DataStateClosing
func dll_DataStateClosing() C.int {
	r1, _, _ := myfuncs["dll_DataStateClosing"].Call()
	// Variables (consts) need to be set explicitly here
	_cgoDataStateClosing = int(r1)
	return C.int(r1)
}

//export dll_DataStateClosed
func dll_DataStateClosed() C.int {
	r1, _, _ := myfuncs["dll_DataStateClosed"].Call()
	// Variables (consts) need to be set explicitly here
	_cgoDataStateClosed = int(r1)
	return C.int(r1)
}

func dll_LoadInt(procName string) int {
	r1, _, _ := myfuncs[procName].Call()
	return int(r1)
}
