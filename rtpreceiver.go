package webrtc

/*
#include "refptr.h"
#include "rtpreceiver.h"
#include "audiotrack.h"
*/
import "C"

type RtpReceiver struct {
	p *refptr
	r C.CGO_RtpReceiver
}

func newRtpReceiver(r C.CGO_RtpReceiver) *RtpReceiver {
	return &RtpReceiver{
		p: newRefPtr(C.CGO_RefPtr(r)),
		r: r,
	}
}

func (r *RtpReceiver) Track() MediaStreamTrack {
	isAudio := false
	t := C.CGO_RtpReceiver_Track(r.r, (*C.bool)(&isAudio))
	if isAudio {
		return newAudioTrack(C.CGO_AudioTrack(t))
	}
	panic("VideoTrack not yet implemented")
}
