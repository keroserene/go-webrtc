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

  CGO_Channel CGO_Channel_RegisterObserver(void *obs, int goChannel);

  void CGO_Channel_Send(CGO_Channel channel, void *data, int size, bool binary);
  void CGO_Channel_Close(CGO_Channel channel);

  const char *CGO_Channel_Label(CGO_Channel);
  const bool CGO_Channel_Ordered(CGO_Channel);
  int CGO_Channel_MaxRetransmitTime(CGO_Channel channel);
  int CGO_Channel_MaxRetransmits(CGO_Channel channel);
  const char *CGO_Channel_Protocol(CGO_Channel);
  const bool CGO_Channel_Negotiated(CGO_Channel channel);
  int CGO_Channel_ID(CGO_Channel channel);
  int CGO_Channel_ReadyState(CGO_Channel);
  int CGO_Channel_BufferedAmount(CGO_Channel channel);

  extern const int CGO_DataStateConnecting;
  extern const int CGO_DataStateOpen;
  extern const int CGO_DataStateClosing;
  extern const int CGO_DataStateClosed;

  // Testing helpers:
  CGO_Channel CGO_getFakeDataChannel();
  void CGO_fakeMessage(CGO_Channel channel, void *data, int size);
  void CGO_fakeStateChange(CGO_Channel channel, int state);
  void CGO_fakeBufferAmount(CGO_Channel channel, int amount);

#ifdef __cplusplus
}
#endif

#endif  // _C_DATACHANNEL_H
