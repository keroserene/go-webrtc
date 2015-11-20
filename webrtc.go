// Go wrapper around libwebrtc.
// TODO
package webrtcs
/*

// #cgo CPPFLAGS: -I/usr/include/c++/5.2.0 -I/usr/include/c++/5.2.0/x86_64-unknown-linux-gnu
#cgo CPPFLAGS: -Ithird_party/libwebrtc/
#cgo CXXFLAGS: -std=gnu++11 -Wno-c++0x-extensions
// #cgo LDFLAGS: -extldflags -static
// #cgo LDFLAGS: -L/home/serene/code/webrtc-check/src/out/Release/lib/
// #cgo LDFLAGS: -ljingle_peerconnection_so -lboringssl -lprotobuf_lite -licuuc
// #cgo LDFLAGS: -L/home/serene/code/webrtc-check/src/out/Release/obj/webrtc/base
// #cgo LDFLAGS: -lrtc_base -lrtc_base_approved -lrtc_base_tests_utils
// #cgo LDFLAGS: /home/serene/code/webrtc-check/src/out/Release/obj/webrtc/base/rtc_base_approved.logging.o
// #cgo LDFLAGS: /home/serene/code/webrtc-check/src/out/Release/obj/webrtc/base/rtc_base_approved.logging.o

// I seem to be using ld -r command to combine the object files, and this solves
// the linker errors. However, there's probably a flag I could switch to build
// a single archive to be included here.

// #cgo LDFLAGS: /home/serene/code/webrtc-check/src/out/Release/obj/webrtc/base/omg_rtc_base.o
// #cgo LDFLAGS: /home/serene/code/webrtc-check/src/out/Release/obj/talk/session/media/libjingle_p2p.a
// #cgo LDFLAGS: /home/serene/code/webrtc-check/src/out/Release/obj/talk/media/webrtc/libjingle_media.a
// #cgo LDFLAGS: /home/serene/code/webrtc-check/src/out/Release/obj/talk/media/sctp/libjingle_media.sctpdataengine.o
// #cgo LDFLAGS: /home/serene/code/webrtc-check/src/out/Release/obj/talk/media/base/libjingle_media.a
// #cgo LDFLAGS: /home/serene/code/webrtc-check/src/out/Release/obj/talk/app/webrtc/libjingle_pc.a
// #cgo LDFLAGS: /home/serene/code/webrtc-check/src/out/Release/obj/chromium/src/third_party/jsoncpp/source/src/lib_json/jsoncpp.json_writer.o
// #cgo LDFLAGS: -L/home/serene/code/webrtc-check/src/out/Release/obj/talk
// #cgo LDFLAGS: -ljingle_peerconnection  -ljingle_media -ljingle_p2p
// #cgo LDFLAGS: -L/home/serene/code/webrtc-check/src/out/Release/obj/chromium/src/third_party/jsoncpp
// #cgo LDFLAGS: -ljsoncpp

// -lwebrtc
#cgo LDFLAGS: -Llib
// #cgo LDFLAGS: -ljingle_p2p
#cgo LDFLAGS: -lrtc_base
// #cgo LDFLAGS: -lwebrtc_common
#cgo LDFLAGS: -lrtc_base_approved
// #cgo LDFLAGS: -ljsoncpp
// #cgo LDFLAGS: -ljingle_media
// #cgo LDFLAGS: -lwebrtc_video_coding
// #cgo LDFLAGS: -ljingle_peerconnection_so
#cgo LDFLAGS: -lboringssl -lprotobuf_lite -licuuc
// #cgo LDFLAGS: -lwebrtc
// #cgo LDFLAGS: -ljingle_peerconnection
// #cgo LDFLAGS: -ljingle_unittest_main
// #cgo LDFLAGS: -lrtc_base_tests_utils

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
