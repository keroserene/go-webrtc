package webrtc

// #include "cpeerconnection.h"
// #include "ctestenums.h"
import "C"

// Spec: https://w3c.github.io/webrtc-pc/#configuration

// Bundle Policy Enum
type RTCBundlePolicy int
type RTCIceTransportPolicy int
type RTCRtcpMuxPolicy int

// These "Enum" consts must match order in: peerconnectioninterface.h
// There doesn't seem to be a way to have a named container for enums
// in go, and the idiomatic way seems to be just prefixes.
const (
	BundlePolicyBalanced RTCBundlePolicy = iota
	BundlePolicyMaxBundle
	BundlePolicyMaxCompat
)

const (
	IceTransportPolicyNone RTCIceTransportPolicy = iota
	IceTransportPolicyRelay
	// TODO: Look into why nohost is not exposed in w3c spec, but is available
	// in native code? If it does need to be exposed, capitalize the i.
	// (It still needs to exist, to ensure the enum values match up.
	iceTransportPolicyNoHost
	IceTransportPolicyAll
)

const (
	RtcpMuxPolicyNegotiate RTCRtcpMuxPolicy = iota
	RtcpMuxPolicyRequire
)

type RTCConfiguration struct {
	// TODO: Implement, and provide as argument to CreatePeerConnection
	IceServers           []string
	IceTransportPolicy   RTCIceTransportPolicy
	BundlePolicy         RTCBundlePolicy
	RtcpMuxPolicy        RTCRtcpMuxPolicy
	PeerIdentity         string   // Target peer identity
	Certificates         []string // TODO: implement to allow key continuity
	IceCandidatePoolSize int

	cgoConfig *C.CGORTCConfiguration // Native code internals
}

// Create a new RTCConfiguration with default values according to spec.
func NewRTCConfiguration() *RTCConfiguration {
	c := new(RTCConfiguration)
	c.IceServers = make([]string, 0)
	c.IceServers = nil
	c.IceTransportPolicy = IceTransportPolicyAll
	c.BundlePolicy = BundlePolicyBalanced
	c.RtcpMuxPolicy = RtcpMuxPolicyRequire
	c.Certificates = make([]string, 0)
	return c
}

func (config *RTCConfiguration) CGO() C.CGORTCConfiguration {
	c := new(C.CGORTCConfiguration)
	// TODO: Fix go slices to C arrays conversion
	// c.IceServers = (C.CGOArray)(unsafe.Pointer(&config.IceServers[0]))
	c.IceTransportPolicy = C.int(config.IceTransportPolicy)
	// c.BundlePolicy = C.CString(config.BundlePolicy)
	c.BundlePolicy = C.int(config.BundlePolicy)
	c.RtcpMuxPolicy = C.int(config.RtcpMuxPolicy)
	c.PeerIdentity = C.CString(config.PeerIdentity)
	// c.Certificates = config.Certificates
	c.IceCandidatePoolSize = C.int(config.IceCandidatePoolSize)
	config.cgoConfig = c
	return *c
}

type RTCSignalingState int

/*
const {
  stable RTCSignallingState = iota
  have-local-offer
  have-remote-offer
  have-local-pranswer
  have-remote-pranswer
  closed
}
*/

type RTCIceCredentialType struct {
}

type RTCIceServer struct {
	Urls       string
	Username   string
	Credential string
	// credentialType   RTCIceCredentialType
}

//
// Below are Go wrappers around intermediary C externs that extract the integer value of enums
// declared in native webrtc. This allows testing that the Go enums are correct.
// They unfortunately cannot be directly applied to the consts above.
//

var _cgoIceTransportPolicyNone = int(C.CGOIceTransportPolicyNone)
var _cgoIceTransportPolicyRelay = int(C.CGOIceTransportPolicyRelay)
var _cgoIceTransportPolicyNoHost = int(C.CGOIceTransportPolicyNoHost)
var _cgoIceTransportPolicyAll = int(C.CGOIceTransportPolicyAll)

var _cgoBundlePolicyBalanced = int(C.CGOBundlePolicyBalanced)
var _cgoBundlePolicyMaxCompat = int(C.CGOBundlePolicyMaxCompat)
var _cgoBundlePolicyMaxBundle = int(C.CGOBundlePolicyMaxBundle)

var _cgoRtcpMuxPolicyNegotiate = int(C.CGORtcpMuxPolicyNegotiate)
var _cgoRtcpMuxPolicyRequire = int(C.CGORtcpMuxPolicyRequire)
