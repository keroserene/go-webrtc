// Go wrapper around libwebrtc.
// TODO
package webrtc

import (
  "fmt"
)


type RTCPeerConnection struct {
  // createOffer
  // createAnswer
  // setLocalDescription
  // localDescription

  // currentLocalDescription
  // pendingLocalDescription
  // setRemoteDescription

  // remoteDescription
  // currentRemoteDescription
  // pendingRemoteDescription

  // addIceCandidate
  // signalingState
  // iceGatheringState
  // iceConnectionState
  // canTrickleIceCandidates
  // getConfiguration
  // setConfiguration
  // close

  // onnegotiationneeded
  // onicecandidate
  // onsignalingstatechange
  // oniceconnectionstatechange
  // onicegatheringstatechange
}

type RTCSignalingState int
/*
const {
  stable RTCSignallingState = iota
  have-local-offer
  have-remote-offer
  have-local-pranswer
  have-remote-pranswer
  closed
}
*/

type RTCConfiguration struct {
}

type RTCIceCredentialType struct {
}

type RTCIceServer struct {
  urls        string
  username    string
  credential  string
  // credentialType   RTCIceCredentialType
}

type RTCIceTransportPolicy struct {
}

type RTCBundlePolicy struct {
}

func init() {
  fmt.Printf("yet to be implemented")
}
