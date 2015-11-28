// Go wrapper around libwebrtc.
// TODO
package webrtcs
/*

// #cgo CPPFLAGS: -I/usr/include/c++/5.2.0 -I/usr/include/c++/5.2.0/x86_64-unknown-linux-gnu
#cgo CPPFLAGS: -Ithird_party/libwebrtc/
#cgo CXXFLAGS: -std=gnu++11 -Wno-c++0x-extensions
#cgo LDFLAGS: -Llib
#cgo LDFLAGS: -Wl,-z,now -Wl,-z,relro -Wl,--fatal-warnings -Wl,-z,defs -pthread
#cgo LDFLAGS: -Wl,-z,noexecstack -fPIC -fuse-ld=gold
#cgo LDFLAGS: -B/home/serene/code/webrtc-check/src/third_party/binutils/Linux_x64/Release/bin
#cgo LDFLAGS: -Wl,--disable-new-dtags -pthread -m64 -Wl,--detect-odr-violations
// #cgo LDFLAGS: -Wl,--icf=all -Wl,-O1 -Wl,--as-needed -Wl,--gc-sections
#cgo LDFLAGS: -lwebrtc_magic
#cgo LDFLAGS: -lgthread-2.0 -lgtk-x11-2.0 -lgdk-x11-2.0 -lpangocairo-1.0
#cgo LDFLAGS: -latk-1.0 -lcairo -lgdk_pixbuf-2.0 -lgio-2.0 -lpangoft2-1.0 -lpango-1.0
#cgo LDFLAGS: -lgobject-2.0 -lglib-2.0 -lfontconfig -lfreetype -lX11 -lXcomposite
#cgo LDFLAGS: -lXext -lXrender -ldl -lrt -lexpat -lm

#include "cpeerconnection.h"
int lol(int x) {
  return x+10;
}
*/
import "C"

func Lol(x int) int {
  return int(C.lol(C.int(x)))
}

type RTCPeerConnection struct {
  // createOffer
  // createAnswer
  // setLocalDescription
  // localDescription

  // currentLocalDescription
  // pendingLocalDescription
  // setRemoteDescription

  // remoteDescription
  // currentRemoteDescription
  // pendingRemoteDescription

  // addIceCandidate
  // signalingState
  // iceGatheringState
  // iceConnectionState
  // canTrickleIceCandidates
  // getConfiguration
  // setConfiguration
  // close

  // onnegotiationneeded
  // onicecandidate
  // onsignalingstatechange
  // oniceconnectionstatechange
  // onicegatheringstatechange
}

type RTCSignalingState int
/*
const {
  stable RTCSignallingState = iota
  have-local-offer
  have-remote-offer
  have-local-pranswer
  have-remote-pranswer
  closed
}
*/

type RTCConfiguration struct {
}

type RTCIceCredentialType struct {
}

type RTCIceServer struct {
  Urls        string
  Username    string
  Credential  string
  // credentialType   RTCIceCredentialType
}

type RTCIceTransportPolicy struct {
}

type RTCBundlePolicy struct {
}
