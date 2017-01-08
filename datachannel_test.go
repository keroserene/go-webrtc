package webrtc

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDataStateEnums(t *testing.T) {

	Convey("Enum: DataState", t, func() {
		So(DataStateConnecting, ShouldEqual, _cgoDataStateConnecting)
		So(DataStateOpen, ShouldEqual, _cgoDataStateOpen)
		So(DataStateClosing, ShouldEqual, _cgoDataStateClosing)
		So(DataStateClosed, ShouldEqual, _cgoDataStateClosed)
	})

	Convey("DataChannel", t, func() {

		c := NewDataChannel(nil)
		So(c, ShouldBeNil)

		c = NewDataChannel(cgoFakeDataChannel())
		So(c, ShouldNotBeNil)
		So(c.Label(), ShouldEqual, "fake")
		So(c.Ordered(), ShouldBeFalse)
		So(c.Protocol(), ShouldEqual, "")
		So(c.MaxPacketLifeTime(), ShouldEqual, 0)
		So(c.MaxRetransmits(), ShouldEqual, 0)
		So(c.Negotiated(), ShouldBeFalse)
		So(c.ID(), ShouldEqual, 12345)
		So(c.BufferedAmount(), ShouldEqual, 1234)

		// There's not a good way to create a DataChannel without first having an
		// available PeerConnection object with a valid session, but that's part of
		// the outer package, so testing the attributes is less useful here.
		Convey("Callbacks should fire correctly", func() {

			Convey("OnMessage", func() {
				success := make(chan []byte, 1)
				c.OnMessage = func(msg []byte) {
					success <- msg
				}
				bytes := []byte("somenumberofbytesinhere")
				size := len(bytes)
				cgoFakeMessage(c, bytes, size)
				select {
				case data := <-success:
					So(c.OnMessage, ShouldNotBeNil)
					So(data, ShouldResemble, bytes)
				case <-time.After(time.Second * 1):
					t.Fatal("Timed out.")
				}
			})

			Convey("StateChangeCallbacks", func() {
				opened := make(chan int, 1)
				closed := make(chan int, 1)
				c.OnOpen = func() {
					opened <- 1
				}
				c.OnClose = func() {
					closed <- 1
				}

				cgoFakeStateChange(c, DataStateOpen)
				select {
				case <-opened:
					So(c.OnOpen, ShouldNotBeNil)
					So(c.ReadyState(), ShouldEqual, DataStateOpen)
				case <-time.After(time.Second * 1):
					t.Fatal("Timed out when waiting for Open.")
				}

				cgoFakeStateChange(c, DataStateClosed)
				select {
				case <-closed:
					So(c.OnClose, ShouldNotBeNil)
					So(c.ReadyState(), ShouldEqual, DataStateClosed)
				case <-time.After(time.Second * 1):
					t.Fatal("Timed out when waiting for Closed.")
				}

				// TODO: Unimplemented
				cgoFakeStateChange(c, DataStateConnecting)
				cgoFakeStateChange(c, DataStateClosing)

				So(func() {
					cgoFakeStateChange(c, 999)
				}, ShouldPanic)
			})

			Convey("OnBufferedAmountLow", func() {
				success := make(chan int, 1)
				c.BufferedAmountLowThreshold = 100
				c.OnBufferedAmountLow = func() {
					success <- 1
				}
				cgoFakeBufferAmount(c, 90)
				select {
				case <-success:
					So(c.OnBufferedAmountLow, ShouldNotBeNil)
				case <-time.After(time.Second * 1):
					t.Fatal("Timed out.")
				}
			})
		})

		Convey("Send", func() {
			messages := make(chan []byte, 1)
			data := []byte("some data to send")
			// Fake data channel routes send to its own onmessage.
			c.OnMessage = func(msg []byte) {
				messages <- msg
			}
			c.Send(data)
			select {
			case recv := <-messages:
				So(c.OnMessage, ShouldNotBeNil)
				So(recv, ShouldResemble, data)
			case <-time.After(time.Second * 1):
				t.Fatal("Timed out.")
			}
			c.Send(nil)
			select {
			case <-messages:
				t.Fatal("Unexpected message when sending nil.")
			case <-time.After(time.Second * 1):
			}
		})

		Convey("SendText", func() {
			messages := make(chan []byte, 1)
			text := "Hello, 世界"
			// Fake data channel routes send to its own onmessage.
			c.OnMessage = func(msg []byte) {
				messages <- msg
			}
			c.SendText(text)
			select {
			case recv := <-messages:
				So(c.OnMessage, ShouldNotBeNil)
				So(recv, ShouldResemble, []byte(text))
			case <-time.After(time.Second * 1):
				t.Fatal("Timed out.")
			}
			c.SendText("")
			select {
			case <-messages:
				t.Fatal("Unexpected message when sending nil.")
			case <-time.After(time.Second * 1):
			}
		})

		Convey("Close", func() {
			closed := make(chan int, 1)
			c.OnClose = func() {
				closed <- 1
			}
			c.Close()
			select {
			case <-closed:
				So(c.OnClose, ShouldNotBeNil)
				So(c.ReadyState(), ShouldEqual, DataStateClosed)
			case <-time.After(time.Second * 1):
				t.Fatal("Timed out during close..")
			}
		})
	})

}
