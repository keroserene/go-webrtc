#ifndef _C_TEST_ENUMS_H_
#define _C_TEST_ENUMS_H_

#define WEBRTC_POSIX 1

#ifdef __cplusplus
extern "C" {
#endif

  // See talk/app/webrtc/peerconnectioninterface.h

  extern const int CGOIceTransportPolicyNone;
  extern const int CGOIceTransportPolicyRelay;
  extern const int CGOIceTransportPolicyNoHost;
  extern const int CGOIceTransportPolicyAll;

  extern const int CGOBundlePolicyBalanced;
  extern const int CGOBundlePolicyMaxBundle;
  extern const int CGOBundlePolicyMaxCompat;

  // TODO: [ED]
  // extern const int CGORtcpMuxPolicyNegotiate;
  // extern const int CGORtcpMuxPolicyRequire;

  extern const int CGOSignalingStateStable;
  extern const int CGOSignalingStateHaveLocalOffer;
  extern const int CGOSignalingStateHaveLocalPrAnswer;
  extern const int CGOSignalingStateHaveRemoteOffer;
  extern const int CGOSignalingStateHaveRemotePrAnswer;
  extern const int CGOSignalingStateClosed;

#ifdef __cplusplus
}
#endif

#endif
