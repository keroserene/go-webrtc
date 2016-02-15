/*
Package webrtc/data contains the go wrapper for the Peer-to-Peer Data API
portion of WebRTC spec.

See: https://w3c.github.io/webrtc-pc/#idl-def-RTCDataChannel
*/
package webrtc

/*
#cgo CXXFLAGS: -std=c++0x
#cgo LDFLAGS: -L${SRCDIR}/lib
#cgo linux,amd64 pkg-config: webrtc-linux-amd64.pc
#cgo darwin,amd64 pkg-config: webrtc-darwin-amd64.pc
#include <stdlib.h>  // Needed for C.free
#include "datachannel.h"
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
type DataChannel struct {
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
// The most reasonable place for this to be created is from PeerConnection,
// which is not available in the subpackage.
func NewDataChannel(o unsafe.Pointer) *DataChannel {
	if o == nil {
		return nil
	}
	c := new(DataChannel)
	c.BinaryType = "blob"
	cgoChannel := C.CGO_Channel_RegisterObserver(o, unsafe.Pointer(c))
	c.cgoChannel = (C.CGO_Channel)(cgoChannel)
	return c
}

func (c *DataChannel) Send(data []byte) {
	if nil == data {
		return
	}
	C.CGO_Channel_Send(c.cgoChannel, unsafe.Pointer(&data[0]), C.int(len(data)))
}

func (c *DataChannel) Close() error {
	C.CGO_Channel_Close(c.cgoChannel)
	return nil
}

func (c *DataChannel) Label() string {
	s := C.CGO_Channel_Label(c.cgoChannel)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

func (c *DataChannel) Ordered() bool {
	return bool(C.CGO_Channel_Ordered(c.cgoChannel))
}

func (c *DataChannel) Protocol() string {
	p := C.CGO_Channel_Protocol(c.cgoChannel)
	defer C.free(unsafe.Pointer(p))
	return C.GoString(p)
}

func (c *DataChannel) MaxPacketLifeTime() uint {
	return uint(C.CGO_Channel_MaxRetransmitTime(c.cgoChannel))
}

func (c *DataChannel) MaxRetransmits() uint {
	return uint(C.CGO_Channel_MaxRetransmits(c.cgoChannel))
}

func (c *DataChannel) Negotiated() bool {
	return bool(C.CGO_Channel_Negotiated(c.cgoChannel))
}

func (c *DataChannel) ID() int {
	return int(C.CGO_Channel_ID(c.cgoChannel))
}

func (c *DataChannel) ReadyState() DataState {
	return (DataState)(C.CGO_Channel_ReadyState(c.cgoChannel))
}

func (c *DataChannel) BufferedAmount() int {
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
	dc := (*DataChannel)(goChannel)
	if nil != dc.OnMessage {
		dc.OnMessage(bytes)
	}
}

//export cgoChannelOnStateChange
func cgoChannelOnStateChange(goChannel unsafe.Pointer) {
	dc := (*DataChannel)(goChannel)
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
	dc := (*DataChannel)(goChannel)
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

func cgoFakeMessage(c *DataChannel, b []byte, size int) {
	C.CGO_fakeMessage((C.CGO_Channel)(c.cgoChannel),
		unsafe.Pointer(&b[0]), C.int(size))
}

func cgoFakeStateChange(c *DataChannel, s DataState) {
	C.CGO_fakeStateChange((C.CGO_Channel)(c.cgoChannel), (C.int)(s))
}

func cgoFakeBufferAmount(c *DataChannel, amount int) {
	C.CGO_fakeBufferAmount((C.CGO_Channel)(c.cgoChannel), (C.int)(amount))
}
