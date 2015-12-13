/*
 * Webrtc chat demo.
 * Send chat messages via webrtc, over go.
 * Can interop with the JS client. (Open chat.html in a browser)
 *
 * To use: `go run chat.go`
 */
package main

// #cgo LDFLAGS: -L../lib
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

type Mode int

const (
	ModeInit Mode = iota
	ModeConnect
	ModeChat
)

var pc *webrtc.PeerConnection
var dc *data.Channel
var mode Mode
var err error
var username = "Alice"

// Signaling channel can be copy paste.
func signalSend(msg string) {
	fmt.Println("\n ---- Please copy the below to peer ---- \n")
	fmt.Println(msg + "\n")
}

func signalRecieve(msg string) {
	fmt.Println("Received: ", msg)
	// recv := msg
	var v interface{}
	json.Unmarshal([]byte(msg), v)

	// Only do this once it's a valid SDP.
	if nil == pc {
		start(false)
	}
}


type (
  Description struct {
		Type string `json:"type"`
		Sdp string  `json:"sdp"`
	}
)

func sendOffer() {
	offer, err := pc.CreateOffer()
	if err != nil {
		fmt.Println(err)
		return
	}
	pc.SetLocalDescription(offer)
	data := Description{
			"offer",
			offer.Description,
		}
	bytes, err := json.Marshal(data)
	message := fmt.Sprintf(`{"desc":%s}`, string(bytes))
	signalSend(message)
}

func receiveDescription(desc string) {
	pc.SetRemoteDescription(desc)
}

// func	startAnswer := func(offer string) {
// sdp := webrtc.NewSessionDescription(offer)
// pc.SetRemoteDescription(sdp)
// answer, _ := pc.CreateAnswer()
// signalSend(answer.Description)
// }
func prepareDataChannel(channel *data.Channel) {
	channel.OnOpen = func() {
		fmt.Println("Data Channel opened!")
	}
	channel.OnClose = func() {
		fmt.Println("Data Channel closed.")
	}
	channel.OnMessage = func(msg []byte) {
		fmt.Println(string(msg))
	}
}

func sendChat(msg string) {
	line := username + ": " + msg
	fmt.Println(line)
	dc.Send([]byte(line))
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

	pc.OnNegotiationNeeded = func() {
		// Initial SDP offer / answer exchange.
		go sendOffer()
	}
	pc.OnSignalingStateChange = func(state webrtc.SignalingState) {
		fmt.Println("signal changed:", state)
	}
	pc.OnIceCandidate = func(candidate webrtc.IceCandidate) {
		// TODO: use the serializing.
		signalSend(candidate.Candidate)
	}
	pc.OnDataChannel = func(channel *data.Channel) {
		fmt.Println("Datachannel established...", channel)
		prepareDataChannel(channel)
	}

	if instigator {
		// Attempting to create the first datachannel triggers ICE.
		fmt.Println("Trying to create a datachannel.")
		dc, err = pc.CreateDataChannel("init", data.Init{})
		if nil != err {
			fmt.Println("Unexpected failure creating data.Channel.")
			return;
		}
		prepareDataChannel(dc)
	}
}

func main() {
	webrtc.SetVerbosity(3)

	mode = ModeInit
	wait := make(chan int, 1)
	fmt.Println("To initiate PeerConnection, type \"start\".")
	fmt.Println("Alternatively, input SDP messages from peer.")
	reader := bufio.NewReader(os.Stdin)

	// Input loop.
	for true {
		text, _ := reader.ReadString('\n')
		switch mode {
		case ModeInit:
			if strings.HasPrefix(text, "start") {
				start(true)
			} else {
				signalRecieve(text)
			}
		case ModeConnect:
			signalRecieve(text)
		case ModeChat:
			sendChat(text)
			break;	
		}
		text = ""
	}
	<-wait
	fmt.Println("done")
}
