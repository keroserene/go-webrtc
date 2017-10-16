#include "mediastream.h"

#include "webrtc/api/mediastreaminterface.h"

CGO_AudioTrack* CGO_MediaStream_GetAudioTracks(CGO_MediaStream s, int* n) {
	auto tracks = ((webrtc::MediaStreamInterface*)s)->GetAudioTracks();
	*n = tracks.size();
	auto ctracks = (CGO_AudioTrack*)malloc(*n * sizeof(CGO_AudioTrack*));
	for (int i = 0; i < *n; ++i) {
		ctracks[i] = tracks[i].release();
	}
	return ctracks;
}

void CGO_MediaStream_AddAudioTrack(CGO_MediaStream s, CGO_AudioTrack t) {
	((webrtc::MediaStreamInterface*)s)->AddTrack((webrtc::AudioTrackInterface*)t);
}
