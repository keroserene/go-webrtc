package webrtc

import (
	"fmt"
	"testing"
)

var pc PeerConnection
var err error

// Prepare a PeerConnection client, and create an offer.
func TestCreatePeerConnection(t *testing.T) {
	fmt.Println("\nTest: [PeerConnection] Creating")
	pc, err = NewPeerConnection()
	if nil != err {
		t.Fatal(err)
	}
	// fmt.Printf("PeerConnection: %+v\n", pc)
}

func TestCreateOffer(t *testing.T) {
	fmt.Println("\nTest: [PeerConnection] CreateOffer")
	r := make(chan bool, 1)

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

// TODO: Test video / audio stream support.
