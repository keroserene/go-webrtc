#include "refptr.h"

#include "webrtc/base/refcount.h"

void CGO_RefPtr_Release(CGO_RefPtr p) {
	((rtc::RefCountInterface*)p)->Release();
}
