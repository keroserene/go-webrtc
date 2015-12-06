#ifndef _C_TEST_ENUMS_H_
#define _C_TEST_ENUMS_H_

#define WEBRTC_POSIX 1

#ifdef __cplusplus
extern "C" {
#endif

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

#ifdef __cplusplus
}
#endif

#endif
