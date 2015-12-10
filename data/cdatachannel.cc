#include "cdatachannel.h"
#include "talk/app/webrtc/datachannel.h"
#include "talk/app/webrtc/peerconnectioninterface.h"
#include "talk/app/webrtc/test/fakeconstraints.h"
#include "webrtc/base/common.h"
#include <iostream>
#include <stdbool.h>

using namespace std;
using namespace webrtc;

class CGoDataChannelObserver : public DataChannelObserver {
 public:

  void OnStateChange() {
    cout << "[C] OnStateChange" << endl;
  }

  void OnMessage(const DataBuffer& buffer) {
    cout << "[C] OnMessage" << endl;
  }

  void OnBufferedAmountChange(uint64_t previous_amount) {
    cout << "[C] OnBufferedAmountChange" << endl;
  }

 protected:
  ~CGoDataChannelObserver() {
    cout << "[C] Destructing DataChannelObserver" << endl;
  }
};  // class DoDataChannelObserver

// Create and register a new DataChannelObserver.
const char *CGO_Channel_RegisterObserver(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  auto obs = new CGoDataChannelObserver();
  dc->RegisterObserver(obs);
}

const char *CGO_Channel_Label(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  assert(NULL != dc);
  return dc->label().c_str();
}

const bool CGO_Channel_Ordered(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  assert(NULL != dc);
  return dc->ordered();
}

int CGO_Channel_MaxRetransmitTime(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  assert(NULL != dc);
  return dc->maxRetransmitTime();
}

int CGO_Channel_MaxRetransmits(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  assert(NULL != dc);
  return dc->maxRetransmits();
}

const char *CGO_Channel_Protocol(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  assert(NULL != dc);
  return dc->protocol().c_str();
}

const bool CGO_Channel_Negotiated(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  assert(NULL != dc);
  return dc->negotiated();
}

int CGO_Channel_ID(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  assert(NULL != dc);
  return dc->id();
}

int CGO_Channel_ReadyState(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  assert(NULL != dc);
  return dc->state();
}

int CGO_Channel_BufferedAmount(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  assert(NULL != dc);
  return dc->buffered_amount();
}


//
// === Testing helpers ===
//

// Real DataChannels can only be created from a correctly configured
// PeerConnection, which is outside the scope of functionality this data
// subpackage. However, we can still need fake DataChannelInterface for testing.

class FakeDataChannel : public DataChannelInterface {
  virtual void RegisterObserver(DataChannelObserver* observer) {};
  virtual void UnregisterObserver() {};

  virtual std::string label() const {
    return "fake";
  };

  virtual bool reliable() const {
    return false;
  };

  virtual int id() const {
    return 12345;
  };

  virtual DataState state() const {
    return DataState::kClosed;
  };

  virtual uint64_t buffered_amount() const {
    return 0;
  };

  virtual bool Send(const DataBuffer& buffer) {
    return false;
  };

  virtual void Close() {};
};
rtc::scoped_refptr<FakeDataChannel> test_dc;

CGO_Channel CGO_getFakeDataChannel() {
  test_dc = new rtc::RefCountedObject<FakeDataChannel>();
  return (void *)test_dc;
}

const int CGO_DataStateConnecting = DataChannelInterface::DataState::kConnecting;
const int CGO_DataStateOpen = DataChannelInterface::DataState::kOpen;
const int CGO_DataStateClosing = DataChannelInterface::DataState::kClosing;
const int CGO_DataStateClosed = DataChannelInterface::DataState::kClosed;
