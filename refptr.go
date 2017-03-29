package webrtc

/*
#include "refptr.h"
*/
import "C"
import "runtime"

/*
refptr holds a reference counted pointer and decrements its reference count
when the refptr is garbage collected.  It is roughly the equivalent of
rtc::scoped_refptr, except it doesn't increment the reference count.
*/
type refptr struct{}

/*
newRefPtr returns a refptr that, upon being garbage collected, decrements p's
reference count.  p must already have its reference count incremented before
being passed to newRefPtr, for example by calling scoped_refptr::release().
*/
func newRefPtr(p C.CGO_RefPtr) *refptr {
	rp := &refptr{}
	runtime.SetFinalizer(rp, func(*refptr) {
		C.CGO_RefPtr_Release(p)
	})
	return rp
}
