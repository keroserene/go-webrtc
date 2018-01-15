#ifndef _C_REFPTR_H_
#define _C_REFPTR_H_

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void* CGO_RefPtr; // rtc::RefCountInterface*

  void CGO_RefPtr_Release(CGO_RefPtr);

#ifdef __cplusplus
}
#endif

#endif  // _C_REFPTR_H_
