package webrtc

import (
	"github.com/keroserene/go-webrtc/data"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
	"unsafe"
)

func TestIceGatheringStateEnums(t *testing.T) {
	Convey(`Enum: IceGatheringState values should match
C++ webrtc::PeerConnectionInterface values`, t, func() {
		So(IceGatheringStateNew, ShouldEqual, _cgoIceGatheringStateNew)
		So(IceGatheringStateGathering, ShouldEqual, _cgoIceGatheringStateGathering)
		So(IceGatheringStateComplete, ShouldEqual, _cgoIceGatheringStateComplete)
	})
}

func TestPeerConnection(t *testing.T) {
	SetLoggingVerbosity(0)

	Convey("PeerConnection", t, func() {
		var offer *SessionDescription
		var answer *SessionDescription
		config := NewConfiguration(
			OptionIceServer("stun:stun.l.google.com:19302, stun:another"))
		So(config, ShouldNotBeNil)

		Convey("Basic functionality", func() {
			pc, err := NewPeerConnection(nil)
			So(pc, ShouldBeNil)
			So(err, ShouldNotBeNil)
			// A Configuration is required.
			pc, err = NewPeerConnection(config)
			So(pc, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(pc.ConnectionState(), ShouldEqual, PeerConnectionStateNew)

			Convey("Set and Get Configuration", func() {
				config := NewConfiguration(
					OptionIceServer("stun:something.else"),
					OptionIceTransportPolicy(IceTransportPolicyRelay))
				pc.SetConfiguration(*config)
				got := pc.GetConfiguration()
				So(got.IceTransportPolicy, ShouldEqual, IceTransportPolicyRelay)
			})

			Convey("Callbacks fire correctly", func() {

				Convey("OnSignalingState", func() {
					success := make(chan SignalingState, 1)
					pc.OnSignalingStateChange = func(s SignalingState) {
						success <- s
					}
					cgoOnSignalingStateChange(unsafe.Pointer(pc), SignalingStateStable)
					select {
					case state := <-success:
						So(state, ShouldEqual, SignalingStateStable)
					case <-time.After(time.Second * 1):
						t.Fatal("Timed out.")
					}
				})

				Convey("OnNegotiationNeeded", func() {
					success := make(chan int, 1)
					pc.OnNegotiationNeeded = func() {
						success <- 0
					}
					cgoOnNegotiationNeeded(unsafe.Pointer(pc))
					select {
					case <-success:
					case <-time.After(time.Second * 1):
						t.Fatal("Timed out.")
					}
				})

				Convey("OnConnectionStateChange", func() {
					success := make(chan PeerConnectionState, 1)
					pc.OnConnectionStateChange = func(state PeerConnectionState) {
						success <- state
					}
					cgoOnConnectionStateChange(unsafe.Pointer(pc),
						PeerConnectionStateDisconnected)
					select {
					case r := <-success:
						So(r, ShouldEqual, PeerConnectionStateDisconnected)
					case <-time.After(time.Second * 1):
						t.Fatal("Timed out.")
					}
				})

				// TODO: Find better way to trigger a fake ICE candidate.
				SkipConvey("OnIceCandidate", func() {
					success := make(chan IceCandidate, 1)
					pc.OnIceCandidate = func(ic IceCandidate) {
						success <- ic
					}
					// ice := DeserializeIceCandidate("fake")
					// cgoOnIceCandidate(unsafe.Pointer(pc), nil)
					select {
					case <-success:
					case <-time.After(time.Second * 1):
						t.Fatal("Timed out.")
					}
				})

				Convey("OnDataChannel", func() {
					success := make(chan *data.Channel, 1)
					pc.OnDataChannel = func(dc *data.Channel) {
						success <- dc
					}
					cgoOnDataChannel(unsafe.Pointer(pc), nil)
					select {
					case <-success:
					case <-time.After(time.Second * 1):
						t.Fatal("Timed out.")
					}
				})
			}) // Callbacks
		}) // Basic Functionality

		Convey("Create PeerConnections for Alice and Bob", func() {
			alice, err := NewPeerConnection(config)
			So(alice, ShouldNotBeNil)
			So(err, ShouldBeNil)
			bob, err := NewPeerConnection(config)
			So(bob, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(alice.SignalingState(), ShouldEqual, SignalingStateStable)
			So(bob.SignalingState(), ShouldEqual, SignalingStateStable)
			So(alice.ConnectionState(), ShouldEqual, PeerConnectionStateNew)
			So(bob.ConnectionState(), ShouldEqual, PeerConnectionStateNew)

			Convey("Alice creates offer", func() {
				offer, err = alice.CreateOffer()
				So(offer, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(alice.SignalingState(), ShouldEqual, SignalingStateStable)
				So(alice.ConnectionState(), ShouldEqual, PeerConnectionStateNew)

				// Shouldn't be able to set nil SDPs.
				err = alice.SetLocalDescription(nil)
				So(err, ShouldNotBeNil)
				err = alice.SetRemoteDescription(nil)
				So(err, ShouldNotBeNil)

				// Shouldn't be able to CreateAnswer
				answer, err = alice.CreateAnswer()
				So(answer, ShouldBeNil)
				So(err, ShouldNotBeNil)

				err = alice.SetLocalDescription(offer)
				So(err, ShouldBeNil)
				So(alice.LocalDescription(), ShouldEqual, offer)
				So(alice.SignalingState(), ShouldEqual, SignalingStateHaveLocalOffer)

				ic := IceCandidate{"fixme", "", 0}
				err = alice.AddIceCandidate(ic)
				So(err, ShouldNotBeNil)

				Convey("Bob receive offer and generates answer", func() {
					err = bob.SetRemoteDescription(offer)
					So(err, ShouldBeNil)
					So(bob.RemoteDescription(), ShouldEqual, offer)
					So(bob.SignalingState(), ShouldEqual, SignalingStateHaveRemoteOffer)

					answer, err = bob.CreateAnswer()
					So(answer, ShouldNotBeNil)
					So(err, ShouldBeNil)

					err = bob.SetLocalDescription(answer)
					So(err, ShouldBeNil)
					So(bob.LocalDescription(), ShouldEqual, answer)
					So(bob.SignalingState(), ShouldEqual, SignalingStateStable)

					Convey("Alice receives Bob's answer", func() {
						err = alice.SetRemoteDescription(answer)
						So(err, ShouldBeNil)
						So(alice.RemoteDescription(), ShouldEqual, answer)
						So(alice.SignalingState(), ShouldEqual, SignalingStateStable)
					})

				})
			})

			Convey("DataChannel", func() {
				channel, err := alice.CreateDataChannel("test", data.Init{})
				So(channel, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(channel.Label(), ShouldEqual, "test")
				channel.Close()
			})

			Convey("Close PeerConnections.", func() {
				success := make(chan int, 1)
				go func() {
					// err = alice.Close()
					// So(err, ShouldBeNil)
					err = bob.Close()
					// So(err, ShouldBeNil)
					success <- 1
				}()
				// TODO: Check the signaling state.
				select {
				case <-success:
				case <-time.After(time.Second * 2):
					WARN.Println("Timed out... something's probably amiss.")
					success <- 0
				}
			})

		})
	})
}

func TestConfiguration(t *testing.T) {

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

func TestIceCandidate(t *testing.T) {
	Convey("IceCandidate", t, func() {

		Convey("Serialize and Deserialize", func() {
			ice := IceCandidate{
				"fake",
				"not real",
				1337,
			}
			expected := `{"candidate":"fake","sdpMid":"not real","sdpMLineIndex":1337}`
			So(ice.Serialize(), ShouldEqual, expected)

			r := DeserializeIceCandidate(`
    		{"candidate":"still fake","sdpMid":"illusory","sdpMLineIndex":1337}`)
			So(r, ShouldNotBeNil)
			So(r.Candidate, ShouldEqual, "still fake")
			So(r.SdpMid, ShouldEqual, "illusory")
			So(r.SdpMLineIndex, ShouldEqual, 1337)

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
