#ifndef _C_DATACHANNEL_H_
#define _C_DATACHANNEL_H_

#define WEBRTC_POSIX 1

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void* CGO_Channel;

  const char *CGO_Channel_Label(CGO_Channel);
  const bool CGO_Channel_Ordered(CGO_Channel);
  int CGO_Channel_ReadyState(CGO_Channel);

  extern const int CGO_DataStateConnecting;
  extern const int CGO_DataStateOpen;
  extern const int CGO_DataStateClosing;
  extern const int CGO_DataStateClosed;

#ifdef __cplusplus
}
#endif

#endif  // _C_DATACHANNEL_H
