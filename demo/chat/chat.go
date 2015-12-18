/*
 * Webrtc chat demo.
 * Send chat messages via webrtc, over go.
 * Can interop with the JS client. (Open chat.html in a browser)
 *
 * To use: `go run chat.go`
 */
package main

import "C"
import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/keroserene/go-webrtc"
	"github.com/keroserene/go-webrtc/data"
	"os"
	"strings"
)

var pc *webrtc.PeerConnection
var dc *data.Channel
var mode Mode
var err error
var username = "Alice"

// Janky state machine.
type Mode int

const (
	ModeInit Mode = iota
	ModeConnect
	ModeChat
)

// Signaling channel can be copy paste.
func signalSend(msg string) {
	fmt.Println("\n ---- Please copy the below to peer ---- \n")
	fmt.Println(msg + "\n")
}

func signalReceive(msg string) {
	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(msg), &parsed)
	if nil != err {
		fmt.Println(err, ", try again.")
		return
	}

	// If this is a valid signal and no PeerConnection has been instantiated,
	// start as the "answerer."
	if nil == pc {
		start(false)
	}

	// TODO: Allow multiple candidates combined with an offer/answer description
	// as a single json string to copy paste.
	if nil != parsed["sdp"] {
		sdp := webrtc.DeserializeSessionDescription(msg)
		if nil == sdp {
			fmt.Println("Invalid SDP.")
			return
		}
		receiveDescription(sdp)
	}
	if nil != parsed["candidate"] {
		ice := webrtc.DeserializeIceCandidate(msg)
		if nil == ice {
			fmt.Println("Invalid ICE candidate.")
			return
		}
		pc.AddIceCandidate(*ice)
		fmt.Println("ICE candidate successfully received.")
	}
}

func sendOffer() {
	offer, err := pc.CreateOffer()
	if err != nil {
		fmt.Println(err)
		return
	}
	pc.SetLocalDescription(offer)
	signalSend(offer.Serialize())
}

func sendAnswer() {
	answer, err := pc.CreateAnswer()
	if err != nil {
		fmt.Println(err)
		return
	}
	pc.SetLocalDescription(answer)
	signalSend(answer.Serialize())
}

func receiveDescription(sdp *webrtc.SessionDescription) {
	err = pc.SetRemoteDescription(sdp)
	if nil != err {
		fmt.Println("ERROR", err)
		return
	}
	fmt.Println("SDP " + sdp.Type + " successfully received.")
	if "offer" == sdp.Type {
		fmt.Println("Generating and replying with an answer.")
		go sendAnswer()
	}
}

func prepareDataChannel(channel *data.Channel) {
	channel.OnOpen = func() {
		fmt.Println("Data Channel opened!")
		startChat()
	}
	channel.OnClose = func() {
		fmt.Println("Data Channel closed.")
		endChat()
	}
	channel.OnMessage = func(msg []byte) {
		receiveChat(string(msg))
	}
}

func sendChat(msg string) {
	line := username + ": " + msg
	fmt.Println("[sent]")
	dc.Send(line)
}

func receiveChat(msg string) {
	fmt.Println("\n" + string(msg))
}

func startChat() {
	mode = ModeChat
	fmt.Println("------- chat enabled! -------")
}

func endChat() {
	mode = ModeInit
	fmt.Println("------- chat disabled -------")
}

func start(instigator bool) {
	mode = ModeConnect
	fmt.Println("Starting up PeerConnection...")
	config := webrtc.NewConfiguration(
		webrtc.OptionIceServer("stun:stun.l.google.com:19302"))
	pc, err = webrtc.NewPeerConnection(config)
	if nil != err {
		fmt.Println("Failed to create PeerConnection.")
		return
	}

	// The below three callbacks are the minimum required.
	pc.OnNegotiationNeeded = func() {
		go sendOffer()
	}
	pc.OnIceCandidate = func(candidate webrtc.IceCandidate) {
		signalSend(candidate.Serialize())
	}
	pc.OnDataChannel = func(channel *data.Channel) {
		fmt.Println("Datachannel established...", channel)
		dc = channel
		prepareDataChannel(channel)
	}

	if instigator {
		// Attempting to create the first datachannel triggers ICE.
		fmt.Println("Trying to create a datachannel.")
		dc, err = pc.CreateDataChannel("test", data.Init{})
		if nil != err {
			fmt.Println("Unexpected failure creating data.Channel.")
			return
		}
		prepareDataChannel(dc)
	}
}

func main() {
	webrtc.SetLoggingVerbosity(1)
	mode = ModeInit
	reader := bufio.NewReader(os.Stdin)

	wait := make(chan int, 1)
	fmt.Println("=== go-webrtc chat demo ===")
	fmt.Println("What is your username?")
	username, _ = reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Println("Welcome, " + username + "!")
	fmt.Println("To initiate a WebRTC PeerConnection, type \"start\".")
	fmt.Println("(Alternatively, immediately input SDP messages from the peer.)")

	// Input loop.
	for true {
		text, _ := reader.ReadString('\n')
		switch mode {
		case ModeInit:
			if strings.HasPrefix(text, "start") {
				start(true)
			} else {
				signalReceive(text)
			}
		case ModeConnect:
			signalReceive(text)
		case ModeChat:
			sendChat(text)
			// fmt.Print(username + ": ")
			break
		}
		text = ""
	}
	<-wait
	fmt.Println("done")
}
