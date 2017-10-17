#include <_cgo_export.h>  // Allow calling certain Go functions.

#include "eventlistener.h"

#include "webrtc/api/mediastreaminterface.h"

using namespace webrtc;

class Observer : public ObserverInterface {
public:
  Observer(NotifierInterface* n, CGO_EventCallback c)
    : n(n), c(c) {
    n->RegisterObserver(this);
  }

  ~Observer() {
    n->UnregisterObserver(this);
  }

  void OnChanged() override {
    cgoObserverOnChanged(c);
  }

private:
  NotifierInterface* n;
  CGO_EventCallback c;
};

CGO_Observer CGO_NewObserver(CGO_Notifier n, CGO_EventCallback c) {
  return new Observer((NotifierInterface*)n, c);
}

void CGO_DeleteObserver(CGO_Observer o) {
  delete (Observer*)o;
}
