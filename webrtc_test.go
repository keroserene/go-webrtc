package webrtc

import (
	"fmt"
	"github.com/keroserene/go-webrtc/datachannel"
	"testing"
	"time"
	"unsafe"
)

// TODO: Try Gucumber or some potential fancy test framework.

// These tests create two PeerConnections objects, which allows a loopback test.
var pcA *PeerConnection
var pcB *PeerConnection
var err error
var sdp *SDPHeader
var config *Configuration

func TestCreatePeerConnection(t *testing.T) {
	config = NewConfiguration()
	if nil == config {
		t.Fatal("Unable to create Configuration")
	}
	pcA, err = NewPeerConnection(config)
	if nil != err {
		t.Fatal(err)
	}
}

func TestCreateSecondPeerConnections(t *testing.T) {
	pcB, err = NewPeerConnection(config)
	if nil != err {
		t.Fatal(err)
	}
}

func TestCreateOffer(t *testing.T) {
	fmt.Println("\n== ALICE's PeerConnection ==")
	sdp, err = pcA.CreateOffer()
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println("SDP Offer:\n", sdp.description)
}

func TestOnSignalingStateChangeCallback(t *testing.T) {
	success := make(chan SignalingState, 1)
	pcA.OnSignalingStateChange = func(s SignalingState) {
		success <- s
	}
	cgoOnSignalingStateChange(unsafe.Pointer(pcA), SignalingStateStable)
	select {
	case state := <-success:
		if SignalingStateStable != state {
			t.Error("Unexpected SignalingState:", state)
		}
	case <-time.After(time.Second * 1):
		t.Fatal("Timed out.")
	}
}

func TestOnIceCandidateCallback(t *testing.T) {
	success := make(chan string, 1)
	pcA.OnIceCandidate = func(c string) {
		success <- c
	}
	// candidate := "not a real ICE candidate";
	cgoOnIceCandidate(unsafe.Pointer(pcA), nil)
	select {
	case <-success:
	case <-time.After(time.Second * 1):
		t.Fatal("Timed out.")
	}
}

func TestSetLocalDescription(t *testing.T) {
	err = pcA.SetLocalDescription(sdp)
	if nil != err {
		t.Fatal(err)
	}

	// Pretend pcA sends the SDP offer to pcB through some signalling channel.
	fmt.Println("\n ~~ Signalling Happens here ~~ \n")
}

func TestSetRemoteDescription(t *testing.T) {
	fmt.Println("\n == BOB's PeerConnection ==")
	err = pcB.SetRemoteDescription(sdp)
	if nil != err {
		t.Fatal(err)
	}
}

func TestAddIceCandidate(t *testing.T) {
	err := pcB.AddIceCandidate("not real")
	// Expected to fail because the ICE candidate is fake.
	if err == nil {
		// TODO: Change this test once a non-stringified version of IceCandidates
		// is implemented.
		t.Error("AddIceCandidate was expecting to fail.")
	}
}

func TestGetSignalingState(t *testing.T) {
	state := pcB.SignalingState()
	if SignalingStateHaveRemoteOffer != state {
		t.Error("Unexected signaling state:", state)
	}
	fmt.Println(SignalingStateString[state])
}

func TestSetAndGetConfiguration(t *testing.T) {
	config := NewConfiguration(
		OptionIceServer("stun:something.else"),
		OptionIceTransportPolicy(IceTransportPolicyRelay))
	pcA.SetConfiguration(*config)
	got := pcA.GetConfiguration()
	if got.IceTransportPolicy != IceTransportPolicyRelay {
		t.Error("Unexpected Configuration: ",
			IceTransportPolicyString[got.IceTransportPolicy])
	}
}

func TestCreateAnswer(t *testing.T) {
	sdp, err := pcB.CreateAnswer()
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println("SDP Answer:\n", sdp.description)
}

func TestOnNegotiationNeededCallback(t *testing.T) {
	success := make(chan int, 1)
	pcA.OnNegotiationNeeded = func() {
		success <- 0
	}
	cgoOnNegotiationNeeded(unsafe.Pointer(pcA))
	select {
	case <-success:
	case <-time.After(time.Second * 1):
		t.Fatal("Timed out.")
	}
}

// TODO: real datachannel tests
func TestCreateDataChannel(t *testing.T) {
	channel, err := pcA.CreateDataChannel("test", datachannel.Init{})
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println("Data channel: ", channel)
}

func TestOnDataChannelCallback(t *testing.T) {
	success := make(chan string, 1)
	pcA.OnDataChannel = func(channel string) {
		success <- c
	}
	cgoOnDataChannel(unsafe.Pointer(pcA), "")
	select {
	case <-success:
	case <-time.After(time.Second * 1):
		t.Fatal("Timed out.")
	}
}

// TODO: tests for video / audio stream support.
