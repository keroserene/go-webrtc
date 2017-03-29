#ifndef _C_MEDIASTREAMTRACK_H_
#define _C_MEDIASTREAMTRACK_H_

#define WEBRTC_POSIX 1

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void* CGO_MediaStreamTrack; // webrtc::MediaStreamTrackInterface*

  const char* CGO_MediaStreamTrack_ID(CGO_MediaStreamTrack);
  bool CGO_MediaStreamTrack_Enabled(CGO_MediaStreamTrack);
  void CGO_MediaStreamTrack_SetEnabled(CGO_MediaStreamTrack, bool);
  bool CGO_MediaStreamTrack_Ended(CGO_MediaStreamTrack);

#ifdef __cplusplus
}
#endif

#endif  // _C_MEDIASTREAMTRACK_H_
