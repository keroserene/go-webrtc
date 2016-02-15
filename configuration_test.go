package webrtc

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConfiguration(t *testing.T) {
	SetLoggingVerbosity(0)

	Convey("Go enum values should correspond to native C++ values.", t, func() {
		// Ensure the Go "enums" generated in the idiomatic iota const way actually
		// match up with actual int values of the underlying native WebRTC Enums.

		Convey("Enum: BundlePolicy", func() {
			So(BundlePolicyBalanced, ShouldEqual, _cgoBundlePolicyBalanced)
			So(BundlePolicyMaxCompat, ShouldEqual, _cgoBundlePolicyMaxCompat)
			So(BundlePolicyMaxBundle, ShouldEqual, _cgoBundlePolicyMaxBundle)
		})

		Convey("Enum: IceTransportPolicy", func() {
			So(IceTransportPolicyNone, ShouldEqual, IceTransportPolicyNone)
			So(IceTransportPolicyRelay, ShouldEqual, IceTransportPolicyRelay)
			So(IceTransportPolicyAll, ShouldEqual, IceTransportPolicyAll)
		})

		Convey("Enum: SignalingState", func() {
			So(SignalingStateStable, ShouldEqual, _cgoSignalingStateStable)
			So(SignalingStateHaveLocalOffer, ShouldEqual, _cgoSignalingStateHaveLocalOffer)
			So(SignalingStateHaveLocalPrAnswer, ShouldEqual, _cgoSignalingStateHaveLocalPrAnswer)
			So(SignalingStateHaveRemoteOffer, ShouldEqual, _cgoSignalingStateHaveRemoteOffer)
			So(SignalingStateHaveRemotePrAnswer, ShouldEqual, _cgoSignalingStateHaveRemotePrAnswer)
			So(SignalingStateClosed, ShouldEqual, _cgoSignalingStateClosed)
		})

		// TODO: [ED]
		// SkipConvey("Enum: RtcpMuxPolicy", func() {
		// So(RtcpMuxPolicyNegotiate, ShouldEqual, _cgoRtcpMuxPolicyNegotiate)
		// So(RtcpMuxPolicyRequire, ShouldEqual, _cgoRtcpMuxPolicyRequire)
		// })

	}) // Enums

	Convey("New IceServer", t, func() {
		s, err := NewIceServer()
		So(err, ShouldNotBeNil) // 0 params
		So(s, ShouldBeNil)

		s, err = NewIceServer("")
		So(err, ShouldNotBeNil) // empty URL
		So(s, ShouldBeNil)

		s, err = NewIceServer("stun:12345, badurl")
		So(err, ShouldNotBeNil) // malformed URL
		So(s, ShouldBeNil)

		s, err = NewIceServer("stun:12345, stun:ok")
		So(err, ShouldBeNil)
		So(s, ShouldNotBeNil)

		s, err = NewIceServer("stun:a, turn:b")
		So(err, ShouldBeNil)
		So(s, ShouldNotBeNil)

		s, err = NewIceServer("stun:a, turn:b", "alice")
		So(err, ShouldBeNil)
		So(s, ShouldNotBeNil)

		s, err = NewIceServer("stun:a, turn:b", "alice", "secret")
		So(err, ShouldBeNil)
		So(s, ShouldNotBeNil)

		s, err = NewIceServer("stun:a, turn:b", "alice", "secret", "extra")
		So(err, ShouldBeNil) // NewIceServer shouldn't fail, only WARN on too many params
		So(s, ShouldNotBeNil)
	})

	Convey("New Configuration", t, func() {
		config := NewConfiguration()
		So(config, ShouldNotBeNil)

		config = NewConfiguration(OptionIceServer("stun:a"))
		So(len(config.IceServers), ShouldEqual, 1)

		config = NewConfiguration(
			OptionIceServer("stun:a"),
			OptionIceServer("stun:b, turn:c"))
		So(len(config.IceServers), ShouldEqual, 2)

		config = NewConfiguration(
			OptionIceServer("stun:d"),
			OptionIceTransportPolicy(IceTransportPolicyAll))
		So(config.IceTransportPolicy, ShouldEqual, IceTransportPolicyAll)

		config = NewConfiguration(
			OptionIceServer("stun:d"),
			OptionBundlePolicy(BundlePolicyMaxCompat))
		So(config.BundlePolicy, ShouldEqual, BundlePolicyMaxCompat)
	})
}
