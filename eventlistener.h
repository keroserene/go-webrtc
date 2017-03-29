#ifndef _C_EVENTLISTENER_H
#define _C_EVENTLISTENER_H

#define WEBRTC_POSIX 1

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void* CGO_Notifier; // webrtc::NotifierInterface*
  typedef void* CGO_Observer; // webrtc::ObserverInterface*
  typedef int CGO_EventCallback; // key into Go eventCallbacks

  CGO_Observer CGO_NewObserver(CGO_Notifier, CGO_EventCallback);
  void CGO_DeleteObserver(CGO_Observer);

#ifdef __cplusplus
}
#endif

#endif  // _C_EVENTLISTENER_H
