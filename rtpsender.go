package webrtc

/*
#include "refptr.h"
#include "rtpsender.h"
#include "audiotrack.h"
*/
import "C"

type RtpSender struct {
	p *refptr
	s C.CGO_RtpSender
}

func newRtpSender(s C.CGO_RtpSender) *RtpSender {
	return &RtpSender{
		p: newRefPtr(C.CGO_RefPtr(s)),
		s: s,
	}
}

func (s *RtpSender) Track() MediaStreamTrack {
	isAudio := false
	t := C.CGO_RtpSender_Track(s.s, (*C.bool)(&isAudio))
	if isAudio {
		return newAudioTrack(C.CGO_AudioTrack(t))
	}
	panic("VideoTrack not yet implemented")
}
