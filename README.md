# go-webrtc

[![Build Status](https://travis-ci.org/keroserene/go-webrtc.svg?branch=master)](https://travis-ci.org/keroserene/go-webrtc)

WebRTC for Golang.

### Current Status:

This repository is currently fluctuating a lot, and the exposed interfaces will
change. **Do not rely** on anything in here yet!

- A PeerConnection can be successfully established between two separate machines
  using this Go library.
- It is possible to exchange bytes over a real DTLS/SCTP datachannel. (See the
  chat demo)
- Video/Audio support from the Media API is not implemented as it's low priority
  for us -- but pull requests will be gladly taken!

There is still lots of work to do!

### Usage

To immediately see some action, try the chat demo from two machines (or one...)

- `git clone https://github.com/keroserene/go-webrtc`
- `cd go-webrtc`
- `go run demo/chat/chat.go`

Type "start" in one of the Peers, and copy the session descriptions.
(This is the "copy-paste" signalling channel). If ICE negotiation succeeds,
a really janky chat session should begin.


To write Go code which requires WebRTC functionality:
```
import "github.com/keroserene/go-webrtc/"
```
And then you can do things like `webrtc.NewPeerConnection(...)`.

If you've never used WebRTC before, there is already plenty of information
online along with javascript examples, but for the Go code here, take a look
within `demo/*` for real usage examples which show how to prepare a
PeerConnection and set up the necessary callbacks and signaling.

Also, here are the [GoDocs](https://godoc.org/github.com/keroserene/go-webrtc).

### Dependencies:

- GCC 5+
- TODO:

#### Package naming

The package name is `webrtc`, even though the repo name is `go-webrtc`.
(This may be slightly contrary to Go convention, unless we consider the suffix
to really begin at the last dash. Reasons:
- Dashes aren't allowed in package names
- Including the word "go" in a Go package name seems redundant
- Just calling this repo `webrtc` wouldn't make sense either.
- Also you can rename imported packages to whatever you like.

(e.g. `import "foo" "github.com/keroserene/go-webrtc"`)

### Building

Latest tested native webrtc archive: `f33698296719f956497d2dbff81b5080864a8804`

There are currently two ways to build gowebrtc: the easy way, and the hard way.

The hard way is to build from scratch, which involves Google's
depot_tools and chromium stuff, gclient syncing, which takes a couple
hours, and possibly many more if you run into problems... along with
writing a custom ninja file and concatenating archives correctly and such.

See [webrtc.org native-code dev](http://webrtc.org/native-code/development/).

The easy way is to use the pre-built archive I've provided in `lib/`.

Once the archive is ready, cgo takes care of everything, and building
is as easy as `go build` or `go install`.

TODO(keroserene): More information / provide a real build script to automate
the hard way so it becomes the easy way.
(See [Issue #23](https://github.com/keroserene/go-webrtc/issues/23))
