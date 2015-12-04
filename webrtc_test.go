package webrtc

import (
	"fmt"
	"testing"
)

// TODO: Try Gucumber or some potential fancy test framework.

// These tests create two PeerConnections objects, which allows a loopback test.
var pcA *PeerConnection
var pcB *PeerConnection
var err error
var sdp *SDPHeader

func TestCreatePeerConnection(t *testing.T) {
	pcA, err = NewPeerConnection()
	if nil != err {
		t.Fatal(err)
	}
}

func TestCreateSecondPeerConnections(t *testing.T) {
	pcB, err = NewPeerConnection()
	if nil != err {
		t.Fatal(err)
	}
}

func TestCreateOffer(t *testing.T) {
	fmt.Println("\n== ALICE's PeerConnection ==")
	sdp, err = pcA.CreateOffer()
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println("SDP Offer:\n", sdp.description)
}

func TestSetLocalDescription(t *testing.T) {
	pcA.SetLocalDescription(sdp)
	// Pretend pcA sends the SDP offer to pcB through some signalling channel.
	fmt.Println("\n ~~ Signalling Happens here ~~ \n")
}

func TestSetRemoteDescription(t *testing.T) {
	fmt.Println("\n == BOB's PeerConnection ==")
	pcB.SetRemoteDescription(sdp)
}

// TODO: Uncomment once SetRemoteDescription is implemented.
func TestCreateAnswer(t *testing.T) {
	sdp, err := pcB.CreateAnswer()
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println("SDP Answer:\n", sdp.description)
}

// TODO: datachannel tests
// TODO: tests for video / audio stream support.
