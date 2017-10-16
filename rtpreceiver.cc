#include "rtpreceiver.h"

#include "webrtc/api/rtpreceiverinterface.h"

CGO_MediaStreamTrack CGO_RtpReceiver_Track(CGO_RtpReceiver r, bool* isAudio) {
	auto t = ((webrtc::RtpReceiverInterface*)r)->track().release();
	*isAudio = t->kind() == "audio";
	return t;
}
