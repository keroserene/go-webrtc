#ifndef _C_TEST_ENUMS_H_
#define _C_TEST_ENUMS_H_

#define WEBRTC_POSIX 1

#ifdef __cplusplus
extern "C" {
#endif

  // See talk/app/webrtc/peerconnectioninterface.h

  extern const int CGO_IceTransportPolicyNone;
  extern const int CGO_IceTransportPolicyRelay;
  extern const int CGO_IceTransportPolicyNoHost;
  extern const int CGO_IceTransportPolicyAll;

  extern const int CGO_BundlePolicyBalanced;
  extern const int CGO_BundlePolicyMaxBundle;
  extern const int CGO_BundlePolicyMaxCompat;

  // TODO: [ED]
  // extern const int CGO_RtcpMuxPolicyNegotiate;
  // extern const int CGO_RtcpMuxPolicyRequire;

  extern const int CGO_SignalingStateStable;
  extern const int CGO_SignalingStateHaveLocalOffer;
  extern const int CGO_SignalingStateHaveLocalPrAnswer;
  extern const int CGO_SignalingStateHaveRemoteOffer;
  extern const int CGO_SignalingStateHaveRemotePrAnswer;
  extern const int CGO_SignalingStateClosed;

  extern const int CGO_IceConnectionStateNew;
  extern const int CGO_IceConnectionStateChecking;
  extern const int CGO_IceConnectionStateConnected;
  extern const int CGO_IceConnectionStateCompleted;
  extern const int CGO_IceConnectionStateFailed;
  extern const int CGO_IceConnectionStateDisconnected;
  extern const int CGO_IceConnectionStateClosed;

  extern const int CGO_IceGatheringStateNew;
  extern const int CGO_IceGatheringStateGathering;
  extern const int CGO_IceGatheringStateComplete;

#ifdef __cplusplus
}
#endif

#endif
