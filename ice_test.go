package webrtc

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIceCandidate(t *testing.T) {
	Convey("IceCandidate", t, func() {
		SetLoggingVerbosity(0)

		Convey("Serialize and Deserialize", func() {
			ice := IceCandidate{
				"fake",
				"not real",
				1337,
			}
			expected := `{"candidate":"fake","sdpMid":"not real","sdpMLineIndex":1337}`
			So(ice.Serialize(), ShouldEqual, expected)

			r := DeserializeIceCandidate(`{"candidate":"still fake","sdpMid":"illusory","sdpMLineIndex":1337}`)
			So(r, ShouldNotBeNil)
			So(r.Candidate, ShouldEqual, "still fake")
			So(r.SdpMid, ShouldEqual, "illusory")
			So(r.SdpMLineIndex, ShouldEqual, 1337)

			r = DeserializeIceCandidate(`not valid {{{`)
			So(r, ShouldBeNil)

			// Missing fields should result in errors.
			r = DeserializeIceCandidate(`{"sdpMid":"foo","sdpMLineIndex":1234}`)
			So(r, ShouldBeNil)
			r = DeserializeIceCandidate(`{"candidate":"something","sdpMLineIndex":1234}`)
			So(r, ShouldBeNil)
			r = DeserializeIceCandidate(`{"candidate":"something","sdpMid":"bar"}`)
			So(r, ShouldBeNil)

			Convey("Roundtrip", func() {
				ice := IceCandidate{
					"totally fake",
					"fabricated",
					1337,
				}
				r := DeserializeIceCandidate(ice.Serialize())
				So(r, ShouldNotBeNil)
				So(r.Candidate, ShouldEqual, ice.Candidate)
				So(r.SdpMid, ShouldEqual, ice.SdpMid)
				So(r.SdpMLineIndex, ShouldEqual, ice.SdpMLineIndex)
			})
		})
	})
}
