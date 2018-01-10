package webrtc

import (
	"encoding/json"
)

// See: https://w3c.github.io/webrtc-pc/#idl-def-RTCIceCandidate

type (
	// IceProtocol is protocol
	IceProtocol int
	// IceCandidateType is candidate type
	IceCandidateType int
	// IceTCPCandidateTyp is tcp candidate type
	IceTCPCandidateTyp int
)

const (
	// IceProtocolUDP is UDP ice protocol
	IceProtocolUDP IceProtocol = iota
	// IceProtocolTCP is TCP ice protocol
	IceProtocolTCP
)

// String is to convert IceProtocol to corresponding string
func (p IceProtocol) String() string {
	return EnumToStringSafe(int(p), []string{
		"udp",
		"tcp",
	})
}

const (
	// IceCandidateTypeHost means ice candidate type is host
	IceCandidateTypeHost IceCandidateType = iota
	// IceCandidateTypeSrflx means ice candidate type is srflx
	IceCandidateTypeSrflx
	// IceCandidateTypePrflx means ice candidate type is prflx
	IceCandidateTypePrflx
	// IceCandidateTypeRelay means ice candidate type is relay
	IceCandidateTypeRelay
)

// String is to convert IceCandidateType to corresponding string
func (t IceCandidateType) String() string {
	return EnumToStringSafe(int(t), []string{
		"host",
		"srflx",
		"prflx",
		"relay",
	})
}

const (
	// IceTCPCandidateTypActive means active TCP candidate
	IceTCPCandidateTypActive IceTCPCandidateTyp = iota
	// IceTCPCandidateTypPassive means passive TCP candidate
	IceTCPCandidateTypPassive
	// IceTCPCandidateTypSo means so TCP candidate
	IceTCPCandidateTypSo
)

// String is to convert IceTCPCandidateTyp to corresponding string
func (t IceTCPCandidateTyp) String() string {
	return EnumToStringSafe(int(t), []string{
		"active",
		"passive",
		"so",
	})
}

// IceCandidate is structure to store candidate related parameter
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
	// TcpType        IceTCPCandidateTyp
	// RelatedAddress string
	// RelatedPort    C.ushort
}

// Serialize serialize an IceCandidate into a JSON string.
func (candidate *IceCandidate) Serialize() string {
	bytes, err := json.Marshal(candidate)
	if nil != err {
		ERROR.Println(err)
		return ""
	}
	return string(bytes)
}

// DeserializeIceCandidate deserialize a received json string into an IceCandidate, if possible.
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
