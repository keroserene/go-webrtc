package webrtc

import (
	"fmt"
	"github.com/keroserene/go-webrtc/data"
	"testing"
	"time"
	"unsafe"
)

// TODO: Try Gucumber or some potential fancy test framework for go.
/*
func checkEnum(t *testing.T, desc string, enum int, expected int) {
	if enum != expected {
		t.Error("Mismatched Enum Value -", desc,
			"\nwas:", enum,
			"\nexpected:", expected)
	}
}

func TestPeerConnectionStateEnums(t *testing.T) {
	checkEnum(t, "PeerConnectionStateNew",
		int(PeerConnectionStateNewd), _cgoBundlePolicyBalanced)
	checkEnum(t, "BundlePolicyMaxCompat",
		int(BundlePolicyMaxCompat), _cgoBundlePolicyMaxCompat)
	checkEnum(t, "BundlePolicyMaxBundle",
		int(BundlePolicyMaxBundle), _cgoBundlePolicyMaxBundle)
} */

// These tests create two PeerConnections objects, which allows a loopback test.
var pcA *PeerConnection
var pcB *PeerConnection
var err error
var sdp *SessionDescription
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
	fmt.Println("SDP Offer:\n", sdp.Sdp)
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
	t.SkipNow() // Can't import C in tests
	success := make(chan IceCandidate, 1)
	pcA.OnIceCandidate = func(ic IceCandidate) {
		success <- ic
	}
	// cgoOnIceCandidate(unsafe.Pointer(pcA), C.CString("not real"), ...)
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
	fmt.Print("\n ~~ Signalling Happens here ~~ \n\n")
}

func TestSetRemoteDescription(t *testing.T) {
	fmt.Println("\n == BOB's PeerConnection ==")
	err = pcB.SetRemoteDescription(sdp)
	if nil != err {
		t.Fatal(err)
	}
}

func TestAddIceCandidate(t *testing.T) {
	t.SkipNow() // Needs real data
	ic := IceCandidate{"fixme", "", 0}
	err = pcB.AddIceCandidate(ic)
	if err != nil {
		t.Fatal(err)
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
	fmt.Println("SDP Answer:\n", sdp.Sdp)
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

func TestOnConnectionStateChangeCallback(t *testing.T) {
	success := make(chan PeerConnectionState, 1)
	pcA.OnConnectionStateChange = func(state PeerConnectionState) {
		success <- state
	}
	cgoOnConnectionStateChange(unsafe.Pointer(pcA),
		PeerConnectionStateDisconnected)
	select {
	case r := <-success:
		if PeerConnectionStateDisconnected != r {
			t.Error("Unexpected PeerConnectionState:", r)
		}
	case <-time.After(time.Second * 1):
		t.Fatal("Timed out.")
	}
}

func TestCreateDataChannel(t *testing.T) {
	channel, err := pcA.CreateDataChannel("test", data.Init{})
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println("Data channel: ", channel)
	label := channel.Label()
	if label != "test" {
		t.Error("Unexpected label:", label)
	}
	channel.Close()
}

func TestOnDataChannelCallback(t *testing.T) {
	success := make(chan *data.Channel, 1)
	pcA.OnDataChannel = func(dc *data.Channel) {
		success <- dc
	}
	cgoOnDataChannel(unsafe.Pointer(pcA), nil)
	select {
	case <-success:
	case <-time.After(time.Second * 1):
		t.Fatal("Timed out.")
	}
}

func TestClose(t *testing.T) {
	success := make(chan int, 1)
	go func() {
		// pcA.Close()
		pcB.Close()
		success <- 1
	}()
	// TODO: Check the signaling state.
	select {
	case <-success:
	case <-time.After(time.Second * 2):
		WARN.Println("Timed out... something's probably amiss.")
		success <- 0
	}
}

// TODO: tests for video / audio stream support.
