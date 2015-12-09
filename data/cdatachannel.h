#ifndef _C_DATACHANNEL_H_
#define _C_DATACHANNEL_H_

#define WEBRTC_POSIX 1

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void* CGO_Channel;

  const char *CGO_Channel_Label(CGO_Channel);

#ifdef __cplusplus
}
#endif

#endif  // _C_DATACHANNEL_H
