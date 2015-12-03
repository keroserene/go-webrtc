package webrtc

import (
	"fmt"
	"testing"
	"runtime"
)

// TODO: Try Gucumber or some other fancier test framework.

var pc PeerConnection
var err error

// runtime.GOMAXPROCS = 3
// Create a PeerConnection object.
func TestCreatePeerConnection(t *testing.T) {
	fmt.Println(runtime.NumCPU())
	x := runtime.GOMAXPROCS(4)
	fmt.Println(x)
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
			fmt.Println("CreateOffer Succeeded.")
		}
		failure := func() {
			fmt.Println("CreateOffer Failed.")
		}
		pc.CreateOffer(success, failure)
		r <- true
	}()
	<-r
	fmt.Println("Done\n")
}

// func Test	
// TODO: Test video / audio stream support.
