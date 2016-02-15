package webrtc

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSessionDescription(t *testing.T) {
	Convey("SessionDescription", t, func() {
		Convey("Serialize and Deserialize", func() {
			expected := `{"type":"answer","sdp":"fake"}`
			r := DeserializeSessionDescription(expected)
			So(r, ShouldNotBeNil)
			So(r.Type, ShouldEqual, "answer")
			So(r.Sdp, ShouldEqual, "fake")
			s := r.Serialize()
			So(s, ShouldEqual, expected)

			r = DeserializeSessionDescription(`invalid json{{`)
			So(r, ShouldBeNil)
			r = DeserializeSessionDescription(`{"sdp":"fake"}`)
			So(r, ShouldBeNil)
			r = DeserializeSessionDescription(`{"type":"answer"}`)
			So(r, ShouldBeNil)

			Convey("Roundtrip", func() {
				sdp := SessionDescription{"pranswer", "not real"}
				r = DeserializeSessionDescription(sdp.Serialize())
				So(r.Type, ShouldEqual, sdp.Type)
				So(r.Sdp, ShouldEqual, sdp.Sdp)
			})
		})
	})
}
