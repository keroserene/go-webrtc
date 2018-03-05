/*
Package webrtc is a golang wrapper on native code WebRTC.

For consistency with the browser-based WebRTCs, the interface here is based
loosely on: w3c.github.io/webrtc-pc

The main goal of this project is to present a golang WebRTC package in the most
idiomatic and simple-to-use way.

However, to provide a better experience for users of this package, there are
differences inherent in the interface written here and the original native code
WebRTC - from the golang requirement of Capitalized identifiers for public
interfaces, to the replacement of certain callbacks with goroutines.

Note that building the necessary libwebrtc static library is excessively
complicated, which is why the necessary platform-specific archives will be
provided in lib/. This also mitigates the possibility that future commits on
native libwebrtc will break go-webrtc, because the interface with the native
code, through the intermediate CGO layer, is relatively fragile.

Due to other external goals of the developers, this package will only be
focused on DataChannels. However, extending this package to allow video/audio
media streams and related functionality, to be a "complete" WebRTC suite,
is entirely possible and will likely happen in the long term. (Issue #7)
This will however have implications for the archives that need to be built
and linked.

Please share any improvements or concerns as issues or pull requests on github.
*/
package webrtc

/*
#cgo CXXFLAGS: -std=c++0x
#cgo LDFLAGS: -L${SRCDIR}/lib
#cgo linux,arm pkg-config: webrtc-linux-arm.pc
#cgo linux,386 pkg-config: webrtc-linux-386.pc
#cgo linux,amd64 pkg-config: webrtc-linux-amd64.pc
#cgo darwin,amd64 pkg-config: webrtc-darwin-amd64.pc
#include <stdlib.h>  // Needed for C.free
#include "peerconnection.h"
#include "ctestenums.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

func init() {
	SetLoggingVerbosity(3) // Default verbosity.
}

type (
	PeerConnectionState int
	IceGatheringState   int
	IceConnectionState  int
)

const (
	PeerConnectionStateNew PeerConnectionState = iota
	PeerConnectionStateConnecting
	PeerConnectionStateConnected
	PeerConnectionStateDisconnected
	PeerConnectionStateFailed
)

func (s PeerConnectionState) String() string {
	return EnumToStringSafe(int(s), []string{
		"New",
		"Connecting",
		"Connected",
		"Disconnected",
		"Failed",
	})
}

const (
	IceConnectionStateNew IceConnectionState = iota
	IceConnectionStateChecking
	IceConnectionStateConnected
	IceConnectionStateCompleted
	IceConnectionStateFailed
	IceConnectionStateDisconnected
	IceConnectionStateClosed
)

func (s IceConnectionState) String() string {
	return EnumToStringSafe(int(s), []string{
		"New",
		"Checking",
		"Connected",
		"Completed",
		"Failed",
		"Disconnected",
		"Closed",
	})
}

const (
	IceGatheringStateNew IceGatheringState = iota
	IceGatheringStateGathering
	IceGatheringStateComplete
)

func (s IceGatheringState) String() string {
	return EnumToStringSafe(int(s), []string{
		"New",
		"Gathering",
		"Complete",
	})
}

var PCMap = NewCGOMap()

/* WebRTC PeerConnection

This is the main container of WebRTC functionality - from handling the ICE
negotiation to setting up Data Channels.

See: https://w3c.github.io/webrtc-pc/#idl-def-RTCPeerConnection
*/
type PeerConnection struct {
	localDescription        *SessionDescription
	remoteDescription       *SessionDescription
	canTrickleIceCandidates bool

	// Event handlers
	OnNegotiationNeeded        func()
	OnIceCandidate             func(IceCandidate)
	OnIceCandidateError        func()
	OnIceComplete              func() // Possibly to be removed.
	OnSignalingStateChange     func(SignalingState)
	OnIceConnectionStateChange func(IceConnectionState)
	OnIceGatheringStateChange  func(IceGatheringState)
	OnConnectionStateChange    func(PeerConnectionState)
	OnDataChannel              func(*DataChannel)

	config Configuration

	cgoPeer C.CGO_Peer // Native code internals
	index   int        // Index into the PCMap
}

/* Construct a WebRTC PeerConnection.

For a successful connection, provide at least one ICE server (stun or turn)
in the |Configuration| struct.
*/
func NewPeerConnection(config *Configuration) (*PeerConnection, error) {
	if nil == config {
		return nil, errors.New("PeerConnection requires a Configuration.")
	}
	pc := new(PeerConnection)
	pc.index = PCMap.Set(pc)
	// Internal CGO Peer wraps the native webrtc::PeerConnectionInterface.
	pc.cgoPeer = C.CGO_InitializePeer(C.int(pc.index))
	if nil == pc.cgoPeer {
		return pc, errors.New("PeerConnection: failed to initialize.")
	}
	pc.config = *config
	cConfig := config._CGO()
	defer freeConfig(cConfig)
	if 0 != C.CGO_CreatePeerConnection(pc.cgoPeer, cConfig) {
		return nil, errors.New("PeerConnection: could not create from config.")
	}
	INFO.Println("Created PeerConnection: ", pc, pc.cgoPeer)
	return pc, nil
}

func (pc *PeerConnection) Destroy() error {
	err := pc.Close()
	PCMap.Delete(pc.index)
	C.CGO_DestroyPeer(pc.cgoPeer)
	return err
}

//
// === Session Description Protocol ===
//

/*
CreateOffer prepares an SDP "offer" message, which should be set as the local
description, then sent to the remote peer over a signalling channel. This
should only be called by the peer initiating the connection.

This method is blocking, and should occur within a separate goroutine.
*/
func (pc *PeerConnection) CreateOffer() (*SessionDescription, error) {
	sdp := C.CGO_CreateOffer(pc.cgoPeer)
	if nil == sdp {
		return nil, errors.New("CreateOffer: could not prepare SDP offer.")
	}
	offer := NewSessionDescription("offer", sdp)
	if offer == nil {
		return nil, errors.New("CreateOffer: could not prepare SDP offer.")
	}
	return offer, nil
}

/*
CreateAnswer prepares an SDP "answer" message. This should only happen in
response to an offer received and set as the remote description. Once generated,
this answer should then be set as the local description and sent back over the
signaling channel to the remote peer.

This method is blocking, and should occur within a separate goroutine.
*/
func (pc *PeerConnection) CreateAnswer() (*SessionDescription, error) {
	sdp := C.CGO_CreateAnswer(pc.cgoPeer)
	if nil == sdp {
		return nil, errors.New("CreateAnswer failed: could not prepare SDP offer.")
	}
	answer := NewSessionDescription("answer", sdp)
	if answer == nil {
		return nil, errors.New("CreateAnswer failed: could not prepare SDP offer.")
	}
	return answer, nil
}

/*
Set a |SessionDescription| as the local description. The description should be
generated from the local peer's CreateOffer or CreateAnswer, and not be a
description received over the signaling channel.
*/
func (pc *PeerConnection) SetLocalDescription(sdp *SessionDescription) error {
	if nil == sdp {
		return errors.New("Cannot use nil SessionDescription.")
	}
	r := C.CGO_SetLocalDescription(pc.cgoPeer, sdp.GoStringToCgoSdp())
	if 0 != r {
		return errors.New("SetLocalDescription failed.")
	}
	pc.localDescription = sdp
	return nil
}

// readonly localDescription
func (pc *PeerConnection) LocalDescription() (sdp *SessionDescription) {
	// Refresh SDP; it might have changed by ICE candidate gathering.
	if pc.localDescription != nil {
		cgoSdp := C.CGO_GetLocalDescription(pc.cgoPeer)
		pc.localDescription.Sdp = CgoSdpToGoString(cgoSdp)
	}
	return pc.localDescription
}

/*
Set a |SessionDescription| as the remote description. This description should
be one generated by the remote peer's CreateOffer or CreateAnswer, received
over the signaling channel, and not a description created locally.

If the local peer is the answerer, this must be called before CreateAnswer.
*/
func (pc *PeerConnection) SetRemoteDescription(sdp *SessionDescription) error {
	if nil == sdp {
		return errors.New("Cannot use nil SessionDescription.")
	}
	r := C.CGO_SetRemoteDescription(pc.cgoPeer, sdp.GoStringToCgoSdp())
	if 0 != r {
		return errors.New("SetRemoteDescription failed.")
	}
	pc.remoteDescription = sdp
	return nil
}

// readonly remoteDescription
func (pc *PeerConnection) RemoteDescription() (sdp *SessionDescription) {
	if pc.remoteDescription != nil {
		cgoSdp := C.CGO_GetRemoteDescription(pc.cgoPeer)
		pc.remoteDescription.Sdp = CgoSdpToGoString(cgoSdp)
	}
	return pc.remoteDescription
}

// readonly signalingState
func (pc *PeerConnection) SignalingState() SignalingState {
	return (SignalingState)(C.CGO_GetSignalingState(pc.cgoPeer))
}

// readonly connectionState
func (pc *PeerConnection) ConnectionState() PeerConnectionState {
	// TODO: Aggregate states according to:
	// https://w3c.github.io/webrtc-pc/#rtcpeerconnectionstate-enum
	return (PeerConnectionState)(C.CGO_IceConnectionState(pc.cgoPeer))
}

// readonly icegatheringstatee
func (pc *PeerConnection) IceGatheringState() IceGatheringState {
	return (IceGatheringState)(C.CGO_IceGatheringState(pc.cgoPeer))
}

// readonly iceconnectionState
func (pc *PeerConnection) IceConnectionState() IceConnectionState {
	return (IceConnectionState)(C.CGO_IceConnectionState(pc.cgoPeer))
}

func (pc *PeerConnection) AddIceCandidate(ic IceCandidate) error {
	sdpMid := C.CString(ic.SdpMid)
	defer C.free(unsafe.Pointer(sdpMid))
	sdp := C.CString(ic.Candidate)
	defer C.free(unsafe.Pointer(sdp))

	cIC := new(C.CGO_IceCandidate)
	cIC.sdp_mid = sdpMid
	cIC.sdp_mline_index = C.int(ic.SdpMLineIndex)
	cIC.sdp = sdp

	r := C.CGO_AddIceCandidate(pc.cgoPeer, cIC)
	if 0 != r {
		return errors.New("AddIceCandidate failed.")
	}
	return nil
}

func (pc *PeerConnection) GetConfiguration() Configuration {
	// There does not appear to be a native code version of GetConfiguration -
	// so we'll keep track of it purely from Go.
	return pc.config
}

func (pc *PeerConnection) SetConfiguration(config Configuration) error {
	cConfig := config._CGO()
	defer freeConfig(cConfig)
	err := C.CGO_SetConfiguration(pc.cgoPeer, cConfig)
	if err != 0 {
		return errors.New(fmt.Sprintf("PeerConnection: could not set configuration. Error ID: %d", err))
	}
	pc.config = config
	return nil
}

/*
Create and return a DataChannel.

This only needs to be called by one side, unless "negotiated" is true.

If creating the first DataChannel, this actually triggers the local
PeerConnection's .OnNegotiationNeeded callback, which should lead to a
user-provided goroutine containing CreateOffer, SetLocalDescription, and the
rest of the signalling exchange.

Once the connection succeeds, .OnDataChannel should trigger on the remote peer's
|PeerConnection|, while .OnOpen should trigger on the local DataChannel returned
by this method. Both DataChannel references should then be open and ready to
exchange data.
*/

// Ordered configures a DataChannels 'ordered' option.
func Ordered(ordered bool) func(*DataChannelInit) {
	return func(i *DataChannelInit) {
		i.Ordered = ordered
	}
}

// MaxPacketLifeTime configures a DataChannels 'maxRetransmitTime' option.
func MaxPacketLifeTime(maxPacketLifeTime int) func(*DataChannelInit) {
	return func(i *DataChannelInit) {
		i.MaxPacketLifeTime = maxPacketLifeTime
	}
}

// MaxRetransmits configures a DataChannels 'maxRetransmits' option.
func MaxRetransmits(maxRetransmits int) func(*DataChannelInit) {
	return func(i *DataChannelInit) {
		i.MaxRetransmits = maxRetransmits
	}
}

// Negotiated configures a DataChannels 'negotiated' option.
func Negotiated(negotiated bool) func(*DataChannelInit) {
	return func(i *DataChannelInit) {
		i.Negotiated = negotiated
	}
}

func (pc *PeerConnection) CreateDataChannel(label string, options ...func(*DataChannelInit)) (
	*DataChannel, error) {

	// These are the defaults taken from include/webrtc/api/datachannelinterface.h
	init := DataChannelInit{
		Ordered:           true,
		MaxPacketLifeTime: -1,
		MaxRetransmits:    -1,
		Negotiated:        false,
		ID:                -1,
	}

	for _, option := range options {
		option(&init)
	}

	cfg := C.CGO_DataChannelInit{}
	cfg.ordered = 1
	if init.Ordered == false {
		cfg.ordered = 0
	}
	cfg.negotiated = 0
	if init.Negotiated == true {
		cfg.negotiated = 1
	}
	cfg.id = C.int(init.ID)
	cfg.maxRetransmits = C.int(init.MaxRetransmits)
	cfg.maxPacketLifeTime = C.int(init.MaxPacketLifeTime)

	l := C.CString(label)
	defer C.free(unsafe.Pointer(l))
	cDataChannel := C.CGO_CreateDataChannel(pc.cgoPeer, l, cfg)
	if nil == cDataChannel {
		return nil, errors.New("Failed to CreateDataChannel")
	}
	// Provide internal Data Channel as reference to create the Go wrapper.
	dc := NewDataChannel(unsafe.Pointer(cDataChannel))
	return dc, nil
}

func (pc *PeerConnection) DeleteDataChannel(dc *DataChannel) {
	dc.Close()
	C.CGO_DeleteDataChannel(pc.cgoPeer, dc.cgoChannelObserver)
	deleteDataChannel(dc.index)
	return
}

func (pc *PeerConnection) Close() error {
	C.CGO_Close(pc.cgoPeer)
	return nil
}

//
// === cgo hooks for user-provided Go funcs fired from C callbacks ===
//

//export cgoOnSignalingStateChange
func cgoOnSignalingStateChange(p int, s SignalingState) {
	INFO.Println("fired OnSignalingStateChange: ", p, s)
	pc := PCMap.Get(p).(*PeerConnection)
	if nil != pc.OnSignalingStateChange {
		pc.OnSignalingStateChange(s)
	}
}

//export cgoOnNegotiationNeeded
func cgoOnNegotiationNeeded(p int) {
	INFO.Println("fired OnNegotiationNeeded: ", p)
	pc := PCMap.Get(p).(*PeerConnection)
	if nil != pc.OnNegotiationNeeded {
		pc.OnNegotiationNeeded()
	}
}

//export cgoOnIceCandidate
func cgoOnIceCandidate(p int, cIC C.CGO_IceCandidate) {
	ic := IceCandidate{
		C.GoString(cIC.sdp),
		C.GoString(cIC.sdp_mid),
		int(cIC.sdp_mline_index),
	}
	INFO.Println("fired OnIceCandidate: ", p, ic.Candidate)
	pc := PCMap.Get(p).(*PeerConnection)
	if nil != pc.OnIceCandidate {
		pc.OnIceCandidate(ic)
	}
}

//export cgoOnIceCandidateError
func cgoOnIceCandidateError(p int) {
	INFO.Println("fired OnIceCandidateError: ", p)
	pc := PCMap.Get(p).(*PeerConnection)
	if nil != pc.OnIceCandidateError {
		pc.OnIceCandidateError()
	}
}

//export cgoOnConnectionStateChange
func cgoOnConnectionStateChange(p int, iceState IceConnectionState) {
	// TODO: This may need to be slightly more complicated...
	// https://w3c.github.io/webrtc-pc/#rtcpeerconnectionstate-enum
	var state PeerConnectionState
	switch iceState {
	case IceConnectionStateNew:
		state = PeerConnectionStateNew
	case IceConnectionStateChecking:
		state = PeerConnectionStateConnecting
	case IceConnectionStateConnected:
		state = PeerConnectionStateConnected
	case IceConnectionStateFailed:
		state = PeerConnectionStateFailed
	case IceConnectionStateDisconnected:
		state = PeerConnectionStateDisconnected
	default:
		return
	}

	INFO.Println("fired OnConnectionStateChange: ", p)
	pc := PCMap.Get(p).(*PeerConnection)
	if nil != pc.OnConnectionStateChange {
		pc.OnConnectionStateChange(state)
	}
}

//export cgoOnIceConnectionStateChange
func cgoOnIceConnectionStateChange(p int, state IceConnectionState) {
	INFO.Println("fired OnIceConnectionStateChange: ", p)
	pc := PCMap.Get(p).(*PeerConnection)
	if nil != pc.OnIceConnectionStateChange {
		pc.OnIceConnectionStateChange(state)
	}
}

//export cgoOnIceGatheringStateChange
func cgoOnIceGatheringStateChange(p int, state IceGatheringState) {
	INFO.Println("fired OnIceGatheringStateChange:", p)
	pc := PCMap.Get(p).(*PeerConnection)
	if nil != pc.OnIceGatheringStateChange {
		pc.OnIceGatheringStateChange(state)
	}
	// Although OnIceComplete is to be deprecated in the native API, and no longer
	// part of the w3 spec, keeping it for go seems easier for the users.
	if IceGatheringStateComplete == state && nil != pc.OnIceComplete {
		pc.OnIceComplete()
	}
}

//export cgoOnDataChannel
func cgoOnDataChannel(p int, o unsafe.Pointer) {
	INFO.Println("fired OnDataChannel: ", p, o)
	pc := PCMap.Get(p).(*PeerConnection)
	dc := NewDataChannel(o)
	if nil != pc.OnDataChannel {
		pc.OnDataChannel(dc)
	}
}

// Test helpers
//
var _cgoIceConnectionStateNew = int(C.CGO_IceConnectionStateNew)
var _cgoIceConnectionStateChecking = int(C.CGO_IceConnectionStateChecking)
var _cgoIceConnectionStateConnected = int(C.CGO_IceConnectionStateConnected)
var _cgoIceConnectionStateCompleted = int(C.CGO_IceConnectionStateCompleted)
var _cgoIceConnectionStateFailed = int(C.CGO_IceConnectionStateFailed)
var _cgoIceConnectionStateDisconnected = int(C.CGO_IceConnectionStateDisconnected)
var _cgoIceConnectionStateClosed = int(C.CGO_IceConnectionStateClosed)

var _cgoIceGatheringStateNew = int(C.CGO_IceGatheringStateNew)
var _cgoIceGatheringStateGathering = int(C.CGO_IceGatheringStateGathering)
var _cgoIceGatheringStateComplete = int(C.CGO_IceGatheringStateComplete)

func cgoFakeIceCandidateError(pc *PeerConnection) {
	C.CGO_fakeIceCandidateError(pc.cgoPeer)
}
