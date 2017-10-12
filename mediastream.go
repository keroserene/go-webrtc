package webrtc

/*
#include "refptr.h"
#include "mediastream.h"
#include "peerconnection.h"

#include <stdlib.h>
*/
import "C"
import "unsafe"

type MediaStream struct {
	p *refptr
	s C.CGO_MediaStream
}

func newMediaStream(s C.CGO_MediaStream) *MediaStream {
	return &MediaStream{
		p: newRefPtr(C.CGO_RefPtr(s)),
		s: s,
	}
}

// TODO: Factor pc_factory out of Peer (and make it a single global instance?)
func (pc *PeerConnection) NewMediaStream(label string) *MediaStream {
	s := C.CGO_NewMediaStream(pc.cgoPeer, C.CString(label))
	return newMediaStream(s)
}

func (s *MediaStream) GetAudioTracks() []*AudioTrack {
	var n C.int
	ctracks := uintptr(unsafe.Pointer(C.CGO_MediaStream_GetAudioTracks(s.s, &n)))
	tracks := make([]*AudioTrack, n)
	for i := range tracks {
		t := *(*C.CGO_AudioTrack)(unsafe.Pointer(ctracks + uintptr(i)))
		tracks[i] = newAudioTrack(t)
	}
	C.free(unsafe.Pointer(ctracks))
	return tracks
}

func (s *MediaStream) AddTrack(t MediaStreamTrack) {
	switch t := t.(type) {
	case *AudioTrack:
		C.CGO_MediaStream_AddAudioTrack(s.s, t.t)
	default:
		panic("unknown track type")
	}
}
