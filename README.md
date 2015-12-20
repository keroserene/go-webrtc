# go-webrtc

[![Circle CI](https://circleci.com/gh/keroserene/go-webrtc.svg?style=svg)](https://circleci.com/gh/keroserene/go-webrtc)

This repository is currently **not ready** and fluctuating a lot.
**Do not rely** on anything in here yet!

This currently only builds on linux, but OSX is in progress.
Actual documentation is on the way.

### Current Status:

- The Go code successfully wraps the C++ code.
- From Go, it is now possible to create a WebRTC PeerConnection, generate SDP
messages and ICE candidates, and exchange bytes over a real DataChannel.
- Video/Audio support from the Media API is not implemented as it's low priority
  for us -- but pull requests will be gladly taken :)

There is still lots of work to do!

### Usage

In .go files that requires WebRTC functionality:
```
import "github.com/keroserene/go-webrtc/"
```
And then you can do things like `webrtc.NewPeerConnection(...)`.

If you've never used WebRTC before, there is already plenty of information
online along with javascript examples, but for the Go code here, take a look
within `demo/*` for real usage examples which show how to prepare a
PeerConnection and set up the necessary callbacks and signaling.

Also, here are the [GoDocs](https://godoc.org/github.com/keroserene/go-webrtc).

#### Package naming

The package name is `webrtc`, even though the repo name is `go-webrtc`.
(This may be slightly contrary to Go convention, unless we consider the suffix
to really begin at the last dash. Reasons:
- Dashes aren't allowed in package names
- Including the word "go" in a Go package name seems redundant
- Just calling this repo `webrtc` wouldn't make sense either.
- Also you can rename imported packages to whatever you like.
(Something like `import foo github.com/keroserene/go-webrtc`)

### Dependencies / Building

To build this from scratch is currently challenging, due to the way building
webrtc is.
See [webrtc.org native-code dev](http://webrtc.org/native-code/development/) pages.

Once the archive is prepared though, cgo takes care of everything, and building
is as easy as `go build` or `go install`.
If you'd like to try the chat demo, you can do `go run demo/chat/chat.go`.

TODO(keroserene): More information / make things easier.

### Conventions

- Please run `go fmt` before every commit.

- There is a `.CGO()` method for every Go struct which expects being passed to
  native code.



