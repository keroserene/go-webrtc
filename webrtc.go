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
#include <stdlib.h>
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

type SessionDescription struct {
	// Keep track of both a pointer to the C++ SessionDescription object,
	// and the serialized string version (which native code generates)
	cgoSdp      C.CGO_sdp
	Description string
}

// Construct a SessionDescription object from a valid msg.
func NewSessionDescription(sdpType string, msg string) *SessionDescription {
	cSdp := C.CGO_DeserializeSDP(C.CString(sdpType), C.CString(msg))
	if nil == cSdp {
		ERROR.Println("Invalid SDP string.")
		return nil
	}
	sdp := new(SessionDescription)
	sdp.cgoSdp = cSdp
	sdp.Description = msg
	return sdp
}

type PeerConnection struct {
	localDescription *SessionDescription
	// currentLocalDescription
	// pendingLocalDescription

	remoteDescription *SessionDescription
	// currentRemoteDescription
	// pendingRemoteDescription

	// iceGatheringState  RTCIceGatheringState
	// iceConnectionState  RTCIceConnectionState
	canTrickleIceCandidates bool

	// Event handlers TODO: The remainder of the callbacks.
	OnIceCandidate      func(IceCandidate)
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
func (pc *PeerConnection) CreateOffer() (*SessionDescription, error) {
	sdp := C.CGO_CreateOffer(pc.cgoPeer)
	if nil == sdp {
		return nil, errors.New("CreateOffer: could not prepare SDP offer.")
	}
	offer := new(SessionDescription)
	offer.cgoSdp = sdp
	offer.Description = C.GoString(C.CGO_SerializeSDP(sdp))
	return offer, nil
}

// CreateAnswer prepares an SDP "answer" message, which should be sent in
// response to a peer that has sent an offer, over the signalling channel.
func (pc *PeerConnection) CreateAnswer() (*SessionDescription, error) {
	sdp := C.CGO_CreateAnswer(pc.cgoPeer)
	if nil == sdp {
		return nil, errors.New("CreateAnswer failed: could not prepare SDP offer.")
	}
	answer := new(SessionDescription)
	answer.cgoSdp = sdp
	answer.Description = C.GoString(C.CGO_SerializeSDP(sdp))
	return answer, nil
}

// TODO: Above methods blocks until success or failure occurs. Maybe there should
// actually be a callback version, so the user doesn't have to make their own
// goroutine.

func (pc *PeerConnection) SetLocalDescription(sdp *SessionDescription) error {
	r := C.CGO_SetLocalDescription(pc.cgoPeer, sdp.cgoSdp)
	if 0 != r {
		return errors.New("SetLocalDescription failed.")
	}
	pc.localDescription = sdp
	return nil
}

// readonly localDescription
func (pc *PeerConnection) LocalDescription() (sdp *SessionDescription) {
	return pc.localDescription
}

func (pc *PeerConnection) SetRemoteDescription(sdp *SessionDescription) error {
	r := C.CGO_SetRemoteDescription(pc.cgoPeer, sdp.cgoSdp)
	if 0 != r {
		return errors.New("SetRemoteDescription failed.")
	}
	pc.remoteDescription = sdp
	return nil
}

// readonly remoteDescription
func (pc *PeerConnection) RemoteDescription() (sdp *SessionDescription) {
	return pc.remoteDescription
}

// readonly signalingState
func (pc *PeerConnection) SignalingState() SignalingState {
	return (SignalingState)(C.CGO_GetSignalingState(pc.cgoPeer))
}

//
// === ICE / Configuration ===
//

type (
	IceProtocol         int
	IceCandidateType    int
	IceTcpCandidateType int
)

const (
	IceProtocolUPD IceProtocol = iota
	IceProtocolTCP
)

var IceProtocolString = []string{"udp", "tcp"}

const (
	IceCandidateTypeHost IceCandidateType = iota
	IceCandidateTypeSrflx
	IceCandidateTypePrflx
	IceCandidateTypeRelay
)

var IceCandidateTypeString = []string{"host", "srflx", "prflx", "relay"}

const (
	IceTcpCandidateTypeActive IceTcpCandidateType = iota
	IceTcpCandidateTypePassive
	IceTcpCandidateTypeSo
)

var IceTcpCandidateTypeString = []string{"active", "passive", "so"}

type IceCandidate struct {
	Candidate     string `json:"candidate"`
	SdpMid        string `json:"sdpMid"`
	SdpMLineIndex int    `json:"sdpMLineIndex"`
	// Foundation     string
	// Priority       C.ulong
	// IP             net.IP
	// Protocol       IceProtocol
	// Port           C.ushort
	// Type           IceCandidateType
	// TcpType        IceTcpCandidateType
	// RelatedAddress string
	// RelatedPort    C.ushort
}

func (pc *PeerConnection) AddIceCandidate(ic IceCandidate) error {
	candidate := C.CString(ic.Candidate)
	defer C.free(unsafe.Pointer(candidate))
	sdpMid := C.CString(ic.SdpMid)
	defer C.free(unsafe.Pointer(sdpMid))
	r := C.CGO_AddIceCandidate(pc.cgoPeer, candidate, sdpMid, C.int(ic.SdpMLineIndex))
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

func (pc *PeerConnection) CreateDataChannel(label string, dict data.Init) (
	*data.Channel, error) {
	cDC := C.CGO_CreateDataChannel(pc.cgoPeer, C.CString(label), unsafe.Pointer(&dict))
	if nil == cDC {
		return nil, errors.New("Failed to CreateDataChannel")
	}
	// Convert cDC and put it in Go DC
	dc := data.NewChannel(unsafe.Pointer(cDC))
	return dc, nil
}

func (pc *PeerConnection) Close() error {
	C.CGO_Close(pc.cgoPeer)
	return nil
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
func cgoOnIceCandidate(p unsafe.Pointer, candidate *C.char, sdpMid *C.char, sdpMLineIndex int) {
	ic := IceCandidate{
		C.GoString(candidate),
		C.GoString(sdpMid),
		sdpMLineIndex,
	}
	INFO.Println("fired OnIceCandidate: ", p, ic.Candidate)
	pc := (*PeerConnection)(p)
	if nil != pc.OnIceCandidate {
		pc.OnIceCandidate(ic)
	}
}

type CDC C.CGO_Channel

// func cgoOnDataChannel(p unsafe.Pointer, cDC C.CGO_Channel) {
//export cgoOnDataChannel
func cgoOnDataChannel(p unsafe.Pointer, cDC CDC) {
	INFO.Println("fired OnDataChannel: ", p, cDC)
	pc := (*PeerConnection)(p)
	// TODO: Convert DataChannel to Go for real.
	dc := data.NewChannel(unsafe.Pointer(cDC))
	// C.CGO_Channel)(cDC))
	if nil != pc.OnDataChannel {
		pc.OnDataChannel(dc)
	}
}
