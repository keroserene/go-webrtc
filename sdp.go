package webrtc

// #include "peerconnection.h"
// #include <stdlib.h>  // Needed for C.free
import "C"
import (
	"encoding/json"
	"unsafe"
)

// SessionDescription is structure of SDP parameters
// https://w3c.github.io/webrtc-pc/#idl-def-RTCSessionDescription
type SessionDescription struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

// SdpTypes is type of Sdp
// TODO: Turn into Enum.
var SdpTypes = []string{"offer", "pranswer", "answer", "rollback"}

// CgoSdpToGoString convert SDP data to string
func CgoSdpToGoString(sdp C.CGO_sdp) string {
	serializedSDP := C.CGO_SerializeSDP(sdp)
	defer C.free(unsafe.Pointer(serializedSDP))
	return C.GoString(serializedSDP)
}

// NewSessionDescription construct a SessionDescription object from a valid msg.
func NewSessionDescription(sdpType string, serializedSDP C.CGO_sdpString) *SessionDescription {
	defer C.free(unsafe.Pointer(serializedSDP))
	in := false
	for i := 0; i < len(SdpTypes); i++ {
		if SdpTypes[i] == sdpType {
			in = true
		}
	}
	if !in {
		ERROR.Println("Invalid SDP type.")
		return nil
	}
	sdp := new(SessionDescription)
	sdp.Type = sdpType
	sdp.Sdp = C.GoString(serializedSDP)
	return sdp
}

// Serialize serialize a SessionDescription into a JSON string.
func (desc *SessionDescription) Serialize() string {
	bytes, err := json.Marshal(desc)
	if nil != err {
		ERROR.Println(err)
		return ""
	}
	return string(bytes)
}

// GoStringToCgoSdp convert string to C
func (desc *SessionDescription) GoStringToCgoSdp() C.CGO_sdp {
	t := C.CString(desc.Type)
	defer C.free(unsafe.Pointer(t))
	s := C.CString(desc.Sdp)
	defer C.free(unsafe.Pointer(s))
	return C.CGO_DeserializeSDP(t, s)
}

// DeserializeSessionDescription deserialize a received json string into a SessionDescription, if possible.
func DeserializeSessionDescription(msg string) *SessionDescription {
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(msg), &parsed)
	if nil != err {
		ERROR.Println(err)
		return nil
	}
	if _, ok := parsed["type"]; !ok {
		ERROR.Println("cannot deserialize SessionDescription without type field")
		return nil
	}
	if _, ok := parsed["sdp"]; !ok {
		ERROR.Println("cannot deserialize SessionDescription without sdp field")
		return nil
	}
	return &SessionDescription{
		Type: parsed["type"].(string),
		Sdp:  parsed["sdp"].(string),
	}
}
