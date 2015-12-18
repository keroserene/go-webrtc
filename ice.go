package webrtc


type (
	IceProtocol         int
	IceCandidateType    int
	IceTcpCandidateType int
)

const (
	IceProtocolUPD IceProtocol = iota
	IceProtocolTCP
)

var IceProtocolString = []string{"udp", "tcp"}

const (
	IceCandidateTypeHost IceCandidateType = iota
	IceCandidateTypeSrflx
	IceCandidateTypePrflx
	IceCandidateTypeRelay
)

var IceCandidateTypeString = []string{"host", "srflx", "prflx", "relay"}

const (
	IceTcpCandidateTypeActive IceTcpCandidateType = iota
	IceTcpCandidateTypePassive
	IceTcpCandidateTypeSo
)

var IceTcpCandidateTypeString = []string{"active", "passive", "so"}

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
https://w3c.github.io/webrtc-pc/#idl-def-RTCIceCandidate