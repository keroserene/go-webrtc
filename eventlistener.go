package webrtc

/*
#include "eventlistener.h"
*/
import "C"

type EventListener struct {
	o C.CGO_Observer
	c int
}

func newEventListener(n C.CGO_Notifier, f func()) *EventListener {
	c := eventCallbacks.Set(f)
	return &EventListener{
		o: C.CGO_NewObserver(n, C.CGO_EventCallback(c)),
		c: c,
	}
}

var eventCallbacks = NewCGOMap()

func (e *EventListener) Cancel() {
	if e.o != nil {
		C.CGO_DeleteObserver(e.o)
		e.o = nil
		eventCallbacks.Delete(e.c)
	}
}

//export cgoObserverOnChanged
func cgoObserverOnChanged(l C.CGO_EventCallback) {
	eventCallbacks.Get(int(l)).(func())()
}
