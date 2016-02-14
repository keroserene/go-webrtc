#ifndef _DATACHANNEL_H_
#define _DATACHANNEL_H_

#include "_cgo_export.h"  // Allow calling certain Go functions.

#include "talk/app/webrtc/peerconnectioninterface.h"
#include "talk/app/webrtc/datachannelinterface.h"

using namespace std;
using namespace webrtc;

typedef rtc::scoped_refptr<DataChannelInterface> DataChannel;

class CGoDataChannelObserver
  : public DataChannelObserver,
    public rtc::RefCountInterface {
 public:
  CGoDataChannelObserver(DataChannel dc) : dc(dc) {
    assert(NULL != dc);
  }

  void OnStateChange() {
    cgoChannelOnStateChange(goChannel);
  }

  void OnMessage(const DataBuffer& buffer) {
    auto data = (uint8_t*)buffer.data.data();
    cgoChannelOnMessage(goChannel, (void *)data, buffer.size());
  }

  void OnBufferedAmountChange(uint64_t previous_amount) {
    cgoChannelOnBufferedAmountChange(goChannel, previous_amount);
  }

  // Reference to external Go data.Channel required for callbacks.
  void *goChannel;
  DataChannel dc;

 protected:
  ~CGoDataChannelObserver() {}
};  // class DoDataChannelObserver

#endif  // _DATACHANNEL_H
