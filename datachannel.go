package webrtc

// #include "cpeerconnection.h"
import "C"

type DataChannel struct {
	Label                      string
	Ordered                    bool
	MaxPacketLifeTime          uint
	MaxRetransmits             uint
	Protocol                   string
	Negotiated                 bool
	ID                         uint
	ReadyState                 string // RTCDataChannelState
	BufferedAmount             int
	BufferedAmountLowThreshold int
	// TODO: Close() and Send()
	// TODO: OnOpen, OnBufferedAmountLow, OnError, OnClose, OnMessage,
	BinaryType string

	cgoDataChannel C.CGODataChannel // Internal PeerConnection functionality.
}

type DataChannelInit struct {
	// TODO: defaults
	Ordered           bool
	MaxPacketLifeTime uint
	MaxRetransmits    uint
	Protocol          string
	Negotiated        bool
	ID                uint
}

// func (dc *DataChannel) CGO() C.CGODataChannel {
  // return (C.CGODataChannel)
// }

func newDataChannel() *DataChannel {
  return new(DataChannel)
}
