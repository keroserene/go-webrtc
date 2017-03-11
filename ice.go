package webrtc

import (
	"encoding/json"
)

// See: https://w3c.github.io/webrtc-pc/#idl-def-RTCIceCandidate

type (
	IceProtocol         int
	IceCandidateType    int
	IceTcpCandidateType int
)

const (
	IceProtocolUPD IceProtocol = iota
	IceProtocolTCP
)

func (p IceProtocol) String() string {
	return EnumToStringSafe(int(p), []string{
		"udp",
		"tcp",
	})
}

const (
	IceCandidateTypeHost IceCandidateType = iota
	IceCandidateTypeSrflx
	IceCandidateTypePrflx
	IceCandidateTypeRelay
)

func (t IceCandidateType) String() string {
	return EnumToStringSafe(int(t), []string{
		"host",
		"srflx",
		"prflx",
		"relay",
	})
}

const (
	IceTcpCandidateTypeActive IceTcpCandidateType = iota
	IceTcpCandidateTypePassive
	IceTcpCandidateTypeSo
)

func (t IceTcpCandidateType) String() string {
	return EnumToStringSafe(int(t), []string{
		"active",
		"passive",
		"so",
	})
}

type IceCandidate struct {
	Candidate     string `json:"candidate"`
	SdpMid        string `json:"sdpMid"`
	SdpMLineIndex int    `json:"sdpMLineIndex"`
	// Foundation     string
	// Priority       C.ulong
	// IP             net.IP
	// Protocol       IceProtocol
	// Port           C.ushort
	// Type           IceCandidateType
	// TcpType        IceTcpCandidateType
	// RelatedAddress string
	// RelatedPort    C.ushort
}

// Serialize an IceCandidate into a JSON string.
func (candidate *IceCandidate) Serialize() string {
	bytes, err := json.Marshal(candidate)
	if nil != err {
		ERROR.Println(err)
		return ""
	}
	return string(bytes)
}

// Deserialize a received json string into an IceCandidate, if possible.
func DeserializeIceCandidate(msg string) *IceCandidate {
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(msg), &parsed)
	if nil != err {
		ERROR.Println(err)
		return nil
	}
	if _, ok := parsed["candidate"]; !ok {
		ERROR.Println("Cannot deserialize IceCandidate without candidate field.")
		return nil
	}
	if _, ok := parsed["sdpMid"]; !ok {
		ERROR.Println("Cannot deserialize IceCandidate without sdpMid field.")
		return nil
	}
	if _, ok := parsed["sdpMLineIndex"]; !ok {
		ERROR.Println("Cannot deserialize IceCandidate without sdpMLineIndex field.")
		return nil
	}
	ice := new(IceCandidate)
	ice.Candidate = parsed["candidate"].(string)
	ice.SdpMid = parsed["sdpMid"].(string)
	ice.SdpMLineIndex = int(parsed["sdpMLineIndex"].(float64))
	return ice
}
