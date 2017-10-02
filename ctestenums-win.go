// +build windows
package webrtc

/*
#ifdef _WIN32
int CGO_IceTransportPolicyNone = 0;
int CGO_IceTransportPolicyRelay = 0;
int CGO_IceTransportPolicyNoHost = 0;
int CGO_IceTransportPolicyAll = 0;
int CGO_BundlePolicyBalanced = 0;
int CGO_BundlePolicyMaxBundle = 0;
int CGO_BundlePolicyMaxCompat = 0;
int CGO_SignalingStateStable = 0;
int CGO_SignalingStateHaveLocalOffer = 0;
int CGO_SignalingStateHaveLocalPrAnswer = 0;
int CGO_SignalingStateHaveRemoteOffer = 0;
int CGO_SignalingStateHaveRemotePrAnswer = 0;
int CGO_SignalingStateClosed = 0;
int CGO_IceConnectionStateNew = 0;
int CGO_IceConnectionStateChecking = 0;
int CGO_IceConnectionStateConnected = 0;
int CGO_IceConnectionStateCompleted = 0;
int CGO_IceConnectionStateFailed = 0;
int CGO_IceConnectionStateDisconnected = 0;
int CGO_IceConnectionStateClosed = 0;
int CGO_IceGatheringStateNew = 0;
int CGO_IceGatheringStateGathering = 0;
int CGO_IceGatheringStateComplete = 0;
#endif
*/
import "C"
