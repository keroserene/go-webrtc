/*
Package webrtc/data contains the go wrapper for the Peer-to-Peer Data API
portion of WebRTC spec.

See: https://w3c.github.io/webrtc-pc/#idl-def-RTCDataChannel
*/
package data

/*
#cgo CXXFLAGS: -std=c++11
#cgo linux,amd64 pkg-config: webrtc-data-linux-amd64.pc
#cgo darwin,amd64 pkg-config: webrtc-data-darwin-amd64.pc
#include <stdlib.h>  // Needed for C.free
#include "cdatachannel.h"
*/
import "C"
import (
	"unsafe"
)

type DataState int

const (
	DataStateConnecting DataState = iota
	DataStateOpen
	DataStateClosing
	DataStateClosed
)

var DataStateString = []string{"Connecting", "Open", "Closing", "Closed"}

/* DataChannel

OnError - is not implemented because the underlying Send
always returns true as specified for SCTP, there is no reasonable
exposure of other specific errors from the native code, and OnClose
already covers the bases.
*/
type Channel struct {
	BufferedAmountLowThreshold int
	BinaryType                 string

	// Event Handlers
	OnOpen              func()
	OnClose             func()
	OnMessage           func([]byte) // byte slice.
	OnBufferedAmountLow func()

	cgoChannel C.CGO_Channel // Internal DataChannel functionality.
}

// Create a Go Channel struct, and prepare internal CGO references / observers.
// Expects cDC to be a pointer to a CGO_Channel object, which ultimately points
// to a DataChannelInterface*.
// The most reasonable place for this to be created is from PeerConnection,
// which is not available in the subpackage.
func NewChannel(cDC unsafe.Pointer) *Channel {
	if nil == cDC {
		return nil
	}
	dc := new(Channel)
	dc.cgoChannel = (C.CGO_Channel)(cDC)
	dc.BinaryType = "blob"
	// "Observer" is required for attaching callbacks correctly.
	C.CGO_Channel_RegisterObserver(dc.cgoChannel, unsafe.Pointer(dc))
	return dc
}

func (c *Channel) Send(data []byte) {
	if nil == data {
		return
	}
	C.CGO_Channel_Send(c.cgoChannel, unsafe.Pointer(&data[0]), C.int(len(data)))
}

func (c *Channel) Close() error {
	C.CGO_Channel_Close(c.cgoChannel)
	return nil
}

func (c *Channel) Label() string {
	s := C.CGO_Channel_Label(c.cgoChannel)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

func (c *Channel) Ordered() bool {
	return bool(C.CGO_Channel_Ordered(c.cgoChannel))
}

func (c *Channel) Protocol() string {
	p := C.CGO_Channel_Protocol(c.cgoChannel)
	defer C.free(unsafe.Pointer(p))
	return C.GoString(p)
}

func (c *Channel) MaxPacketLifeTime() uint {
	return uint(C.CGO_Channel_MaxRetransmitTime(c.cgoChannel))
}

func (c *Channel) MaxRetransmits() uint {
	return uint(C.CGO_Channel_MaxRetransmits(c.cgoChannel))
}

func (c *Channel) Negotiated() bool {
	return bool(C.CGO_Channel_Negotiated(c.cgoChannel))
}

func (c *Channel) ID() int {
	return int(C.CGO_Channel_ID(c.cgoChannel))
}

func (c *Channel) ReadyState() DataState {
	return (DataState)(C.CGO_Channel_ReadyState(c.cgoChannel))
}

func (c *Channel) BufferedAmount() int {
	return int(C.CGO_Channel_BufferedAmount(c.cgoChannel))
}

// TODO: Variadic options constructor, probably makes more sense for
// CreateDataChannel in parent package PeerConnection.
type Init struct {
	Ordered           bool
	MaxPacketLifeTime uint
	MaxRetransmits    uint
	Protocol          string
	Negotiated        bool
	ID                uint
}

//
// === cgo hooks for user-provided Go callbacks, and enums ===
//

//export cgoChannelOnMessage
func cgoChannelOnMessage(goChannel unsafe.Pointer, cBytes unsafe.Pointer, size int) {
	bytes := C.GoBytes(cBytes, C.int(size))
	dc := (*Channel)(goChannel)
	if nil != dc.OnMessage {
		dc.OnMessage(bytes)
	}
}

//export cgoChannelOnStateChange
func cgoChannelOnStateChange(goChannel unsafe.Pointer) {
	dc := (*Channel)(goChannel)
	switch dc.ReadyState() {
	// Picks between different Go callbacks...
	case DataStateConnecting:
	case DataStateClosing:
		// golang switches don't fallthrough
	case DataStateOpen:
		if nil != dc.OnOpen {
			dc.OnOpen()
		}
	case DataStateClosed:
		if nil != dc.OnClose {
			dc.OnClose()
		}
	default:
		panic("fired an un-implemented data.Channel StateChange.")
	}
}

//export cgoChannelOnBufferedAmountChange
func cgoChannelOnBufferedAmountChange(goChannel unsafe.Pointer, amount int) {
	dc := (*Channel)(goChannel)
	if nil != dc.OnBufferedAmountLow {
		if amount <= dc.BufferedAmountLowThreshold {
			dc.OnBufferedAmountLow()
		}
	}
}

var _cgoDataStateConnecting = int(C.CGO_DataStateConnecting)
var _cgoDataStateOpen = int(C.CGO_DataStateOpen)
var _cgoDataStateClosing = int(C.CGO_DataStateClosing)
var _cgoDataStateClosed = int(C.CGO_DataStateClosed)

// Testing helpers

func cgoFakeDataChannel() unsafe.Pointer {
	return unsafe.Pointer(C.CGO_getFakeDataChannel())
}

func cgoFakeMessage(c *Channel, b []byte, size int) {
	C.CGO_fakeMessage((C.CGO_Channel)(c.cgoChannel),
		unsafe.Pointer(&b[0]), C.int(size))
}

func cgoFakeStateChange(c *Channel, s DataState) {
	C.CGO_fakeStateChange((C.CGO_Channel)(c.cgoChannel), (C.int)(s))
}

func cgoFakeBufferAmount(c *Channel, amount int) {
	C.CGO_fakeBufferAmount((C.CGO_Channel)(c.cgoChannel), (C.int)(amount))
}
