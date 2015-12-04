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
var config *RTCConfiguration

func TestCreatePeerConnection(t *testing.T) {
	config = NewRTCConfiguration()
	// ["stun:stun.l.google.com:19302"]
	// )
	pcA, err = NewPeerConnection(config)
	if nil != err {
		t.Fatal(err)
	}
}

func TestCreateSecondPeerConnections(t *testing.T) {
	pcB, err = NewPeerConnection(config)
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
	err = pcA.SetLocalDescription(sdp)
	if nil != err {
		t.Fatal(err)
	}
	// Pretend pcA sends the SDP offer to pcB through some signalling channel.
	fmt.Println("\n ~~ Signalling Happens here ~~ \n")
}

func TestSetRemoteDescription(t *testing.T) {
	fmt.Println("\n == BOB's PeerConnection ==")
	err = pcB.SetRemoteDescription(sdp)
	if nil != err {
		t.Fatal(err)
	}
}

// TODO: Uncomment once SetRemoteDescription is implemented.
func TestCreateAnswer(t *testing.T) {
	sdp, err := pcB.CreateAnswer()
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println("SDP Answer:\n", sdp.description)
}

// TODO: real datachannel tests
func TestCreateDataChannel(t *testing.T) {
	channel, err := pcA.CreateDataChannel("test", DataChannelInit{})
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println("Data channel: ", channel)
}

// TODO: tests for video / audio stream support.
