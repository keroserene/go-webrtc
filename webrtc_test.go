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
	header := pc.CreateOffer()
	fmt.Println("SDP created: ", header)
}

// TODO: Test video / audio stream support.
