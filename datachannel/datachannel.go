package datachannel

// #include "../cpeerconnection.h"
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

	// TODO: Think about visibility and the implications of having
	// multiple packages like this...
	cgoDataChannel C.CGODataChannel // Internal PeerConnection functionality.
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

// func (dc *DataChannel) CGO() C.CGODataChannel {
  // return (C.CGODataChannel)
// }

func New() *DataChannel {
  return new(DataChannel)
}
