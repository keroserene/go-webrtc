// Go wrapper for libwebrtc.
package webrtc

/*
// #cgo CPPFLAGS: -I/usr/include/c++/5.2.0 -I/usr/include/c++/5.2.0/x86_64-unknown-linux-gnu
#cgo CPPFLAGS: -Ithird_party/libwebrtc/
#cgo CXXFLAGS: -std=gnu++11 -Wno-c++0x-extensions
#cgo LDFLAGS: -Wl,-z,now -Wl,-z,relro -Wl,--fatal-warnings -Wl,-z,defs -pthread
#cgo LDFLAGS: -Wl,-z,noexecstack -fPIC -fuse-ld=gold
#cgo LDFLAGS: -B/home/serene/code/webrtc-check/src/third_party/binutils/Linux_x64/Release/bin
#cgo LDFLAGS: -Wl,--disable-new-dtags -pthread -m64 -Wl,--detect-odr-violations
// #cgo LDFLAGS: -Wl,--icf=all -Wl,-O1 -Wl,--as-needed -Wl,--gc-sections
#cgo LDFLAGS: -L/home/serene/code/go/src/github.com/keroserene/webrtc/lib
// #cgo LDFLAGS: -Llib
// libwebrtc_magic is a custom built archive based on many other ninja files.
#cgo LDFLAGS: -lwebrtc_magic
#cgo LDFLAGS: -lgthread-2.0 -lgtk-x11-2.0 -lgdk-x11-2.0 -lpangocairo-1.0
#cgo LDFLAGS: -latk-1.0 -lcairo -lgdk_pixbuf-2.0 -lgio-2.0 -lpangoft2-1.0 -lpango-1.0
#cgo LDFLAGS: -lgobject-2.0 -lglib-2.0 -lfontconfig -lfreetype -lX11 -lXcomposite
#cgo LDFLAGS: -lXext -lXrender -ldl -lrt -lexpat -lm
#include "cpeerconnection.h"
*/
import "C"
import (
  "unsafe"
  "fmt"
)

type Callback func()

type RTCPeerConnection struct {
  CreateOffer func(Callback, Callback)
  CreateAnswer func(Callback)
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
  OnIceCandidate func()
  // onsignalingstatechange
  // oniceconnectionstatechange
  // onicegatheringstatechange
  pc C.PeerConnection
  IceServers string
}

func createOffer(pc RTCPeerConnection, success Callback, fail Callback) {
  fmt.Println("[go] creating offer...")
  C.CreateOffer(pc.pc, unsafe.Pointer(&c))
}

// func createAnswer(pc RTCPeerConnection, c Callback) {
  // C.CreateAnswer(pc.pc, c)
// }

// Install a handler for receiving ICE Candidates.
// func OnIceCandidate(pc RTCPeerConnection) {
// }

func NewPeerConnection() RTCPeerConnection {
  var ret RTCPeerConnection
  ret.pc = C.NewPeerConnection()
  // ret.IceServers = C.GetIceServers(ret.pc)
  // Assign "methods"
  ret.CreateOffer = func(c Callback) {
    createOffer(ret, c)
  }
  // ret.CreateAnswer = func(c Callback) {
    // createAnswer(ret, c)
  // }
  return ret
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

// func main() {
// }
