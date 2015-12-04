package webrtc

import (
	"fmt"
	"testing"
)

// TODO: Try Gucumber or some potential fancy test framework.

var pc *PeerConnection
var err error
var sdp *SDPHeader

func TestPeerConnection(t *testing.T) {
	pc, err = NewPeerConnection()
	if nil != err {
		t.Fatal(err)
	}
}

func TestCreateOffer(t *testing.T) {
	sdp, err = pc.CreateOffer()
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println("SDP Offer:\n", sdp.description)
}

func TestSetLocalDescription(t *testing.T) {
	pc.SetLocalDescription(sdp)
}

/*
TODO: Uncomment once SetRemoteDescription is implemented.
func TestCreateAnswer(t *testing.T) {
	header, err := pc.CreateAnswer()
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println("SDP Answer: ", header.description)
}
*/

// TODO: datachannel tests
// TODO: tests for video / audio stream support.
