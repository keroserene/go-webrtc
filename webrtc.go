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

TODO(keroserene): More / better documentation everywhere.
*/
package webrtc

/*
#cgo CXXFLAGS: -std=c++0x
// rpath to allow deploying same-directory corrected libstd++ in certain distros
#cgo CFLAGS: -Wl,rpath,$ORIGIN
#cgo linux,amd64 pkg-config: webrtc-linux-amd64.pc
#cgo darwin,amd64 pkg-config: webrtc-darwin-amd64.pc
#include <stdlib.h>  // Needed for C.free
#include "cpeerconnection.h"
*/
import "C"
import (
	"errors"
	"github.com/keroserene/go-webrtc/data"
	"unsafe"
)

func init() {
	SetLoggingVerbosity(3) // Default verbosity.
}

/* WebRTC PeerConnection

This is the main container of WebRTC functionality - from handling the ICE
negotiation to setting up Data Channels.

See: https://w3c.github.io/webrtc-pc/#idl-def-RTCPeerConnection
*/
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

	// Event handlers
	// TODO: Implement remainder of callbacks.
	OnIceCandidate      func(IceCandidate)
	OnIceComplete       func()
	OnNegotiationNeeded func()
	// onicecandidateerror
	OnSignalingStateChange func(SignalingState)
	// onicegatheringstatechange
	OnConnectionStateChange func(PeerConnectionState)
	OnDataChannel func(*data.Channel)

	config Configuration

	cgoPeer C.CGO_Peer // Native code internals
}

type PeerConnectionState int

/* Construct a WebRTC PeerConnection.

For a successful connection, provide at least one ICE server (stun or turn)
in the |Configuration| struct.
*/
func NewPeerConnection(config *Configuration) (*PeerConnection, error) {
	pc := new(PeerConnection)
	INFO.Println("PC at ", unsafe.Pointer(pc))
	// Internal CGO Peer wraps the native webrtc::PeerConnectionInterface.
	pc.cgoPeer = C.CGO_InitializePeer(unsafe.Pointer(pc))
	if nil == pc.cgoPeer {
		return pc, errors.New("PeerConnection: failed to initialize.")
	}
	pc.config = *config
	cConfig := config._CGO()
	if 0 != C.CGO_CreatePeerConnection(pc.cgoPeer, &cConfig) {
		return nil, errors.New("PeerConnection: could not create from config.")
	}
	INFO.Println("Created PeerConnection: ", pc, pc.cgoPeer)
	return pc, nil
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
	offer := new(SessionDescription)
	offer.cgoSdp = sdp
	offer.Type = "offer"
	offer.Sdp = C.GoString(C.CGO_SerializeSDP(sdp))
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
	answer := new(SessionDescription)
	answer.cgoSdp = sdp
	answer.Type = "answer"
	answer.Sdp = C.GoString(C.CGO_SerializeSDP(sdp))
	return answer, nil
}

/*
Set a |SessionDescription| as the local description. The description should be
generated from the local peer's CreateOffer or CreateAnswer, and not be a
description received over the signaling channel.
*/
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
	// Refresh SDP; it might have changed by ICE candidate gathering.
	pc.localDescription.Sdp = C.GoString(C.CGO_SerializeSDP(
		pc.localDescription.cgoSdp))
	return pc.localDescription
}

/*
Set a |SessionDescription| as the remote description. This description should
be one generated by the remote peer's CreateOffer or CreateAnswer, received
over the signaling channel, and not a description created locally.

If the local peer is the answerer, this must be called before CreateAnswer.
*/
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

TODO: Implement the "negotiated" flag?
*/
func (pc *PeerConnection) CreateDataChannel(label string, dict data.Init) (
	*data.Channel, error) {
	l := C.CString(label)
	defer C.free(unsafe.Pointer(l))
	cDataChannel := C.CGO_CreateDataChannel(pc.cgoPeer, l, unsafe.Pointer(&dict))
	if nil == cDataChannel {
		return nil, errors.New("Failed to CreateDataChannel")
	}
	// Provide internal Data Channel as reference to create the Go wrapper.
	dc := data.NewChannel(unsafe.Pointer(cDataChannel))
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
func cgoOnIceCandidate(p unsafe.Pointer, cIC C.CGO_IceCandidate) {
	ic := IceCandidate{
		C.GoString(cIC.sdp),
		C.GoString(cIC.sdp_mid),
		int(cIC.sdp_mline_index),
	}
	INFO.Println("fired OnIceCandidate: ", p, ic.Candidate)
	pc := (*PeerConnection)(p)
	if nil != pc.OnIceCandidate {
		pc.OnIceCandidate(ic)
	}
}

//export cgoOnIceComplete
func cgoOnIceComplete(p unsafe.Pointer) {
	INFO.Println("fired OnIceComplete: ", p)
	pc := (*PeerConnection)(p)
	if nil != pc.OnIceComplete {
		pc.OnIceComplete()
	}
}

//export cgoOnConnectionStateChange
func cgoOnConnectionStateChange(p unsafe.Pointer, state PeerConnectionState) {
	INFO.Println("fired OnConnectionStateChange: ", p)
	pc := (*PeerConnection)(p)
	if nil != pc.OnConnectionStateChange {
		pc.OnConnectionStateChange(state)
	}
}

//export cgoOnDataChannel
func cgoOnDataChannel(p unsafe.Pointer, cDC C.CGO_Channel) {
	INFO.Println("fired OnDataChannel: ", p, cDC)
	pc := (*PeerConnection)(p)
	dc := data.NewChannel(unsafe.Pointer(cDC))
	if nil != pc.OnDataChannel {
		pc.OnDataChannel(dc)
	}
}

//
// test helpers
//

func cgoFakeConnectionStateChange(p *PeerConnection, state PeerConnectionState) {
	// TODO
}
