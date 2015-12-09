/*
Package webrtc/data contains the go wrapper for the Peer-to-Peer Data API
portion of WebRTC spec.
*/
package data

/*
#cgo CXXFLAGS: -std=c++0x
#cgo CXXFLAGS: -I../include
#cgo linux,amd64 pkg-config: ../webrtc-linux-amd64.pc

#include "../cpeerconnection.h"
#include "cdatachannel.h"
// #cgo LDFLAGS: -l../include
*/
import "C"
import (
	"unsafe"
	// "fmt"
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
	MaxPacketLifeTime          uint
	MaxRetransmits             uint
	Protocol                   string
	Negotiated                 bool
	ID                         uint
	BufferedAmount             int
	BufferedAmountLowThreshold int
	// TODO: Close() and Send()
	// TODO: OnOpen, OnBufferedAmountLow, OnError, OnClose, OnMessage,
	BinaryType string

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

func (c *Channel) ReadyState() DataState {
	return (DataState)(C.CGO_Channel_ReadyState(c.cgoChannel))
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

// func (dc *DataChannel) CGO_() C.CGO_Channel {
  // return (C.CGO_Channel)
// }

// func NewChannel(cDC C.CGO_Channel) *Channel {
func NewChannel(cDC unsafe.Pointer) *Channel {
  dc := new(Channel)
	dc.cgoChannel = (C.CGO_Channel)(cDC)
	return dc
}

// func (channel *Channel) _CGO() C.CGO_Channel {
	// return (C.CGO_Channel)(unsafePointer(channel))	
// }

var _cgoDataStateConnecting = int(C.CGO_DataStateConnecting)
var _cgoDataStateOpen = int(C.CGO_DataStateOpen)
var _cgoDataStateClosing = int(C.CGO_DataStateClosing)
var _cgoDataStateClosed = int(C.CGO_DataStateClosed)
