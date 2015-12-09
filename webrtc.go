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
#cgo CXXFLAGS: -std=c++0x
#cgo linux,amd64 pkg-config: webrtc-linux-amd64.pc
#cgo darwin,amd64 pkg-config: webrtc-darwin-amd64.pc
#include "cpeerconnection.h"
*/
import "C"
import (
	"errors"
	// "fmt"
	"github.com/keroserene/go-webrtc/data"
	"unsafe"
	// "io"
	"io/ioutil"
	"log"
	"os"
)

var (
	INFO  log.Logger
	WARN  log.Logger
	ERROR log.Logger
	TRACE log.Logger
)

// Logging verbosity level, from 0 (nothing) upwards.
func SetVerbosity(level int) {
	// handle io.Writer
	infoOut := ioutil.Discard
	warnOut := ioutil.Discard
	errOut := ioutil.Discard
	traceOut := ioutil.Discard

	// TODO: Better logging levels
	if level > 0 {
		errOut = os.Stdout
	}
	if level > 1 {
		warnOut = os.Stdout
	}
	if level > 2 {
		infoOut = os.Stdout
	}
	if level > 3 {
		traceOut = os.Stdout
	}

	INFO = *log.New(infoOut,
		"INFO: ",
		// log.Ldate|log.Ltime|log.Lshortfile)
		log.Lshortfile)
	WARN = *log.New(warnOut,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	ERROR = *log.New(errOut,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	TRACE = *log.New(traceOut,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func init() {
	SetVerbosity(3)
}

type SDPHeader struct {
	// Keep track of both a pointer to the C++ SessionDescription object,
	// and the serialized string version (which native code generates)
	cgoSdp      C.CGO_sdp
	description string
}

type PeerConnection struct {
	localDescription *SDPHeader
	// currentLocalDescription
	// pendingLocalDescription

	remoteDescription *SDPHeader
	// currentRemoteDescription
	// pendingRemoteDescription

	// iceGatheringState  RTCIceGatheringState
	// iceConnectionState  RTCIceConnectionState
	canTrickleIceCandidates bool
	// close

	// Event handlers
	OnIceCandidate      func(string)
	OnNegotiationNeeded func()
	// onicecandidateerror
	OnSignalingStateChange func(SignalingState)
	// onicegatheringstatechange
	// oniceconnectionstatechange
	OnDataChannel func(*data.Channel)

	config Configuration

	cgoPeer C.CGO_Peer // Native code internals
}

// PeerConnection constructor.
func NewPeerConnection(config *Configuration) (*PeerConnection, error) {
	pc := new(PeerConnection)
	INFO.Println("PC at ", unsafe.Pointer(pc))
	pc.cgoPeer = C.CGO_InitializePeer(unsafe.Pointer(pc)) // internal CGO_ Peer.
	if nil == pc.cgoPeer {
		return pc, errors.New("PeerConnection: failed to initialize.")
	}
	pc.config = *config
	cConfig := config._CGO() // Convert for CGO_
	if 0 != C.CGO_CreatePeerConnection(pc.cgoPeer, &cConfig) {
		return nil, errors.New("PeerConnection: could not create from config.")
	}
	INFO.Println("Created PeerConnection: ", pc, pc.cgoPeer)
	return pc, nil
}

//
// === Session Description Protocol ===
//

// CreateOffer prepares an SDP "offer" message, which should be sent to the target
// peer over a signalling channel.
func (pc *PeerConnection) CreateOffer() (*SDPHeader, error) {
	sdp := C.CGO_CreateOffer(pc.cgoPeer)
	if nil == sdp {
		return nil, errors.New("CreateOffer: could not prepare SDP offer.")
	}
	offer := new(SDPHeader)
	offer.cgoSdp = sdp
	offer.description = C.GoString(C.CGO_SerializeSDP(sdp))
	return offer, nil
}

// CreateAnswer prepares an SDP "answer" message, which should be sent in
// response to a peer that has sent an offer, over the signalling channel.
func (pc *PeerConnection) CreateAnswer() (*SDPHeader, error) {
	sdp := C.CGO_CreateAnswer(pc.cgoPeer)
	if nil == sdp {
		return nil, errors.New("CreateAnswer failed: could not prepare SDP offer.")
	}
	answer := new(SDPHeader)
	answer.cgoSdp = sdp
	answer.description = C.GoString(C.CGO_SerializeSDP(sdp))
	return answer, nil
}

func (pc *PeerConnection) SetLocalDescription(sdp *SDPHeader) error {
	r := C.CGO_SetLocalDescription(pc.cgoPeer, sdp.cgoSdp)
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
	r := C.CGO_SetRemoteDescription(pc.cgoPeer, sdp.cgoSdp)
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

// readonly signalingState
func (pc *PeerConnection) SignalingState() SignalingState {
	return (SignalingState)(C.CGO_GetSignalingState(pc.cgoPeer))
}

//
// === ICE ===
//

// TODO: change candidate into a real IceCandidate type.
func (pc *PeerConnection) AddIceCandidate(candidate string) error {
	r := C.CGO_AddIceCandidate(pc.cgoPeer, C.CString(candidate))
	if 0 != r {
		return errors.New("AddIceCandidate failed.")
	}
	return nil
}


func (pc *PeerConnection) GetConfiguration() Configuration {
	// There does not appear to be a native code version of GetConfiguration -
	// so we'll keep track of it purely from Go.
	return pc.config
	// return (Configuration)(C.CGO_GetConfiguration(pc.cgoPeer))
}

func (pc *PeerConnection) SetConfiguration(config Configuration) error {
	cConfig := config._CGO()
	if 0 != C.CGO_SetConfiguration(pc.cgoPeer, &cConfig) {
		return errors.New("PeerConnection: could not set configuration.")
	}
	pc.config = config
	return nil
}


// TODO: Above methods blocks until success or failure occurs. Maybe there should
// actually be a callback version, so the user doesn't have to make their own
// goroutine.

func (pc *PeerConnection) CreateDataChannel(label string, dict data.Init) (
	*data.Channel, error) {
	cDC := C.CGO_CreateDataChannel(pc.cgoPeer, C.CString(label), unsafe.Pointer(&dict))
	if nil == cDC {
		return nil, errors.New("Failed to CreateDataChannel")
	}
	// Convert cDC and put it in Go DC
	dc := data.NewChannel()
	return dc, nil
}

//
// === cgo hooks for user-provided Go funcs fired from C callbacks ===
//

//export cgoOnSignalingStateChange
func cgoOnSignalingStateChange(p unsafe.Pointer, s SignalingState) {
	INFO.Println("fired OnSignalingStateChange: ", p,
		s, SignalingStateString[s])
	pc := (*PeerConnection)(p)
	if nil != pc.OnSignalingStateChange {
		pc.OnSignalingStateChange(s)
	}
}

//export cgoOnNegotiationNeeded
func cgoOnNegotiationNeeded(p unsafe.Pointer) {
	INFO.Println("fired OnNegotiationNeeded: ", p)
	pc := (*PeerConnection)(p)
	if nil != pc.OnNegotiationNeeded {
		pc.OnNegotiationNeeded()
	}
}

//export cgoOnIceCandidate
func cgoOnIceCandidate(p unsafe.Pointer, candidate C.CGO_sdpString) {
	c := C.GoString(candidate)
	INFO.Println("fired OnIceCandidate: ", p, c)
	pc := (*PeerConnection)(p)
	if nil != pc.OnIceCandidate {
		pc.OnIceCandidate(c)
	}
}

//export cgoOnDataChannel
func cgoOnDataChannel(p unsafe.Pointer, cDC C.CGO_DataChannel) {
	INFO.Println("fired OnDataChannel: ", p, cDC)
	pc := (*PeerConnection)(p)
	// TODO: Convert DataChannel to Go for real.
	dc := data.NewChannel()
	if nil != pc.OnDataChannel {
		pc.OnDataChannel(dc)
	}
}
