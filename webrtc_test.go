package webrtc

import (
	"fmt"
	"testing"
)

// Prepare a PeerConnection client, and create an offer.
func TestCreateOffer(*testing.T) {
	fmt.Println("Test - [RTCPeerConnection]: CreateOffer")
	r := make(chan bool, 1)

	// Pretend to be the "Signalling" thread.
	go func() {

		pc := NewPeerConnection()
		fmt.Println("PeerConnection: ", pc)

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
	fmt.Println("Done")
}
