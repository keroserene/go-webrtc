package webrtc

// #include "cpeerconnection.h"
import "C"

// Bundle Policy Enum. Must have same order as in: peerconnectioninterface.h
type RTCBundlePolicy int

const (
	Balanced RTCBundlePolicy = iota
	MaxCompat
	MaxBundle
)

type RTCConfiguration struct {
	// TODO: Implement, and provide as argument to CreatePeerConnection
	IceServers           []string
	IceTransportPolicy   string
	BundlePolicy         RTCBundlePolicy
	RtcpMuxPolicy        string
	PeerIdentity         string   // Target peer identity
	Certificates         []string // TODO: implement to allow key continuity
	IceCandidatePoolSize int

	cgoConfig *C.CGORTCConfiguration // Native code internals
}

// Create a new RTCConfiguration with spec default values.
func NewRTCConfiguration() *RTCConfiguration {
	c := new(RTCConfiguration)
	c.IceServers = make([]string, 0)
	c.IceServers = nil
	c.IceTransportPolicy = "all"
	c.BundlePolicy = Balanced
	c.RtcpMuxPolicy = "require"
	c.Certificates = make([]string, 0)
	return c
}

func (config *RTCConfiguration) CGO() C.CGORTCConfiguration {
	c := new(C.CGORTCConfiguration)
	// TODO: Fix go slices to C arrays conversion
	// c.IceServers = (C.CGOArray)(unsafe.Pointer(&config.IceServers[0]))
	c.IceTransportPolicy = C.CString(config.IceTransportPolicy)
	// c.BundlePolicy = C.CString(config.BundlePolicy)
	c.BundlePolicy = C.int(config.BundlePolicy)
	c.RtcpMuxPolicy = C.CString(config.RtcpMuxPolicy)
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

type RTCIceTransportPolicy struct {
}
