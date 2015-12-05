# go-webrtc

[![Circle CI](https://circleci.com/gh/keroserene/go-webrtc.svg?style=svg)](https://circleci.com/gh/keroserene/go-webrtc)

This repository is currently fluctuating a lot.
This currently only builds on linux, but OSX is in progress.
Actual documentation is on the way.

## Current Status:

- The Go code successfully wraps the C++ code.
- From Go, you can now create a WebRTC PeerConnection, create SDP messages, and
create DataChannels.

There is still lots of work to do!


## Usage in Go code
In .go files that requires WebRTC functionality:
```

import "github.com/keroserene/go-webrtc/"
```
And then you can do things like `webrtc.NewPeerConnection(...)`.

The package name is `webrtc`, even though the repo name is `go-webrtc`.
(This may be slightly contrary to Go convention, unless we consider the suffix
to really begin at the last dash. Reasons:
- Dashes aren't allowed in package names
- Including the word "go" in a Go package name seems redundant
- Just calling this repo `webrtc` wouldn't make sense either.
- Also you can rename imported packages to whatever you like.

Look within `demo/*` for further usage examples.


## Conventions

- Please run `go fmt` before every commit.

- There is a `.CGO()` method for every Go struct which expects being passed to
  native code.


## Explanation of the cgo wrapper

TODO(keroserene)
