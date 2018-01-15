#ifndef _C_MEDIASTREAM_H_
#define _C_MEDIASTREAM_H_

#define WEBRTC_POSIX 1

#include "audiotrack.h"

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void* CGO_MediaStream; // webrtc::MediaStreamInterface*

  CGO_AudioTrack* CGO_MediaStream_GetAudioTracks(CGO_MediaStream, int*);
  void CGO_MediaStream_AddAudioTrack(CGO_MediaStream, CGO_AudioTrack);

#ifdef __cplusplus
}
#endif

#endif  // _C_MEDIASTREAM_H_
