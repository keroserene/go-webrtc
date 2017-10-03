#include "peerconnection.h"
#include "datachannel.h"
#include "ctestenums.h"
#include "dllExports.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Callback forward declarations
void(*callbackChannelOnMessagePtr)(int channel, void* data, int size);
void(*callbackChannelOnStateChangePtr)(int channel);
void(*callbackChannelOnBufferedAmountChangePtr)(int channel, int amount);
void(*callbackOnIceCandidatePtr)(int, CGO_IceCandidate*);
void(*callbackOnIceCandidateErrorPtr)(int);
void(*callbackOnSignalingStateChangePtr)(int, int);
void(*callbackOnNegotiationNeededPtr)(int);
void(*callbackOnIceConnectionStateChangePtr)(int p, int state);
void(*callbackOnConnectionStateChangePtr)(int p, int state);
void(*callbackOnIceGatheringStateChangePtr)(int p, int state);
void(*callbackOnDataChannelPtr)(int index, void* o);

extern "C" {
	// Set callback function pointers
	void SetCallbackChannelOnMessage(void(*cb)(int channel, void* data, int size)) {
		callbackChannelOnMessagePtr = cb;
	}

	void SetCallbackChannelOnStateChange(void(*cb)(int channel)) {
		callbackChannelOnStateChangePtr = cb;
	}

	void SetCallbackChannelOnBufferedAmountChange(void(*cb)(int channel, int amount)) {
		callbackChannelOnBufferedAmountChangePtr = cb;
	}

	void SetCallbackOnIceCandidate(void(*cb)(int, CGO_IceCandidate*)) {
		callbackOnIceCandidatePtr = cb;
	}

	void SetCallbackOnIceCandidateError(void(*cb)(int)) {
		callbackOnIceCandidateErrorPtr = cb;
	}

	void SetCallbackOnSignalingStateChange(void(*cb)(int, int)) {
		callbackOnSignalingStateChangePtr = cb;
	}

	void SetCallbackOnNegotiationNeeded(void(*cb)(int)) {
		callbackOnNegotiationNeededPtr = cb;
	}

	void SetCallbackOnIceConnectionStateChange(void(*cb)(int, int)) {
		callbackOnIceConnectionStateChangePtr = cb;
	}

	void SetCallbackOnConnectionStateChange(void(*cb)(int, int)) {
		callbackOnConnectionStateChangePtr = cb;
	}

	void SetCallbackOnIceGatheringStateChange(void(*cb)(int, int)) {
		callbackOnIceGatheringStateChangePtr = cb;
	}

	void SetCallbackOnDataChannel(void(*cb)(int, void*)) {
		callbackOnDataChannelPtr = cb;
	}

	// Enums
	int dll_DataStateConnecting()		{ return CGO_DataStateConnecting; }
	int dll_DataStateOpen()				{ return CGO_DataStateOpen; }
	int dll_DataStateClosing()			{ return CGO_DataStateClosing; }
	int dll_DataStateClosed()			{ return CGO_DataStateClosed; }
	int dll_IceTransportPolicyNone()	{ return CGO_IceTransportPolicyNone; }
	int dll_IceTransportPolicyRelay()	{ return CGO_IceTransportPolicyRelay; }
	int dll_IceTransportPolicyNoHost()	{ return CGO_IceTransportPolicyNoHost; }
	int dll_IceTransportPolicyAll()		{ return CGO_IceTransportPolicyAll; }
	int dll_BundlePolicyBalanced()		{ return CGO_BundlePolicyBalanced; }
	int dll_BundlePolicyMaxBundle()		{ return CGO_BundlePolicyMaxBundle; }
	int dll_BundlePolicyMaxCompat()		{ return CGO_BundlePolicyMaxCompat; }
	int dll_SignalingStateStable()				{ return CGO_SignalingStateStable; }
	int dll_SignalingStateHaveLocalOffer()		{ return CGO_SignalingStateHaveLocalOffer; }
	int dll_SignalingStateHaveLocalPrAnswer()	{ return CGO_SignalingStateHaveLocalPrAnswer; }
	int dll_SignalingStateHaveRemoteOffer()		{ return CGO_SignalingStateHaveRemoteOffer; }
	int dll_SignalingStateHaveRemotePrAnswer()	{ return CGO_SignalingStateHaveRemotePrAnswer; }
	int dll_SignalingStateClosed()				{ return CGO_SignalingStateClosed; }
	int dll_IceConnectionStateNew()				{ return CGO_IceConnectionStateNew; }
	int dll_IceConnectionStateChecking()		{ return CGO_IceConnectionStateChecking; }
	int dll_IceConnectionStateConnected()		{ return CGO_IceConnectionStateConnected; }
	int dll_IceConnectionStateCompleted()		{ return CGO_IceConnectionStateCompleted; }
	int dll_IceConnectionStateFailed()			{ return CGO_IceConnectionStateFailed; }
	int dll_IceConnectionStateDisconnected()	{ return CGO_IceConnectionStateDisconnected; }
	int dll_IceConnectionStateClosed()			{ return CGO_IceConnectionStateClosed; }
	int dll_IceGatheringStateNew()				{ return CGO_IceGatheringStateNew; }
	int dll_IceGatheringStateGathering()		{ return CGO_IceGatheringStateGathering; }
	int dll_IceGatheringStateComplete()			{ return CGO_IceGatheringStateComplete; }

	void* dll_Channel_RegisterObserver(void* o, int index) {
		return CGO_Channel_RegisterObserver(o, index);
	}

	void* dll_getFakeDataChannel() {
		return CGO_getFakeDataChannel();
	}

	const char *dll_Channel_Label(CGO_Channel channel) {
		const char* res = CGO_Channel_Label(channel);
		return res;
	}

	const char *dll_Channel_Protocol(CGO_Channel channel) {
		const char* res = CGO_Channel_Protocol(channel);
		return res;
	}

	bool dll_Channel_Ordered(CGO_Channel channel) {
		return CGO_Channel_Ordered(channel);
	}

	bool dll_Channel_Negotiated(CGO_Channel channel) {
		return CGO_Channel_Negotiated(channel);
	}

	int dll_Channel_MaxRetransmitTime(CGO_Channel channel) {
		return CGO_Channel_MaxRetransmitTime(channel);
	}

	int dll_Channel_MaxRetransmits(CGO_Channel channel) {
		return CGO_Channel_MaxRetransmits(channel);
	}

	int dll_Channel_ID(CGO_Channel channel) {
		return CGO_Channel_ID(channel);
	}

	int dll_Channel_BufferedAmount(CGO_Channel channel) {
		return CGO_Channel_BufferedAmount(channel);
	}

	int dll_Channel_ReadyState(CGO_Channel channel) {
		return CGO_Channel_ReadyState(channel);
	}

	void dll_Channel_Send(CGO_Channel channel, void *data, int size, bool binary) {
		return CGO_Channel_Send(channel, data, size, binary);
	}

	void dll_Channel_Close(CGO_Channel channel) {
		CGO_Channel_Close(channel);
	}

	void dll_fakeMessage(CGO_Channel channel, void *data, int size) {
		return CGO_fakeMessage(channel, data, size);
	}

	void dll_fakeStateChange(CGO_Channel channel, int state) {
		return CGO_fakeStateChange(channel, state);
	}

	void dll_fakeBufferAmount(CGO_Channel channel, int amount) {
		return CGO_fakeBufferAmount(channel, amount);
	}

	void dll_Free(void* ptr) {
		free(ptr); 
	}

	// PeerConnection
	void* dll_InitializePeer(int goPc) {
		return CGO_InitializePeer(goPc);
	}

	int dll_CreatePeerConnection(CGO_Peer peer, CGO_Configuration *config) {
		return CGO_CreatePeerConnection(peer, config);
	}

	CGO_sdpString dll_CreateOffer(CGO_Peer peer) {
		return CGO_CreateOffer(peer);
	}

	CGO_sdpString dll_CreateAnswer(CGO_Peer peer) {
		return CGO_CreateAnswer(peer);
	}

	int dll_SetLocalDescription(CGO_Peer peer, CGO_sdp sdp) {
		return CGO_SetLocalDescription(peer, sdp);
	}

	int dll_SetRemoteDescription(CGO_Peer peer, CGO_sdp sdp) {
		return CGO_SetRemoteDescription(peer, sdp);
	}

	CGO_sdp dll_GetLocalDescription(CGO_Peer peer) {
		return CGO_GetLocalDescription(peer);
	}

	CGO_sdp dll_GetRemoteDescription(CGO_Peer peer) {
		return CGO_GetRemoteDescription(peer);
	}

	int dll_GetSignalingState(CGO_Peer peer) {
		return CGO_GetSignalingState(peer);
	}

	int dll_IceConnectionState(CGO_Peer peer) {
		return CGO_IceConnectionState(peer);
	}

	int dll_IceGatheringState(CGO_Peer peer) {
		return CGO_IceGatheringState(peer);
	}

	int dll_SetConfiguration(CGO_Peer peer, CGO_Configuration *config) {
		return CGO_SetConfiguration(peer, config);
	}

	int dll_AddIceCandidate(CGO_Peer cgoPeer, CGO_IceCandidate *cgoIC) {
		return CGO_AddIceCandidate(cgoPeer, cgoIC);
	}

	void* dll_CreateDataChannel(CGO_Peer peer, char *label, void *dict) {
		return CGO_CreateDataChannel(peer, label, dict);
	}

	CGO_sdpString dll_SerializeSDP(CGO_sdp sdp) {
		return CGO_SerializeSDP(sdp);
	}

	CGO_sdp dll_DeserializeSDP(const char *type, const char *msg) {
		return CGO_DeserializeSDP(type, msg);
	}

	void dll_Close(CGO_Peer peer) {
		CGO_Close(peer);
	}

	void dll_fakeIceCandidateError(CGO_Peer peer) {
		CGO_fakeIceCandidateError(peer);
	}
}

/*
* Callbacks
*/
	void cgoOnIceCandidate(int p, CGO_IceCandidate cIC) {
		callbackOnIceCandidatePtr(p, &cIC);
	}
	void cgoOnIceCandidateError(int p) {
		callbackOnIceCandidateErrorPtr(p);
	}
	void cgoOnSignalingStateChange(int p, int s) {
		callbackOnSignalingStateChangePtr(p, s);
	}
	void cgoOnNegotiationNeeded(int p) {
		callbackOnNegotiationNeededPtr(p);
	}
	void cgoChannelOnStateChange(int channel) {
		callbackChannelOnStateChangePtr(channel);
	}
	void cgoChannelOnMessage(int channel, void* data, int size) {
		callbackChannelOnMessagePtr(channel, data, size);
	}
	void cgoChannelOnBufferedAmountChange(int channel, int amount) {
		callbackChannelOnBufferedAmountChangePtr(channel, amount);
	}
	void cgoOnIceConnectionStateChange(int p, int state) {
		callbackOnIceConnectionStateChangePtr(p, state);
	}
	void cgoOnConnectionStateChange(int p, int state) {
		callbackOnConnectionStateChangePtr(p, state);
	}
	void cgoOnIceGatheringStateChange(int p, int state) {
		callbackOnIceGatheringStateChangePtr(p, state);
	}
	void cgoOnDataChannel(int index, void *o) {
		callbackOnDataChannelPtr(index, o);
	}
