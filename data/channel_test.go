package data

import (
	// "fmt"
	"testing"
	"time"
	"reflect"
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
	c = NewChannel(cgoFakeDataChannel())
	if nil == c {
		t.Fatal("Could not create NewChannel")
	}
}

// TODO: There's not a good way to create a DataChannel without first having
// an available PeerConnection object with a valid session, but that's part of
// the outer package, making these tests pretty useless. To fix.

func TestChannelLabel(t *testing.T) {
	if "fake" != c.Label() {
		t.Error()
	}
}

func TestChannelOrdered(t *testing.T) {
	if false != c.Ordered() {
		t.Error()
	}
}

func TestChannelReadyState(t *testing.T) {
	if DataStateClosed != c.ReadyState() {
		t.Error()
	}
}

func TestOnMessageCallback(t *testing.T) {
	success := make(chan []byte, 1)
	c.OnMessage = func(msg []byte) {
		success <- msg
	}
	bytes := []byte("somenumberofbytesinhere")
	size := len(bytes)
	cgoFakeMessage(c, bytes, size)
	select {
	case data := <-success:
		if !reflect.DeepEqual(data, bytes) {
			t.Error("Unexpected bytes: ", data)
		}
	case <-time.After(time.Second * 1):
		t.Fatal("Timed out.")
	}
}


