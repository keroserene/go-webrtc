package webrtc

/*
#include "eventlistener.h"
#include "refptr.h"
#include "mediastreamtrack.h"
*/
import "C"

// A MediaStreamTrack is either an *AudioTrack or a *VideoTrack.
type MediaStreamTrack interface {
	ID() string
	Enabled() bool
	SetEnabled(bool)
	Ended() bool
	OnEnded(func()) *EventListener

	cgoMediaStreamTrack() C.CGO_MediaStreamTrack
}

// A mediaStreamTrack is the common implementation shared by AudioTrack and VideoTrack.
type mediaStreamTrack struct {
	p *refptr
	t C.CGO_MediaStreamTrack
}

func newMediaStreamTrack(t C.CGO_MediaStreamTrack) *mediaStreamTrack {
	return &mediaStreamTrack{
		p: newRefPtr(C.CGO_RefPtr(t)),
		t: t,
	}
}

func (t *mediaStreamTrack) cgoMediaStreamTrack() C.CGO_MediaStreamTrack {
	return t.t
}

// ID returns id of the media stream track
func (t *mediaStreamTrack) ID() string {
	return C.GoString(C.CGO_MediaStreamTrack_ID(t.t))
}

// Enabled returns enabled flag of the media stream track
func (t *mediaStreamTrack) Enabled() bool {
	return bool(C.CGO_MediaStreamTrack_Enabled(t.t))
}

// SetEnabled enable or disable media stream track
func (t *mediaStreamTrack) SetEnabled(x bool) {
	C.CGO_MediaStreamTrack_SetEnabled(t.t, C.bool(x))
}

// Ended returns ended flag of one media stream
func (t *mediaStreamTrack) Ended() bool {
	return bool(C.CGO_MediaStreamTrack_Ended(t.t))
}

// OnEnded pass ended event to callback
func (t *mediaStreamTrack) OnEnded(f func()) *EventListener {
	var e *EventListener
	e = newEventListener(C.CGO_Notifier(t.t), func() {
		if t.Ended() {
			f()
			// Once a track has ended, it will never become live again.
			e.Cancel()
		}
	})
	return e
}
