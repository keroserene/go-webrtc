package webrtc

// #include "cpeerconnection.h"
// #include "ctestenums.h"
import "C"
import (
	"errors"
	"fmt"
	"strings"
	// "encoding/json"
)

// Working draft spec: http://www.w3.org/TR/webrtc/#idl-def-RTCConfiguration
// There are quite a few differences in the latest Editor's draft, but
// for now they are omitted from this Go interface, or commented out with
// an [ED] above.)
// See https://w3c.github.io/webrtc-pc/#idl-def-RTCConfiguration

type (
	RTCBundlePolicy       int
	RTCIceTransportPolicy int
	RTCRtcpMuxPolicy      int
	RTCIceCredentialType  int
	RTCSignalingState     int
)

type RTCConfiguration struct {
	// TODO: Implement, and provide as argument to CreatePeerConnection
	IceServers         []RTCIceServer
	IceTransportPolicy RTCIceTransportPolicy
	BundlePolicy       RTCBundlePolicy
	// [ED] RtcpMuxPolicy        RTCRtcpMuxPolicy
	PeerIdentity string // Target peer identity

	// This would allow key continuity.
	// [ED] Certificates         []string
	// [ED] IceCandidatePoolSize int

	cgoConfig *C.CGORTCConfiguration // Native code internals
}

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
	SignalingStateStable RTCSignalingState = iota
	SignalingStateHaveLocalOffer
	SignalingStateHaveLocalPrAnswer
	SignalingStateHaveRemoteOffer
	SignalingStateHaveRemotePrAnswer
	SignalingStateClosed
)

var RTCSignalingStateString = []string{"Stable",
	"HaveLocalOffer", "HaveLocalPrAnswer",
	"HaveRemoteOffer", "HaveRemotePrAnswer",
	"Closed" }

// TODO: [ED]
/* const (
	RtcpMuxPolicyNegotiate RTCRtcpMuxPolicy = iota
	RtcpMuxPolicyRequire
) */

// TODO: [ED]
/* const (
	IceCredentialTypePassword RTCIceCredentialType = iota
	IceCredentialTypeToken
) */

type RTCIceServer struct {
	Urls       []string // The only "required" element.
	Username   string
	Credential string
	// [ED] CredentialType RTCIceCredentialType
}

// Create a new IceServer object.
// Expects anywhere from one to three strings, in this order:
// - comma-separated list of urls.
// - username
// - credential
// TODO: For the ED version, may need to support CredentialType.
func NewIceServer(params ...string) (*RTCIceServer, error) {
	if len(params) < 1 {
		return nil, errors.New("IceServer: missing first comma-separated Urls string.")
	}
	if len(params) > 3 {
		WARN.Printf("IceServer: got %d strings, expect <= 3. Ignoring extras.\n",
			len(params))
	}
	urls := strings.Split(params[0], ",")
	username := ""
	credential := ""
	if 0 == len(urls) {
		return nil, errors.New("IceServer: requires at least one Url")
	}
	for i, url := range urls {
		url = strings.TrimSpace(url)
		// TODO: Better url validation.
		if !strings.HasPrefix(url, "stun:") &&
			!strings.HasPrefix(url, "turn:") {
			msg := fmt.Sprintf("IceServer: received malformed url: <%s>", url)
			ERROR.Println(msg)
			return nil, errors.New(msg)
		}
		urls[i] = url
	}
	if len(params) > 1 {
		username = params[1]
	}
	if len(params) > 2 {
		credential = params[2]
	}
	return &RTCIceServer{
		Urls:       urls,
		Username:   username,
		Credential: credential,
	}, nil
}

// Create a new RTCConfiguration with default values according to spec.
// Accepts any number of |RTCIceServer|s.
// Returns nil if there's an error.
func NewRTCConfiguration(options ...RTCConfigurationOption) *RTCConfiguration {
	c := new(RTCConfiguration)
	c.IceTransportPolicy = IceTransportPolicyAll
	c.BundlePolicy = BundlePolicyBalanced
	for _, op := range options {
		err := op(c)
		if nil != err {
			ERROR.Println(err)
		}
	}
	// [ED] c.RtcpMuxPolicy = RtcpMuxPolicyRequire
	// [ED] c.Certificates = make([]string, 0)
	INFO.Println("Created RTCConfiguration at ", c)
	INFO.Println("# IceServers: ", len(c.IceServers))
	// TODO: Determine whether the below is true.
	// if 0 == len(c.IceServers) {
	// ERROR.Println("Need at least one ICE server.")
	// return nil
	// }
	return c
}

// Used in RTCConfiguration's variadic functional constructor
type RTCConfigurationOption func(c *RTCConfiguration) error

func OptionIceServer(params ...string) RTCConfigurationOption {
	return func(config *RTCConfiguration) error {
		return config.AddIceServer(params...)
	}
}

func OptionIceTransportPolicy(policy RTCIceTransportPolicy) RTCConfigurationOption {
	return func(config *RTCConfiguration) error {
		config.IceTransportPolicy = policy
		return nil
	}
}

func OptionBundlePolicy(policy RTCBundlePolicy) RTCConfigurationOption {
	return func(config *RTCConfiguration) error {
		config.BundlePolicy = policy
		return nil
	}
}

func (config *RTCConfiguration) AddIceServer(params ...string) error {
	server, err := NewIceServer(params...)
	if nil != err {
		return err
	}
	config.IceServers = append(config.IceServers, *server)
	return nil
}

// Helpers which prepare Go-side of cast to eventual C++ RTCConfiguration struct.
func (server *RTCIceServer) CGO() C.CGOIceServer {
	cServer := new(C.CGOIceServer)
	cServer.numUrls = C.int(len(server.Urls))
	total := C.int(len(server.Urls))
	// TODO: Make this conversion nicer.
	cUrls := make([](*C.char), total)
	for i, url := range server.Urls {
		cUrls[i] = C.CString(url)
	}
	cServer.urls = &cUrls[0]
	cServer.numUrls = total
	cServer.username = C.CString(server.Username)
	cServer.credential = C.CString(server.Credential)
	return *cServer
}

// The C side of things will still need to allocate memory, due to the slices.
// Assumes RTCConfiguration is valid.
func (config *RTCConfiguration) CGO() C.CGORTCConfiguration {
	INFO.Println("Converting Config: ", config)
	c := new(C.CGORTCConfiguration)

	// Need to convert each IceServer struct individually.
	total := len(config.IceServers)
	if total > 0 {
		cServers := make([]C.CGOIceServer, total)
		for i, server := range config.IceServers {
			cServers[i] = server.CGO()
		}
		c.iceServers = &cServers[0]
	}
	c.numIceServers = C.int(total)

	// c.iceServers = (*C.CGOIceServer)(unsafe.Pointer(&config.IceServers))
	c.iceTransportPolicy = C.int(config.IceTransportPolicy)
	c.bundlePolicy = C.int(config.BundlePolicy)
	// [ED] c.RtcpMuxPolicy = C.int(config.RtcpMuxPolicy)
	c.peerIdentity = C.CString(config.PeerIdentity)
	// [ED] c.Certificates = config.Certificates
	// [ED] c.IceCandidatePoolSize = C.int(config.IceCandidatePoolSize)
	config.cgoConfig = c
	return *c
}

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

// [ED]
// var _cgoRtcpMuxPolicyNegotiate = int(C.CGORtcpMuxPolicyNegotiate)
// var _cgoRtcpMuxPolicyRequire = int(C.CGORtcpMuxPolicyRequire)
