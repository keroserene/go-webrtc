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

#include "../cpeerconnection.h"
#include "cdatachannel.h"
// #cgo LDFLAGS: -l../include
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
	OnError func()
	OnClose func()
	OnMessage func([]byte)  // byte slice.

	// TODO: Think about visibility and the implications of having
	// multiple packages like this...
	cgoChannel C.CGO_Channel // Internal DataChannel functionality.
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

func NewChannel(cDC unsafe.Pointer) *Channel {
  dc := new(Channel)
	dc.cgoChannel = (C.CGO_Channel)(cDC)
	return dc
}

//
// === cgo hooks for user-provided Go callbacks, and enums ===
//

//export cgoChannelOnMessage
func cgoChannelOnMessage(c unsafe.Pointer, b []byte) {
	fmt.Println("fired data.Channel.OnMessage: ", c, b)
	dc := (*Channel)(c)
	if nil != dc.OnMessage {
		dc.OnMessage(b)
	}
}

//export cgoChannelOnStateChange
func cgoChannelOnStateChange(c unsafe.Pointer) {
	// This event handler picks between different Go callbacks, depending
	// on the state.
	fmt.Println("fired data.Channel.OnStateChange:", c)
	// dc := (*Channel)(c)
	// TODO: look at state.
	// if nil != dc.OnClosed {
		// pc.OnClosed
	// }
}


var _cgoDataStateConnecting = int(C.CGO_DataStateConnecting)
var _cgoDataStateOpen = int(C.CGO_DataStateOpen)
var _cgoDataStateClosing = int(C.CGO_DataStateClosing)
var _cgoDataStateClosed = int(C.CGO_DataStateClosed)
