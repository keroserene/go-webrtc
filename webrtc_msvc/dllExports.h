#include <stdlib.h>

extern "C" {

	// Enums
	__declspec(dllexport) int dll_DataStateConnecting();
	__declspec(dllexport) int dll_DataStateOpen();
	__declspec(dllexport) int dll_DataStateClosing();
	__declspec(dllexport) int dll_DataStateClosed();
	__declspec(dllexport) int dll_IceTransportPolicyNone();
	__declspec(dllexport) int dll_IceTransportPolicyRelay();
	__declspec(dllexport) int dll_IceTransportPolicyNoHost();
	__declspec(dllexport) int dll_IceTransportPolicyAll();
	__declspec(dllexport) int dll_BundlePolicyBalanced();
	__declspec(dllexport) int dll_BundlePolicyMaxBundle();
	__declspec(dllexport) int dll_BundlePolicyMaxCompat();
	__declspec(dllexport) int dll_SignalingStateStable();
	__declspec(dllexport) int dll_SignalingStateHaveLocalOffer();
	__declspec(dllexport) int dll_SignalingStateHaveLocalPrAnswer();
	__declspec(dllexport) int dll_SignalingStateHaveRemoteOffer();
	__declspec(dllexport) int dll_SignalingStateHaveRemotePrAnswer();
	__declspec(dllexport) int dll_SignalingStateClosed();
	__declspec(dllexport) int dll_IceConnectionStateNew();
	__declspec(dllexport) int dll_IceConnectionStateChecking();
	__declspec(dllexport) int dll_IceConnectionStateConnected();
	__declspec(dllexport) int dll_IceConnectionStateCompleted();
	__declspec(dllexport) int dll_IceConnectionStateFailed();
	__declspec(dllexport) int dll_IceConnectionStateDisconnected();
	__declspec(dllexport) int dll_IceConnectionStateClosed();
	__declspec(dllexport) int dll_IceGatheringStateNew();
	__declspec(dllexport) int dll_IceGatheringStateGathering();
	__declspec(dllexport) int dll_IceGatheringStateComplete();

	// DataChannel
	__declspec(dllexport) void* dll_Channel_RegisterObserver(void* o, int index);
	__declspec(dllexport) void* dll_getFakeDataChannel();
	__declspec(dllexport) const char *dll_Channel_Label(CGO_Channel channel);
	__declspec(dllexport) const char *dll_Channel_Protocol(CGO_Channel channel);
	__declspec(dllexport) int dll_Channel_ID(CGO_Channel channel);
	__declspec(dllexport) int dll_Channel_BufferedAmount(CGO_Channel channel);
	__declspec(dllexport) int dll_Channel_ReadyState(CGO_Channel channel);
	__declspec(dllexport) bool dll_Channel_Ordered(CGO_Channel channel);
	__declspec(dllexport) bool dll_Channel_Negotiated(CGO_Channel channel);
	__declspec(dllexport) int dll_Channel_MaxRetransmitTime(CGO_Channel channel);
	__declspec(dllexport) int dll_Channel_MaxRetransmits(CGO_Channel channel);
	__declspec(dllexport) void dll_Channel_Send(CGO_Channel channel, void *data, int size, bool binary);
	__declspec(dllexport) void dll_Channel_Close(CGO_Channel channel);
	__declspec(dllexport) void dll_fakeMessage(CGO_Channel channel, void *data, int size);
	__declspec(dllexport) void dll_fakeStateChange(CGO_Channel channel, int state);
	__declspec(dllexport) void dll_fakeBufferAmount(CGO_Channel channel, int amount);

	// PeerConnection
	__declspec(dllexport) void* dll_InitializePeer(int goPc);
	__declspec(dllexport) void dll_Close(CGO_Peer peer);
	__declspec(dllexport) int dll_CreatePeerConnection(CGO_Peer, CGO_Configuration *);
	__declspec(dllexport) CGO_sdpString dll_CreateOffer(CGO_Peer);
	__declspec(dllexport) CGO_sdpString dll_CreateAnswer(CGO_Peer);
	__declspec(dllexport) int dll_SetLocalDescription(CGO_Peer, CGO_sdp);
	__declspec(dllexport) int dll_SetRemoteDescription(CGO_Peer, CGO_sdp);
	__declspec(dllexport) CGO_sdp dll_GetLocalDescription(CGO_Peer);
	__declspec(dllexport) CGO_sdp dll_GetRemoteDescription(CGO_Peer);
	__declspec(dllexport) int dll_GetSignalingState(CGO_Peer);
	__declspec(dllexport) int dll_IceConnectionState(CGO_Peer);
	__declspec(dllexport) int dll_IceGatheringState(CGO_Peer);
	__declspec(dllexport) int dll_SetConfiguration(CGO_Peer, CGO_Configuration *);
	__declspec(dllexport) int dll_AddIceCandidate(CGO_Peer cgoPeer, CGO_IceCandidate *cgoIC);
	__declspec(dllexport) void* dll_CreateDataChannel(CGO_Peer, char *, void *);
	__declspec(dllexport) CGO_sdpString dll_SerializeSDP(CGO_sdp);
	__declspec(dllexport) CGO_sdp dll_DeserializeSDP(const char *type, const char *msg);
	__declspec(dllexport) void dll_fakeIceCandidateError(CGO_Peer peer);

	// Callbacks
	__declspec(dllexport) void SetCallbackChannelOnMessage(void(*cb)(int channel, void* data, int size));
	__declspec(dllexport) void SetCallbackChannelOnStateChange(void(*cb)(int channel));
	__declspec(dllexport) void SetCallbackChannelOnBufferedAmountChange(void(*cb)(int channel, int amount));
	__declspec(dllexport) void SetCallbackOnIceCandidate(void(*cb)(int, CGO_IceCandidate*));
	__declspec(dllexport) void SetCallbackOnIceCandidateError(void(*cb)(int));
	__declspec(dllexport) void SetCallbackOnSignalingStateChange(void(*cb)(int, int));
	__declspec(dllexport) void SetCallbackOnNegotiationNeeded(void(*cb)(int));
	__declspec(dllexport) void SetCallbackOnIceConnectionStateChange(void(*cb)(int, int));
	__declspec(dllexport) void SetCallbackOnConnectionStateChange(void(*cb)(int, int));
	__declspec(dllexport) void SetCallbackOnIceGatheringStateChange(void(*cb)(int, int));
	__declspec(dllexport) void SetCallbackOnDataChannel(void(*cb)(int, void*));

	// Misc
	__declspec(dllexport) void dll_Free(void* ptr);
}

