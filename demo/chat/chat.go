/*
 * Webrtc chat demo.
 * Send chat messages via webrtc, over go.
 * Can interop with the JS client. (Open chat.html in a browser)
 *
 * To use: `go run chat.go`
 */
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/keroserene/go-webrtc"
)

var pc *webrtc.PeerConnection
var dc *webrtc.DataChannel
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

//
// Preparing SDP messages for signaling.
// generateOffer and generateAnswer are expected to be called within goroutines.
// It is possible to send the serialized offers or answers immediately upon
// creation, followed by subsequent individual ICE candidates.
//
// However, to ease the user's copy & paste experience, in this case we forgo
// the trickle ICE and wait for OnIceComplete to fire, which will contain
// a full SDP mesasge with all ICE candidates, so the user only has to copy
// one message.
//

func generateOffer() {
	fmt.Println("Generating offer...")
	offer, err := pc.CreateOffer() // blocking
	if err != nil {
		fmt.Println(err)
		return
	}
	pc.SetLocalDescription(offer)
}

func generateAnswer() {
	fmt.Println("Generating answer...")
	answer, err := pc.CreateAnswer() // blocking
	if err != nil {
		fmt.Println(err)
		return
	}
	pc.SetLocalDescription(answer)
}

func receiveDescription(sdp *webrtc.SessionDescription) {
	err = pc.SetRemoteDescription(sdp)
	if nil != err {
		fmt.Println("ERROR", err)
		return
	}
	fmt.Println("SDP " + sdp.Type + " successfully received.")
	if "offer" == sdp.Type {
		go generateAnswer()
	}
}

// Manual "copy-paste" signaling channel.
func signalSend(msg string) {
	fmt.Println("\n ---- Please copy the below to peer ---- \n")
	fmt.Println(msg + "\n")
}

func signalReceive(msg string) {
	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(msg), &parsed)
	if nil != err {
		// fmt.Println(err, ", try again.")
		return
	}

	// If this is a valid signal and no PeerConnection has been instantiated,
	// start as the "answerer."
	if nil == pc {
		start(false)
	}

	if nil != parsed["sdp"] {
		sdp := webrtc.DeserializeSessionDescription(msg)
		if nil == sdp {
			fmt.Println("Invalid SDP.")
			return
		}
		receiveDescription(sdp)
	}

	// Allow individual ICE candidate messages, but this won't be necessary if
	// the remote peer also doesn't use trickle ICE.
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

// Attach callbacks to a newly created data channel.
// In this demo, only one data channel is expected, and is only used for chat.
// But it is possible to send any sort of bytes over a data channel, for many
// more interesting purposes.
func prepareDataChannel(channel *webrtc.DataChannel) {
	channel.OnOpen = func() {
		fmt.Println("Data Channel Opened!")
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

func startChat() {
	mode = ModeChat
	fmt.Println("------- chat enabled! -------")
}

func endChat() {
	mode = ModeInit
	fmt.Println("------- chat disabled -------")
}

func sendChat(msg string) {
	line := username + ": " + msg
	fmt.Println("[sent]")
	dc.SendText(line)
}

func receiveChat(msg string) {
	fmt.Println("\n" + string(msg))
}

// Janky /command inputs.
func parseCommands(input string) bool {
	if !strings.HasPrefix(input, "/") {
		return false
	}
	cmd := strings.TrimSpace(strings.TrimLeft(input, "/"))
	switch cmd {
	case "quit":
		fmt.Println("Disconnecting chat session...")
		dc.Close()
	case "status":
		fmt.Println("WebRTC PeerConnection Configuration:\n", pc.GetConfiguration())
		fmt.Println("Signaling State: ", pc.SignalingState())
		fmt.Println("Connection State: ", pc.ConnectionState())
	case "help":
		showCommands()
	default:
		fmt.Println("Unknown command:", cmd)
		showCommands()
	}
	return true
}

func showCommands() {
	fmt.Println("Possible commands: help status quit")
}

// Create a PeerConnection.
// If |instigator| is true, create local data channel which causes a
// negotiation-needed, leading to preparing an SDP offer to be sent to the
// remote peer. Otherwise, await an SDP offer from the remote peer, and send an
// answer back.
func start(instigator bool) {
	mode = ModeConnect
	fmt.Println("Starting up PeerConnection...")
	// TODO: Try with TURN servers.
	config := webrtc.NewConfiguration(
		webrtc.OptionIceServer("stun:stun.l.google.com:19302"))

	pc, err = webrtc.NewPeerConnection(config)
	if nil != err {
		fmt.Println("Failed to create PeerConnection.")
		return
	}

	// OnNegotiationNeeded is triggered when something important has occurred in
	// the state of PeerConnection (such as creating a new data channel), in which
	// case a new SDP offer must be prepared and sent to the remote peer.
	pc.OnNegotiationNeeded = func() {
		go generateOffer()
	}
	// Once all ICE candidates are prepared, they need to be sent to the remote
	// peer which will attempt reaching the local peer through NATs.
	pc.OnIceComplete = func() {
		fmt.Println("Finished gathering ICE candidates.")
		sdp := pc.LocalDescription().Serialize()
		signalSend(sdp)
	}
	/*
		pc.OnIceGatheringStateChange = func(state webrtc.IceGatheringState) {
			fmt.Println("Ice Gathering State:", state)
			if webrtc.IceGatheringStateComplete == state {
				// send local description.
			}
		}
	*/
	// A DataChannel is generated through this callback only when the remote peer
	// has initiated the creation of the data channel.
	pc.OnDataChannel = func(channel *webrtc.DataChannel) {
		fmt.Println("Datachannel established by remote... ", channel.Label())
		dc = channel
		prepareDataChannel(channel)
	}

	if instigator {
		// Attempting to create the first datachannel triggers ICE.
		fmt.Println("Initializing datachannel....")
		dc, err = pc.CreateDataChannel("test")
		if nil != err {
			fmt.Println("Unexpected failure creating Channel.")
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Println("Demo interrupted. Disconnecting...")
		if nil != dc {
			dc.Close()
		}
		if nil != pc {
			pc.Close()
		}
		os.Exit(1)
	}()

	// Input loop.
	for {
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
			// TODO: make chat interface nicer.
			if !parseCommands(text) {
				sendChat(text)
			}
			// fmt.Print(username + ": ")
			break
		}
	}
	<-wait
	fmt.Println("done")
}
