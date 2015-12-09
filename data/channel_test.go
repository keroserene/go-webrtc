package data

import (
	"fmt"
	"testing"
)

var c *Channel

func checkEnum(t *testing.T, desc string, enum int, expected int) {
	if enum != expected {
		t.Error("Mismatched Enum Value -", desc,
			"\nwas:", enum,
			"\nexpected:", expected)
	}
}

func TestDataStateEnums(t *testing.T) {
	checkEnum(t, "DataStateConnecting",
		int(DataStateConnecting), _cgoDataStateConnecting)
	checkEnum(t, "DataStateOpen",
		int(DataStateOpen), _cgoDataStateOpen)
	checkEnum(t, "DataStateClosing",
		int(DataStateClosing), _cgoDataStateClosing)
	checkEnum(t, "DataStateClosed",
		int(DataStateClosed), _cgoDataStateClosed)
}

func TestNewChannel(t *testing.T) {
	c = NewChannel(nil)
	if nil == c {
		t.Fatal("Could not create NewChannel")
	}
}

// TODO: There's not a good way to create a DataChannel without first having
// an available PeerConnection object with a valid session, but that's part of
// the outer package, making these tests pretty useless. To fix.

func TestChannelLabel(t *testing.T) {
	s := c.Label()
	fmt.Println(s)
}

func TestChannelOrdered(t *testing.T) {}

func TestChannelReadyState(t *testing.T) {
	// state := c.ReadyState()
	// fmt.Println(state)
}
