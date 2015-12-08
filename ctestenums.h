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

#ifdef __cplusplus
}
#endif

#endif
