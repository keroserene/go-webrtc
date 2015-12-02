package webrtc

import (
  "testing"
  "fmt"
)

// Prepare a PeerConnection client, and create an offer.
func TestCreateOffer(t *testing.T) {
  fmt.Println("Test - [RTCPeerConnection]: CreateOffer")
  r := make(chan bool, 1)

  // Pretend to be the "Signalling" thread.
  go func() {

    pc, err := NewPeerConnection()
    if (nil != err) {
      t.Error(err)
      r <- false
      t.FailNow()
    }
    fmt.Printf("PeerConnection: %+v\n", pc)

    success := func () {
      fmt.Println("CreateOffer Succeeded.")
    }
    failure := func () {
      fmt.Println("CreateOffer Failed.")
    }

    pc.CreateOffer(success, failure)
    r <- true
  }()
  <-r
  fmt.Println("Done")
}

// TODO: Test video / audio stream support.
