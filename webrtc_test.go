package webrtc

import (
	"fmt"
	"testing"
)

// TODO: Try Gucumber or some potential fancy test framework.

var pc *PeerConnection
var err error

func TestPeerConnection(t *testing.T) {
	InitializePeer()
	pc, err = NewPeerConnection()
	if nil != err {
		t.Fatal(err)
	}
}

func TestCreateOffer(t *testing.T) {
	header, err := pc.CreateOffer()
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println("SDP Offer: ", header.description)
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

// TODO: Test video / audio stream support.
