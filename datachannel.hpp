#ifndef _DATACHANNEL_H_
#define _DATACHANNEL_H_

#include <_cgo_export.h>  // Allow calling certain Go functions.
#include <assert.h>

#include "webrtc/api/peerconnectioninterface.h"
#include "webrtc/api/datachannelinterface.h"

typedef rtc::scoped_refptr<webrtc::DataChannelInterface> DataChannel;

class CGoDataChannelObserver
  : public webrtc::DataChannelObserver,
    public rtc::RefCountInterface {
 public:
  explicit CGoDataChannelObserver(DataChannel dc) : dc(dc) {
    assert(NULL != dc);
  }

  void OnStateChange() {
    cgoChannelOnStateChange(goChannel);
  }

  void OnMessage(const webrtc::DataBuffer& buffer) {
    auto data = reinterpret_cast<void*>(const_cast<unsigned char*>(
                                        buffer.data.data()));
    cgoChannelOnMessage(goChannel, data, buffer.size());
  }

  void OnBufferedAmountChange(uint64_t previous_amount) {
    cgoChannelOnBufferedAmountChange(goChannel, previous_amount);
  }

  // Reference to external Go data.Channel required for callbacks.
  int goChannel;
  DataChannel dc;

 protected:
  ~CGoDataChannelObserver() {
    dc->UnregisterObserver();
  }
};  // class DoDataChannelObserver

#endif  // _DATACHANNEL_H
