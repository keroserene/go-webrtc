package webrtc

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIceGatheringStateEnums(t *testing.T) {
	Convey(`Enum: IceGatheringState values should match
C++ webrtc::PeerConnectionInterface values`, t, func() {
		So(IceGatheringStateNew, ShouldEqual, _cgoIceGatheringStateNew)
		So(IceGatheringStateGathering, ShouldEqual, _cgoIceGatheringStateGathering)
		So(IceGatheringStateComplete, ShouldEqual, _cgoIceGatheringStateComplete)
	})
}

func TestIceConnectionStateEnums(t *testing.T) {
	Convey(`Enum: IceConnectionState values should match
C++ webrtc::PeerConnectionInterface values`, t, func() {
		So(IceConnectionStateNew, ShouldEqual, _cgoIceConnectionStateNew)
		So(IceConnectionStateChecking, ShouldEqual, _cgoIceConnectionStateChecking)
		So(IceConnectionStateConnected, ShouldEqual,
			_cgoIceConnectionStateConnected)
		So(IceConnectionStateCompleted, ShouldEqual,
			_cgoIceConnectionStateCompleted)
		So(IceConnectionStateFailed, ShouldEqual, _cgoIceConnectionStateFailed)
		So(IceConnectionStateDisconnected, ShouldEqual,
			_cgoIceConnectionStateDisconnected)
		So(IceConnectionStateClosed, ShouldEqual, _cgoIceConnectionStateClosed)
	})
}

func TestPeerConnection(t *testing.T) {
	SetLoggingVerbosity(0)

	Convey("PeerConnection", t, func() {
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
			So(pc.IceGatheringState(), ShouldEqual, IceGatheringStateNew)
			So(pc.IceConnectionState(), ShouldEqual, IceConnectionStateNew)

			Convey("Set and Get Configuration", func() {
				config := NewConfiguration(
					OptionIceServer("stun:something.else"),
					OptionIceTransportPolicy(IceTransportPolicyRelay))
				So(config.IceTransportPolicy, ShouldEqual, IceTransportPolicyRelay)
				// FIXME: Should be calling SetConfiguration here
				pc, err = NewPeerConnection(config)
				So(err, ShouldBeNil)
				got := pc.GetConfiguration()
				So(got.IceTransportPolicy, ShouldEqual, IceTransportPolicyRelay)
			})

			Convey("Callbacks fire correctly", func() {

				Convey("OnNegotiationNeeded", func() {
					success := make(chan int, 1)
					pc.OnNegotiationNeeded = func() {
						success <- 0
					}
					cgoOnNegotiationNeeded(pc.index)
					select {
					case <-success:
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

				Convey("OnIceCandidateError", func() {
					success := make(chan int, 1)
					pc.OnIceCandidateError = func() {
						success <- 1
					}
					cgoFakeIceCandidateError(pc)
					select {
					case <-success:
					case <-time.After(time.Second * 1):
						t.Fatal("Timed out.")
					}
				})

				Convey("OnSignalingState", func() {
					success := make(chan SignalingState, 1)
					pc.OnSignalingStateChange = func(s SignalingState) {
						success <- s
					}
					cgoOnSignalingStateChange(pc.index, SignalingStateStable)
					select {
					case state := <-success:
						So(state, ShouldEqual, SignalingStateStable)
					case <-time.After(time.Second * 1):
						t.Fatal("Timed out.")
					}
				})

				Convey("OnIceConnectionStateChange", func() {
					success := make(chan IceConnectionState, 1)
					pc.OnIceConnectionStateChange = func(state IceConnectionState) {
						success <- state
					}
					cgoOnIceConnectionStateChange(pc.index,
						IceConnectionStateDisconnected)
					select {
					case r := <-success:
						So(r, ShouldEqual, IceConnectionStateDisconnected)
					case <-time.After(time.Second * 1):
						t.Fatal("Timed out.")
					}
				})

				Convey("OnConnectionStateChange", func() {
					success := make(chan PeerConnectionState, 1)
					expectPeerConnectionState := func(state PeerConnectionState) {
						select {
						case r := <-success:
							So(r, ShouldEqual, state)
						case <-time.After(time.Second * 1):
							t.Fatal("Timed out.")
						}
					}
					pc.OnConnectionStateChange = func(state PeerConnectionState) {
						success <- state
					}
					cgoOnConnectionStateChange(pc.index,
						IceConnectionStateNew)
					expectPeerConnectionState(PeerConnectionStateNew)
					cgoOnConnectionStateChange(pc.index,
						IceConnectionStateConnected)
					expectPeerConnectionState(PeerConnectionStateConnected)
					cgoOnConnectionStateChange(pc.index,
						IceConnectionStateFailed)
					expectPeerConnectionState(PeerConnectionStateFailed)
					cgoOnConnectionStateChange(pc.index,
						IceConnectionStateDisconnected)
					expectPeerConnectionState(PeerConnectionStateDisconnected)
				})

				Convey("OnDataChannel", func() {
					success := make(chan *DataChannel, 1)
					pc.OnDataChannel = func(dc *DataChannel) {
						success <- dc
					}
					cgoOnDataChannel(pc.index, nil)
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
				offer, err := alice.CreateOffer()
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
				answer, err := alice.CreateAnswer()
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
				channel, err := alice.CreateDataChannel("test")
				So(channel, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(channel.Label(), ShouldEqual, "test")
				alice.DeleteDataChannel(channel)
			})

			Convey("Destroy PeerConnections.", func() {
				success := make(chan int, 1)
				go func() {
					// err = alice.Destroy()
					// So(err, ShouldBeNil)
					err = bob.Destroy()
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
