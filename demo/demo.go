/*
 * Basic webrtc demo, two local clients speaking over go, to show the most basic
 * usage of webrtc in Go.
 *
 * For a more real-world example, see the chat demo.
 *
 * Do: `go run demo.go`
 *
 * TODO: This is in-progress.
 */
package main

import (
	"fmt"
	"github.com/keroserene/go-webrtc"
)

func main() {

	webrtc.SetLoggingVerbosity(3)

	config := webrtc.NewConfiguration(
		/// There can be as many as you like.
		webrtc.OptionIceServer("stun:stun.l.google.com:19302, stun:another"),
		webrtc.OptionIceServer("stun:another.server"),
	)

	// You can also add IceServers at a different point.
	config.AddIceServer("stun:and.another.server")

	alice, err1 := webrtc.NewPeerConnection(config)
	bob, err2 := webrtc.NewPeerConnection(config)
	if nil != err1 || nil != err2 {
		fmt.Println("Failed to create PeerConnections for both Alice and Bob.")
		return
	}
	fmt.Println("Alice and Bob created PeerConnections.\n")

	// Prepare callback handlers.
	alice.OnSignalingStateChange = func(state webrtc.SignalingState) {
		fmt.Println("Alice signal changed:", state)
	}

	// Let Alice and Bob use go channels as the signaling channel.
	// Must be bidirectional.
	a2b := make(chan *webrtc.SessionDescription, 1)
	b2a := make(chan *webrtc.SessionDescription, 1)

	wait := make(chan int, 1)

	// Start separate goroutines for Alice and Bob.
	// TODO: This will probably change, as the go webrtc interface will also
	// change.
	go func() {
		// Alice initiates the offer.
		offer, _ := alice.CreateOffer()
		alice.SetLocalDescription(offer)
		a2b <- offer
		fmt.Println("\n  Alice created and sent offer:\n", offer)

		// Now Alice waits for Bob's reply.
		answer := <-b2a
		fmt.Println("\n  Alice received answer:\n", answer)
		wait <- 1
	}()

	go func() {
		// Bob waits for alice's initial offer.
		offer := <-a2b
		bob.SetRemoteDescription(offer)
		fmt.Println("\n  Bob received offer:\n", offer)

		answer, _ := bob.CreateAnswer()
		b2a <- answer
	}()

	<-wait
	fmt.Println("done")
}
