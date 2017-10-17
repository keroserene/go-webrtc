#ifndef _C_RTPRECEIVER_H_
#define _C_RTPRECEIVER_H_

#define WEBRTC_POSIX 1

#include <stdbool.h>

#include "mediastreamtrack.h"

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void* CGO_RtpReceiver; // webrtc::RtpReceiverInterface*

  CGO_MediaStreamTrack CGO_RtpReceiver_Track(CGO_RtpReceiver, bool*);

#ifdef __cplusplus
}
#endif

#endif  // _C_RTPRECEIVER_H_
