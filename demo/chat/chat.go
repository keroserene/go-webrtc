/*
 * Webrtc chat demo.
 * Send chat messages via webrtc, over go.
 *
 * To use: `go run chat.go`
 */
package main

// #cgo LDFLAGS: -L../lib
import "C"
import (
	"github.com/keroserene/go-webrtc"
	"bufio"
	"fmt"
	"os"
)

// Signaling channel can be copy paste.
func SignalSend(msg string) {
	fmt.Println("Please provide this string to the other peer:\n")
	fmt.Println(msg)
}

func main() {

	webrtc.SetVerbosity(3)
	config := webrtc.NewConfiguration(
		webrtc.OptionIceServer("stun:stun.l.google.com:19302"))

	pc, err := webrtc.NewPeerConnection(config)
	if nil != err {
		fmt.Println("Failed to create PeerConnection.")
		return
	}

	pc.OnSignalingStateChange = func(state webrtc.SignalingState) {
		fmt.Println("signal changed:", state)
	}
	pc.OnIceCandidate = func(candidate webrtc.IceCandidate) {
		// TODO: use the serializing.
		SignalSend(candidate.Candidate)
	}
	pc.OnNegotiationNeeded = func() {
	}

	wait := make(chan int, 1)
	initiateOffer := func() {
		offer, _ := pc.CreateOffer()
		pc.SetLocalDescription(offer)
		SignalSend(offer.Description)
	}

	startAnswer := func(offer string) {
		sdp := webrtc.NewSessionDescription(offer)
		pc.SetRemoteDescription(sdp)
		answer, _ := pc.CreateAnswer()
		SignalSend(answer.Description)
	}

	fmt.Println("Initiate offer? (y/n)")
	// scan := bufio.NewScanner(os.Stdin)
	// for scan.Scan() {
		// char := scan.Text()
		// fmt.Println(scan.Text())
		// if "y" == char {
			// fmt.Println("woot")
		// }
	// }
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

	go initiateOffer()

	if false {
		go startAnswer("bad")
	}

	<-wait
	fmt.Println("done")
}
