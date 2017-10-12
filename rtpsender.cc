#include "rtpsender.h"

#include "webrtc/api/rtpsenderinterface.h"

CGO_MediaStreamTrack CGO_RtpSender_Track(CGO_RtpSender s, bool* isAudio) {
	auto t = ((webrtc::RtpSenderInterface*)s)->track().release();
	*isAudio = t->kind() == "audio";
	return t;
}
