package webrtc

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSessionDescription(t *testing.T) {
	Convey("SessionDescription", t, func() {
		r := NewSessionDescription("offer", "fake")
		So(r, ShouldNotBeNil)
		So(r.Type, ShouldEqual, "offer")
		So(r.Sdp, ShouldEqual, "fake")

		Convey("Serialize and Deserialize", func() {
			sdp := NewSessionDescription("answer", "fake")
			s := sdp.Serialize()
			So(s, ShouldEqual, `{"type":"answer","sdp":"fake"}`)

			r := DeserializeSessionDescription(`{"type":"answer","sdp":"fake"}`)
			So(r, ShouldNotBeNil)
			So(r.Type, ShouldEqual, "answer")
			So(r.Sdp, ShouldEqual, "fake")

			Convey("Roundtrip", func() {
				sdp = NewSessionDescription("pranswer", "not real")
				r = DeserializeSessionDescription(sdp.Serialize())
				So(r.Type, ShouldEqual, sdp.Type)
				So(r.Sdp, ShouldEqual, sdp.Sdp)
			})
		})
	})
}
