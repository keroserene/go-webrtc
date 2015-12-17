package webrtc

import (
	"fmt"
	"testing"
)

// Ensure the Go "enums" generated in the idiomatic iota const way actually
// match up with actual int values of the underlying native WebRTC Enums.
func checkEnum(t *testing.T, desc string, enum int, expected int) {
	if enum != expected {
		t.Error("Mismatched Enum Value -", desc,
			"\nwas:", enum,
			"\nexpected:", expected)
	}
}

func TestBundlePolicyEnums(t *testing.T) {
	checkEnum(t, "BundlePolicyBalanced",
		int(BundlePolicyBalanced), _cgoBundlePolicyBalanced)
	checkEnum(t, "BundlePolicyMaxCompat",
		int(BundlePolicyMaxCompat), _cgoBundlePolicyMaxCompat)
	checkEnum(t, "BundlePolicyMaxBundle",
		int(BundlePolicyMaxBundle), _cgoBundlePolicyMaxBundle)
}

func TestIceTransportPolicyEnums(t *testing.T) {
	checkEnum(t, "IceTransportPolicyNone",
		int(IceTransportPolicyNone), _cgoIceTransportPolicyNone)
	checkEnum(t, "IceTransportPolicyRelay",
		int(IceTransportPolicyRelay), _cgoIceTransportPolicyRelay)
	checkEnum(t, "IceTransportPolicyAll",
		int(IceTransportPolicyAll), _cgoIceTransportPolicyAll)
}

// TODO: [ED]
/* func TestRtcpMuxPolicy(t *testing.T) {
	checkEnum(t, "RtcpMuxPolicyNegotiate",
		int(RtcpMuxPolicyNegotiate), _cgoRtcpMuxPolicyNegotiate)
	checkEnum(t, "RtcpMuxPolicyRequire",
		int(RtcpMuxPolicyRequire), _cgoRtcpMuxPolicyRequire)
} */

func TestSignalingStateEnums(t *testing.T) {
	checkEnum(t, "SignalingStateStable",
		int(SignalingStateStable), _cgoSignalingStateStable)
	checkEnum(t, "SignalingStateHaveLocalOffer",
		int(SignalingStateHaveLocalOffer), _cgoSignalingStateHaveLocalOffer)
	checkEnum(t, "SignalingStateHaveLocalPrAnswer",
		int(SignalingStateHaveLocalPrAnswer), _cgoSignalingStateHaveLocalPrAnswer)
	checkEnum(t, "SignalingStateHaveRemoteOffer",
		int(SignalingStateHaveRemoteOffer), _cgoSignalingStateHaveRemoteOffer)
	checkEnum(t, "SignalingStateHaveRemotePrAnswer",
		int(SignalingStateHaveRemotePrAnswer), _cgoSignalingStateHaveRemotePrAnswer)
	checkEnum(t, "SignalingStateClosed",
		int(SignalingStateClosed), _cgoSignalingStateClosed)
}

func TestIceServer(t *testing.T) {
	s, err := NewIceServer()
	if nil == err {
		t.Error("NewIceServer should have failed given 0 params",
			s.Urls)
	}
	s, err = NewIceServer("")
	if nil == err {
		t.Error("NewIceServer should have failed given empty urls.")
	}
	s, err = NewIceServer("stun:12345, badurl")
	if nil == err {
		t.Error("NewIceServer should have failed given malformed url.")
	}
	s, err = NewIceServer("stun:12345, stun:ok")
	if nil != err {
		t.Error(err)
	}
	s, err = NewIceServer("stun:a, turn:b")
	if nil != err {
		t.Error(err)
	}
	s, err = NewIceServer("stun:a, turn:b", "alice")
	if nil != err {
		t.Error(err)
	}
	s, err = NewIceServer("stun:a, turn:b", "alice", "secret")
	if nil != err {
		t.Error(err)
	}
	s, err = NewIceServer("stun:a, turn:b", "alice", "secret", "extra")
	if nil != err {
		t.Error("NewIceServer shouldn't fail, only WARN on too many params.")
	}
	fmt.Println(s)
}

func TestNewConfiguration(t *testing.T) {
	config := NewConfiguration()
	if nil == config {
		t.Error("NewConfiguration could not generate basic config.")
	}
	config = NewConfiguration(OptionIceServer("stun:a"))
	if len(config.IceServers) != 1 {
		t.Error("NewConfiguration should have 1 ICE server.")
	}
	config = NewConfiguration(
		OptionIceServer("stun:a"),
		OptionIceServer("stun:b, turn:c"))
	if len(config.IceServers) != 2 {
		t.Error("NewConfiguration should have 2 ICE servers.")
	}

	config = NewConfiguration(
		OptionIceServer("stun:d"),
		OptionIceTransportPolicy(IceTransportPolicyAll))
	if IceTransportPolicyAll != config.IceTransportPolicy {
		t.Error("OptionIceTransportPolicy failed, was ", config.IceTransportPolicy)
	}

	config = NewConfiguration(
		OptionIceServer("stun:d"),
		OptionBundlePolicy(BundlePolicyMaxCompat))
	if BundlePolicyMaxCompat != config.BundlePolicy {
		t.Error("OptionBundlePolicy failed, was ", config.BundlePolicy)
	}
}

func TestIceServerCGO(t *testing.T) {
	// TODO
}

func TestConfigurationCGO(t *testing.T) {
	// TODO
}
