package webrtc

// #include "cpeerconnection.h"
// #include <stdlib.h>  // Needed for C.free
import "C"
import (
	"encoding/json"
	"unsafe"
)

/* WebRTC SessionDescription

See: https://w3c.github.io/webrtc-pc/#idl-def-RTCSessionDescription
*/
type SessionDescription struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`

	// Keep track of internal C++ *webrtc::SessionDescriptionInterface.
	cgoSdp C.CGO_sdp
}

// TODO: Turn into Enum.
var SdpTypes = []string{"offer", "pranswer", "answer", "rollback"}

// Construct a SessionDescription object from a valid msg.
func NewSessionDescription(sdpType string, msg string) *SessionDescription {
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
	s := C.CString(sdpType)
	defer C.free(unsafe.Pointer(s))
	m := C.CString(msg)
	defer C.free(unsafe.Pointer(m))
	cSdp := C.CGO_DeserializeSDP(s, m)
	if nil == cSdp {
		ERROR.Println("Invalid SDP string.")
		return nil
	}
	sdp := new(SessionDescription)
	sdp.cgoSdp = cSdp
	sdp.Type = sdpType
	sdp.Sdp = msg
	return sdp
}

// Serialize a SessionDescription into a JSON string.
func (desc *SessionDescription) Serialize() string {
	bytes, err := json.Marshal(desc)
	if nil != err {
		ERROR.Println(err)
		return ""
	}
	return string(bytes)
}

// Deserialize a received json string into a SessionDescription, if possible.
func DeserializeSessionDescription(msg string) *SessionDescription {
	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(msg), &parsed)
	if nil != err {
		ERROR.Println(err)
		return nil
	}
	if _, ok := parsed["type"]; !ok {
		ERROR.Println("Cannot deserialize SessionDescription without type field.")
		return nil
	}
	if _, ok := parsed["sdp"]; !ok {
		ERROR.Println("Cannot deserialize SessionDescription without sdp field.")
		return nil
	}
	return NewSessionDescription(parsed["type"].(string), parsed["sdp"].(string))
}
