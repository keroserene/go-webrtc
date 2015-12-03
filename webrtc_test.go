package webrtc

import (
	"fmt"
	"testing"
)

// TODO: Try Gucumber or some other fancier test framework.

var pc PeerConnection
var err error

// runtime.GOMAXPROCS = 3
// Create a PeerConnection object.
func TestCreatePeerConnection(t *testing.T) {

	StartPeerLoop()

	pc, err = NewPeerConnection()
	if nil != err {
		t.Fatal(err)
	}
}

// Use the PeerConnection client from above to create an offer,
// and see if an SDP is generated in the callback.
func TestCreateOffer(t *testing.T) {
	r := make(chan bool, 2)

	// Pretend to be the "Signalling" thread.
	go func() {
		success := func() {
			fmt.Println("success!")
		}
		failure := func() {
			t.Error("CreateOffer failed...")
		}
		pc.CreateOffer(success, failure)
		r <- true
	}()
	<-r
	fmt.Println("Done\n")
}

// func Test	
// TODO: Test video / audio stream support.
