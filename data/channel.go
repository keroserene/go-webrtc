/*
Package webrtc/data contains the go wrapper for the Peer-to-Peer Data API
portion of WebRTC spec.
*/
package data

/*
#cgo CXXFLAGS: -std=c++0x
#cgo CXXFLAGS: -I../include
#cgo LDFLAGS: -L../lib
#cgo linux,amd64 pkg-config: webrtc-data-linux-amd64.pc
#include "cdatachannel.h"
*/
import "C"
import (
	"unsafe"
	"fmt"
)


type DataState int

const (
	DataStateConnecting DataState = iota
	DataStateOpen
	DataStateClosing
	DataStateClosed
)

var DataStateString = []string{"Connecting", "Open", "Closing", "Closed"}

// data.Channel
type Channel struct {
	BufferedAmountLowThreshold int
	// TODO: Close() and Send()
	// TODO: OnOpen, OnBufferedAmountLow, OnError, OnClose, OnMessage,
	BinaryType string

	// Event Handlers
	OnOpen func()
	// OnError func()
	OnClose func()
	OnMessage func([]byte)  // byte slice.

	// TODO: Think about visibility and the implications of having
	// multiple packages like this...
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
	fmt.Println("Go channel at: ", unsafe.Pointer(dc))
	dc.cgoChannel = (C.CGO_Channel)(cDC)
	// Observer is required for attaching callbacks correctly.
	C.CGO_Channel_RegisterObserver(dc.cgoChannel, unsafe.Pointer(dc))
	return dc
}

func (c *Channel) Close() error {
	C.CGO_Channel_Close(c.cgoChannel)
	return nil
}

func (c *Channel) Label() string {
	s := C.CGO_Channel_Label(c.cgoChannel)
	return C.GoString(s)
}

func (c *Channel) Ordered() bool {
	return bool(C.CGO_Channel_Ordered(c.cgoChannel))
}

func (c *Channel) Protocol() string {
	return C.GoString(C.CGO_Channel_Protocol(c.cgoChannel))
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


type Init struct {
	// TODO: defaults
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
	fmt.Println("fired data.Channel.OnMessage: ", goChannel, bytes, size)
	dc := (*Channel)(goChannel)
	if nil != dc.OnMessage {
		dc.OnMessage(bytes)
	}
}

//export cgoChannelOnStateChange
func cgoChannelOnStateChange(c unsafe.Pointer) {
	dc := (*Channel)(c)
	// This event handler picks between different Go callbacks, depending
	// on the state.
	// TODO: look at state.
	switch dc.ReadyState() {
		case DataStateOpen:
			fmt.Println("fired data.Channel.OnOpen", c)
    	if nil != dc.OnOpen {
    		dc.OnOpen()
    	}
		case DataStateClosed:
			fmt.Println("fired data.Channel.OnClose", c)
    	if nil != dc.OnClose {
    		dc.OnClose()
    	}
		default:
			fmt.Println("fired an un-implemented data.Channel StateChange.", c)
	}
}

var _cgoDataStateConnecting = int(C.CGO_DataStateConnecting)
var _cgoDataStateOpen = int(C.CGO_DataStateOpen)
var _cgoDataStateClosing = int(C.CGO_DataStateClosing)
var _cgoDataStateClosed = int(C.CGO_DataStateClosed)

// Testing helpers 

func cgoFakeDataChannel() unsafe.Pointer {
	return unsafe.Pointer(C.CGO_getFakeDataChannel());
}

func cgoFakeMessage(c *Channel, b []byte, size int) {
	C.CGO_fakeMessage((C.CGO_Channel)(c.cgoChannel),
		unsafe.Pointer(&b[0]), C.int(size));
}

func cgoFakeStateChange(c *Channel, s DataState) {
	C.CGO_fakeStateChange((C.CGO_Channel)(c.cgoChannel), (C.int)(s))
}
