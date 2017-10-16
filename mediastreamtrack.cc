#include <_cgo_export.h>  // Allow calling certain Go functions.

#include "mediastreamtrack.h"

#include "webrtc/api/mediastreaminterface.h"

using namespace webrtc;

const char* CGO_MediaStreamTrack_ID(CGO_MediaStreamTrack t) {
	return ((MediaStreamTrackInterface*)t)->id().c_str();
}

bool CGO_MediaStreamTrack_Enabled(CGO_MediaStreamTrack t) {
	return ((MediaStreamTrackInterface*)t)->enabled();
}

void CGO_MediaStreamTrack_SetEnabled(CGO_MediaStreamTrack t, bool x) {
	((MediaStreamTrackInterface*)t)->set_enabled(x);
}

bool CGO_MediaStreamTrack_Ended(CGO_MediaStreamTrack t) {
	return ((MediaStreamTrackInterface*)t)->state() == MediaStreamTrackInterface::TrackState::kEnded;
}
