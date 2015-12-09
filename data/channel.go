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
	"fmt"
)

// data.Channel
type Channel struct {
	// Label                      string
	Ordered                    bool
	MaxPacketLifeTime          uint
	MaxRetransmits             uint
	Protocol                   string
	Negotiated                 bool
	ID                         uint
	ReadyState                 string   // RTCDataChannelState
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
	fmt.Println("label: ", s)
	return C.GoString(s)
	// return C.CGO_Channel_Label((C.CGO_Channel)(c.cgoChannel))
	// return c.cgoChannel.label()
	// return "not implemented yet"
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
