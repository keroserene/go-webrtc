/*
Package webrtc is a golang wrapper on native code WebRTC.

To provide an easier experience for users of this package, there are differences
inherent in the interface written here and the original native code WebRTC. This
allows users to use WebRTC in a more idiomatic golang way. For example, callback
mechanism has a layer of indirection that allows goroutines instead.

The interface here is based mostly on: w3c.github.io/webrtc-pc

There is also a complication in building the dependent static library for this
to work. Furthermore it is possible that this will break on future versions
of libwebrtc, because the interface with the native code is be fragile.

TODO(keroserene): More package documentation, and more documentation in general.
*/
package webrtc

/*
#cgo CPPFLAGS: -Iinclude/
#cgo CXXFLAGS: -std=c++0x
// -Wno-c++0x-extensions
#cgo LDFLAGS: -Wl,-z,now -Wl,-z,relro -Wl,--fatal-warnings -Wl,-z,defs -pthread
#cgo LDFLAGS: -Wl,-z,noexecstack -fPIC -fuse-ld=gold
#cgo LDFLAGS: -Bbin
#cgo LDFLAGS: -Wl,--disable-new-dtags -pthread -m64 -Wl,--detect-odr-violations
// #cgo LDFLAGS: -Wl,--icf=all -Wl,-O1 -Wl,--as-needed -Wl,--gc-sections
#cgo LDFLAGS: -Llib
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
	"github.com/keroserene/go-webrtc/datachannel"
	"errors"
	"fmt"
	"unsafe"
)

type SDPHeader struct {
	// Keep track of both a pointer to the C++ SessionDescription object,
	// and the serialized string version (which native code generates)
	cgoSdp      C.CGOsdp
	description string
}

type PeerConnection struct {

	localDescription				*SDPHeader
	// currentLocalDescription
	// pendingLocalDescription

	remoteDescription				*SDPHeader
	// currentRemoteDescription
	// pendingRemoteDescription

	// addIceCandidate func()
	// signalingState  RTCSignalingState
	// iceGatheringState  RTCIceGatheringState
	// iceConnectionState  RTCIceConnectionState
	canTrickleIceCandidates  bool
	// getConfiguration
	// setConfiguration
	// close
	OnIceCandidate func()

	// Event handlers:
	// onnegotiationneeded
	// onicecandidate
	// onicecandidateerror
	// onsignalingstatechange
	// onicegatheringstatechange
	// oniceconnectionstatechange

	cgoPeer C.CGOPeer // Native code internals
}

// PeerConnection constructor.
func NewPeerConnection(config *RTCConfiguration) (*PeerConnection, error) {
	pc := new(PeerConnection)
	pc.cgoPeer = C.CGOInitializePeer() // internal CGO Peer.
	if nil == pc.cgoPeer {
		return pc, errors.New("PeerConnection: failed to initialize.")
	}
	cConfig := config.CGO() // Convert for CGO
	if 0 != C.CGOCreatePeerConnection(pc.cgoPeer, &cConfig) {
		return nil, errors.New("PeerConnection: could not create from config.")
	}
	fmt.Println("Created PeerConnection: ", pc)
	return pc, nil
}

// CreateOffer prepares an SDP "offer" message, which should be sent to the target
// peer over a signalling channel.
func (pc *PeerConnection) CreateOffer() (*SDPHeader, error) {
	sdp := C.CGOCreateOffer(pc.cgoPeer)
	if nil == sdp {
		return nil, errors.New("CreateOffer: could not prepare SDP offer.")
	}
	offer := new(SDPHeader)
	offer.cgoSdp = sdp
	offer.description = C.GoString(C.CGOSerializeSDP(sdp))
	return offer, nil
}

func (pc *PeerConnection) SetLocalDescription(sdp *SDPHeader) error {
	r := C.CGOSetLocalDescription(pc.cgoPeer, sdp.cgoSdp)
	if 0 != r {
		return errors.New("SetLocalDescription failed.")
	}
	pc.localDescription = sdp
	return nil
}

// readonly localDescription
func (pc *PeerConnection) LocalDescription() (sdp *SDPHeader) {
	return pc.localDescription
}

func (pc *PeerConnection) SetRemoteDescription(sdp *SDPHeader) error {
	r := C.CGOSetRemoteDescription(pc.cgoPeer, sdp.cgoSdp)
	if 0 != r {
		return errors.New("SetRemoteDescription failed.")
	}
	pc.remoteDescription = sdp
	return nil
}

// readonly remoteDescription
func (pc *PeerConnection) RemoteDescription() (sdp *SDPHeader) {
	return pc.remoteDescription
}

// CreateAnswer prepares an SDP "answer" message, which should be sent in
// response to a peer that has sent an offer, over the signalling channel.
func (pc *PeerConnection) CreateAnswer() (*SDPHeader, error) {
	sdp := C.CGOCreateAnswer(pc.cgoPeer)
	if nil == sdp {
		return nil, errors.New("CreateAnswer failed: could not prepare SDP offer.")
	}
	answer := new(SDPHeader)
	answer.cgoSdp = sdp
	answer.description = C.GoString(C.CGOSerializeSDP(sdp))
	return answer, nil
}

// TODO: Above methods blocks until success or failure occurs. Maybe there should
// actually be a callback version, so the user doesn't have to make their own
// goroutine.

func (pc *PeerConnection) CreateDataChannel(label string, dict datachannel.Init) (
	*datachannel.DataChannel, error) {
	cDC := C.CGOCreateDataChannel(pc.cgoPeer, C.CString(label), unsafe.Pointer(&dict))
	if nil == cDC {
		return nil, errors.New("Failed to CreateDataChannel")
	}
	dc := datachannel.New()
	// dc.cgoDataChannel = cDC
	return dc, nil
}

// Install a handler for receiving ICE Candidates.
// func OnIceCandidate(pc PeerConnection) {
// }
