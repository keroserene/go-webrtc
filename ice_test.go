package webrtc

import (
	"testing"
)

func TestSerializeIceCandidate(t *testing.T) {
	ice := IceCandidate{
		"fake",
		"not real",
		1337,
	}
	s := ice.Serialize()
	expected := `{"candidate":"fake","sdpMid":"not real","sdpMLineIndex":1337}`
	if expected != s {
		t.Error("Incorrect Serializing of SessionDescription", s)
	}
}

func TestDeserializeIceCandidate(t *testing.T) {
	r := DeserializeIceCandidate(`
		{"candidate":"still fake","sdpMid":"illusory","sdpMLineIndex":1337}`)
	if nil == r {
		t.Fatal("Failed to deserialize IceCandidate.")
	}
	if "still fake" != r.Candidate {
		t.Error("Unexpected candidate:", r.Candidate)
	}
	if "illusory" != r.SdpMid {
		t.Error("Unexpected sdpMid:", r.SdpMid)
	}
	if 1337 != r.SdpMLineIndex {
		t.Error("Unexpected sdpMLineIndex:", r.SdpMLineIndex)
	}
}

func TestRoundtripSerializeDeserializeICE(t *testing.T) {
	ice := IceCandidate{
		"totally fake",
		"fabricated",
		1337,
	}
	r := DeserializeIceCandidate(ice.Serialize())
	if r.Candidate != ice.Candidate || r.SdpMid != ice.SdpMid ||
		r.SdpMLineIndex != ice.SdpMLineIndex {
		t.Error("Incorrect roundtrip serialize and deserialization.")
	}
}
