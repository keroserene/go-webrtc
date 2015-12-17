package webrtc

import (
	"testing"
)

func TestNewSessionDescription(t *testing.T) {
	r := NewSessionDescription("offer", "fake")
	if nil == r {
		t.Fatal("Unable to create new SessionDescription.")
	}
	if "offer" != r.Type {
		t.Error("Unexpected Type:", r.Type)
	}
	if "fake" != r.Sdp {
		t.Error("Unexpected Sdp:", r.Sdp)
	}
}

func TestSerializeSessionDescription(t *testing.T) {
	sdp := NewSessionDescription("answer", "fake")
	s := sdp.Serialize()
	expected := `{"type":"answer","sdp":"fake"}`
	if expected != s {
		t.Error("Incorrect Serializing of SessionDescription", s)
	}
}

func TestDeserializeSessionDescription(t *testing.T) {
	r := DeserializeSessionDescription(`{"type":"answer","sdp":"fake"}`)
	if nil == r {
		t.Fatal("Failed to deserialize SessionDescription.")
	}
	if "answer" != r.Type {
		t.Error("Unexpected type:", r.Type)
	}
	if "fake" != r.Sdp {
		t.Error("Unexpected sdp:", r.Sdp)
	}
}

func TestRoundtripSerializeDeserialize(t *testing.T) {
	sdp := NewSessionDescription("pranswer", "not real")
	r := DeserializeSessionDescription(sdp.Serialize())
	if r.Type != sdp.Type || r.Sdp != sdp.Sdp {
		t.Error("Incorrect roundtrip serialize and deserialization.")
	}
}
