package webrtc

// #include "cpeerconnection.h"
// #include "ctestenums.h"
import "C"
import (
	"fmt"
	"strings"
	"errors"
	// "encoding/json"
)

// Working draft spec: http://www.w3.org/TR/webrtc/#idl-def-RTCConfiguration
// There are quite a few differences in the latest Editor's draft, but
// for now they are omitted from this Go interface, or commented out with
// an [ED] above.)
// See https://w3c.github.io/webrtc-pc/#idl-def-RTCConfiguration

type (
	RTCBundlePolicy int
	RTCIceTransportPolicy int
	RTCRtcpMuxPolicy int
	RTCIceCredentialType int
	RTCSignalingState int
)

type RTCConfiguration struct {
	// TODO: Implement, and provide as argument to CreatePeerConnection
	IceServers           []RTCIceServer
	IceTransportPolicy   RTCIceTransportPolicy
	BundlePolicy         RTCBundlePolicy
	// [ED] RtcpMuxPolicy        RTCRtcpMuxPolicy
	PeerIdentity         string   // Target peer identity

	// This would allow key continuity.
	// [ED] Certificates         []string
	// [ED] IceCandidatePoolSize int

	cgoConfig *C.CGORTCConfiguration  // Native code internals
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
	Urls           []string  // The only "required" element.
	Username       string
	Credential     string
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
		return nil, errors.New("IceServer: received more strings than expected.")
	}
	urls := strings.Split(params[0], ",")
	username := ""
	credential:= ""
	if 0 == len(urls) {
		return nil, errors.New("IceServer: requires at least one Url")
	}
	for i, url := range urls{
		url = strings.TrimSpace(url)
		// TODO: Better url validation.
		if !strings.HasPrefix(url, "stun:") && !strings.HasPrefix(url, "turn:") {
			return nil, errors.New(
					fmt.Sprintf("IceServer: received malformed url: <%s>", url))
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
		Urls: urls,
		Username: username,
		Credential: credential,
	}, nil
}

// Create a new RTCConfiguration with default values according to spec.
// Accepts any number of |RTCIceServer|s.
func NewRTCConfiguration(options... RTCConfigurationOption) *RTCConfiguration {
	c := new(RTCConfiguration)
	c.IceTransportPolicy = IceTransportPolicyAll
	c.BundlePolicy = BundlePolicyBalanced
	for _, op := range options {
		err := op(c)
		if nil != err {
			fmt.Println(err)
		}
	}
	// [ED] c.RtcpMuxPolicy = RtcpMuxPolicyRequire
	// [ED] c.Certificates = make([]string, 0)

	// fmt.Println(c)
	// b, _ := json.Marshal(c)
	// fmt.Printf("%q\n", b)
	// var c2 RTCConfiguration
	// _ = json.Unmarshal(b, &c2)
	// fmt.Println(c2)
	// fmt.Println("RTCConfiguration: ", c)
	return c
}

// Used in RTCConfiguration's variadic functional constructor
type RTCConfigurationOption func(c *RTCConfiguration) error

func OptionIceServer(params ...string) RTCConfigurationOption {
	return func(config *RTCConfiguration) error {
	// return RTCConfigurationOption {
  	// for _, server := range servers {
		return config.AddIceServer(params...)
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

func (config *RTCConfiguration) CGO() C.CGORTCConfiguration {
	c := new(C.CGORTCConfiguration)
	// c.IceServers = (C.CGOArray)(unsafe.Pointer(&config.IceServers[0]))
	c.IceTransportPolicy = C.int(config.IceTransportPolicy)
	// c.BundlePolicy = C.CString(config.BundlePolicy)
	c.BundlePolicy = C.int(config.BundlePolicy)
	// c.RtcpMuxPolicy = C.int(config.RtcpMuxPolicy)
	c.PeerIdentity = C.CString(config.PeerIdentity)
	// c.Certificates = config.Certificates
	// c.IceCandidatePoolSize = C.int(config.IceCandidatePoolSize)
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
