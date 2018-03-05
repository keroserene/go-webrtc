package webrtc

// #include <stdlib.h>
// #include "peerconnection.h"
// #include "ctestenums.h"
import "C"
import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

// Working draft spec: http://www.w3.org/TR/webrtc/#idl-def-Configuration
// There are quite a few differences in the latest Editor's draft, but
// for now they are omitted from this Go interface, or commented out with
// an [ED] above.)
// See https://w3c.github.io/webrtc-pc/#idl-def-Configuration

type (
	BundlePolicy       int
	IceTransportPolicy int
	RtcpMuxPolicy      int
	IceCredentialType  int
	SignalingState     int
)

type Configuration struct {
	IceServers []IceServer
	IceTransportPolicy
	BundlePolicy
	// [ED] RtcpMuxPolicy        RtcpMuxPolicy
	PeerIdentity string // Target peer identity

	// This would allow key continuity.
	// [ED] Certificates         []string
	// [ED] IceCandidatePoolSize int
}

// These "Enum" consts must match order in: peerconnectioninterface.h
// There doesn't seem to be a way to have a named container for enums
// in go, and the idiomatic way seems to be just prefixes.
const (
	BundlePolicyBalanced BundlePolicy = iota
	BundlePolicyMaxBundle
	BundlePolicyMaxCompat
)

func (p BundlePolicy) String() string {
	return EnumToStringSafe(int(p), []string{
		"Balanced",
		"MaxBundle",
		"MaxCompat",
	})
}

const (
	IceTransportPolicyNone IceTransportPolicy = iota
	IceTransportPolicyRelay
	// TODO: Look into why nohost is not exposed in w3c spec, but is available
	// in native code? If it does need to be exposed, capitalize the i.
	// (It still needs to exist, to ensure the enum values match up.
	iceTransportPolicyNoHost
	IceTransportPolicyAll
)

func (p IceTransportPolicy) String() string {
	return EnumToStringSafe(int(p), []string{
		"None",
		"Relay",
		"NoHost",
		"All",
	})
}

const (
	SignalingStateStable SignalingState = iota
	SignalingStateHaveLocalOffer
	SignalingStateHaveLocalPrAnswer
	SignalingStateHaveRemoteOffer
	SignalingStateHaveRemotePrAnswer
	SignalingStateClosed
)

func (s SignalingState) String() string {
	return EnumToStringSafe(int(s), []string{
		"Stable",
		"HaveLocalOffer",
		"HaveLocalPrAnswer",
		"HaveRemoteOffer",
		"HaveRemotePrAnswer",
		"Closed",
	})
}

// TODO: [ED]
/* const (
	RtcpMuxPolicyNegotiate RtcpMuxPolicy = iota
	RtcpMuxPolicyRequire
) */

// TODO: [ED]
/* const (
	IceCredentialTypePassword IceCredentialType = iota
	IceCredentialTypeToken
) */

type IceServer struct {
	Urls       []string // The only "required" element.
	Username   string
	Credential string
	// [ED] CredentialType IceCredentialType
}

// Create a new IceServer object.
// Expects anywhere from one to three strings, in this order:
// - comma-separated list of urls.
// - username
// - credential
// TODO: For the ED version, may need to support CredentialType.
func NewIceServer(params ...string) (*IceServer, error) {
	if len(params) < 1 {
		return nil, errors.New("IceServer: missing first comma-separated Urls string.")
	}
	if len(params) > 3 {
		WARN.Printf("IceServer: got %d strings, expect <= 3. Ignoring extras.\n",
			len(params))
	}
	if "" == params[0] {
		return nil, errors.New("IceServer: requires at least one Url")
	}
	urls := strings.Split(params[0], ",")
	username := ""
	credential := ""
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
	return &IceServer{
		Urls:       urls,
		Username:   username,
		Credential: credential,
	}, nil
}

// Create a new Configuration with default values according to spec.
// Accepts any number of |IceServer|s.
// Returns nil if there's an error.
func NewConfiguration(options ...ConfigurationOption) *Configuration {
	c := new(Configuration)
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
	INFO.Println("Created Configuration at ", c)
	// TODO: Determine whether the below is true.
	// if 0 == len(c.IceServers) {
	// ERROR.Println("Need at least one ICE server.")
	// return nil
	// }
	return c
}

// Used in Configuration's variadic functional constructor
type ConfigurationOption func(c *Configuration) error

func OptionIceServer(params ...string) ConfigurationOption {
	return func(config *Configuration) error {
		return config.AddIceServer(params...)
	}
}

func OptionIceTransportPolicy(policy IceTransportPolicy) ConfigurationOption {
	return func(config *Configuration) error {
		INFO.Println("OptionIceTransportPolicy: ", policy)
		config.IceTransportPolicy = policy
		return nil
	}
}

func OptionBundlePolicy(policy BundlePolicy) ConfigurationOption {
	return func(config *Configuration) error {
		config.BundlePolicy = policy
		return nil
	}
}

func (config *Configuration) AddIceServer(params ...string) error {
	server, err := NewIceServer(params...)
	if nil != err {
		return err
	}
	config.IceServers = append(config.IceServers, *server)
	return nil
}

// Helpers which prepare Go-side of cast to eventual C++ Configuration struct.
func (server *IceServer) _CGO() C.CGO_IceServer {
	cServer := new(C.CGO_IceServer)

	// TODO: Make this conversion nicer.
	total := len(server.Urls)
	if total > 0 {
		sizeof := unsafe.Sizeof(uintptr(0)) // FIXME(arlolra): sizeof *void
		cUrls := unsafe.Pointer(C.malloc(C.size_t(sizeof * uintptr(total))))
		for i, url := range server.Urls {
			*(**C.char)(unsafe.Pointer(uintptr(cUrls) + sizeof*uintptr(i))) = C.CString(url)
		}
		cServer.urls = (**C.char)(cUrls)
	}

	cServer.numUrls = C.int(total)
	cServer.username = C.CString(server.Username)
	cServer.credential = C.CString(server.Credential)
	return *cServer
}

const maxUrls = 1 << 24

func freeIceServer(cServer C.CGO_IceServer) {
	total := int(cServer.numUrls)
	if total > maxUrls {
		panic("Too many urls. Something went wrong.")
	}
	cUrls := (*[maxUrls](*C.char))(unsafe.Pointer(cServer.urls))
	for i := 0; i < total; i++ {
		C.free(unsafe.Pointer(cUrls[i]))
	}
	C.free(unsafe.Pointer(cServer.username))
	C.free(unsafe.Pointer(cServer.credential))
	C.free(unsafe.Pointer(cServer.urls))
}

// The C side of things will still need to allocate memory, due to the slices.
// Assumes Configuration is valid.
func (config *Configuration) _CGO() *C.CGO_Configuration {
	INFO.Println("Converting Config: ", config)
	size := C.size_t(unsafe.Sizeof(C.CGO_Configuration{}))
	c := (*C.CGO_Configuration)(C.malloc(size))

	// Need to convert each IceServer struct individually.
	total := len(config.IceServers)
	if total > 0 {
		sizeof := unsafe.Sizeof(C.CGO_IceServer{})
		cServers := unsafe.Pointer(C.malloc(C.size_t(sizeof * uintptr(total))))
		for i, server := range config.IceServers {
			*(*C.CGO_IceServer)(unsafe.Pointer(uintptr(cServers) + sizeof*uintptr(i))) = server._CGO()
		}
		c.iceServers = (*C.CGO_IceServer)(cServers)
	}
	c.numIceServers = C.int(total)

	// c.iceServers = (*C.CGO_IceServer)(unsafe.Pointer(&config.IceServers))
	c.iceTransportPolicy = C.int(config.IceTransportPolicy)
	c.bundlePolicy = C.int(config.BundlePolicy)
	// [ED] c.RtcpMuxPolicy = C.int(config.RtcpMuxPolicy)
	c.peerIdentity = C.CString(config.PeerIdentity)
	// [ED] c.Certificates = config.Certificates
	// [ED] c.IceCandidatePoolSize = C.int(config.IceCandidatePoolSize)
	return c
}

const maxIceServers = 1 << 24

func freeConfig(cConfig *C.CGO_Configuration) {
	total := int(cConfig.numIceServers)
	if total > maxIceServers {
		panic("Too many ice servers. Something went wrong.")
	} else if total > 0 {
		cServers := (*[maxIceServers]C.CGO_IceServer)(unsafe.Pointer(cConfig.iceServers))
		for i := 0; i < total; i++ {
			freeIceServer(cServers[i])
		}
		C.free(unsafe.Pointer(cConfig.iceServers))
	}
	C.free(unsafe.Pointer(cConfig.peerIdentity))
	C.free(unsafe.Pointer(cConfig))
}

/*
const {
  stable SignallingState = iota
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

var _cgoIceTransportPolicyNone = int(C.CGO_IceTransportPolicyNone)
var _cgoIceTransportPolicyRelay = int(C.CGO_IceTransportPolicyRelay)
var _cgoIceTransportPolicyNoHost = int(C.CGO_IceTransportPolicyNoHost)
var _cgoIceTransportPolicyAll = int(C.CGO_IceTransportPolicyAll)

var _cgoBundlePolicyBalanced = int(C.CGO_BundlePolicyBalanced)
var _cgoBundlePolicyMaxCompat = int(C.CGO_BundlePolicyMaxCompat)
var _cgoBundlePolicyMaxBundle = int(C.CGO_BundlePolicyMaxBundle)

// [ED]
// var _cgoRtcpMuxPolicyNegotiate = int(C.CGO_RtcpMuxPolicyNegotiate)
// var _cgoRtcpMuxPolicyRequire = int(C.CGO_RtcpMuxPolicyRequire)

var _cgoSignalingStateStable = int(C.CGO_SignalingStateStable)
var _cgoSignalingStateHaveLocalOffer = int(C.CGO_SignalingStateHaveLocalOffer)
var _cgoSignalingStateHaveLocalPrAnswer = int(C.CGO_SignalingStateHaveLocalPrAnswer)
var _cgoSignalingStateHaveRemoteOffer = int(C.CGO_SignalingStateHaveRemoteOffer)
var _cgoSignalingStateHaveRemotePrAnswer = int(C.CGO_SignalingStateHaveRemotePrAnswer)
var _cgoSignalingStateClosed = int(C.CGO_SignalingStateClosed)
