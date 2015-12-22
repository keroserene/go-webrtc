#include "cdatachannel.h"
#include "talk/app/webrtc/datachannel.h"
#include "talk/app/webrtc/peerconnectioninterface.h"
#include "talk/app/webrtc/test/fakeconstraints.h"
#include "webrtc/base/common.h"
#include <iostream>
#include <stdbool.h>
#include "_cgo_export.h"

using namespace std;
using namespace webrtc;

class CGoDataChannelObserver : public DataChannelObserver {
 public:
  CGoDataChannelObserver(void *goPtr) : goChannel(goPtr) {
    assert(NULL != goChannel);
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

 protected:

  // Reference to external Go data.Channel required for callbacks.
  void *goChannel;

  ~CGoDataChannelObserver() {
    cout << "[C] CgoDataChannelObserver destructing." << endl;
  }
};  // class DoDataChannelObserver

// Create and register a new DataChannelObserver.
void CGO_Channel_RegisterObserver(CGO_Channel channel, void *goChannel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  auto obs = new CGoDataChannelObserver(goChannel);
  dc->RegisterObserver(obs);
}

void CGO_Channel_Send(CGO_Channel channel, void *data, int size) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  auto bytes = new rtc::Buffer((uint8_t*)data, size);
  auto buffer = DataBuffer(*bytes, true);
  dc->Send(buffer);
}

void CGO_Channel_Close(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  assert(NULL != dc);
  return dc->Close();
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
 public:
  virtual void RegisterObserver(DataChannelObserver* observer) {
    obs_ = observer;
  };
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
    return state_;
  };

  virtual uint64_t buffered_amount() const {
    return 0;
  };

  // Sends data to self.
  bool Send(const DataBuffer& buffer) {
    obs_->OnMessage(buffer);
    return false;
  };

  void Close() {
    SetState(DataState::kClosed);
  };

  void SetState(DataChannelInterface::DataState state) {
    state_ = state;
    obs_->OnStateChange();
  }

  void SetBufferedAmount(int amount) {
    obs_->OnBufferedAmountChange(amount);
  }

 protected:
  DataChannelObserver* obs_;
  DataState state_ = DataState::kClosed;
};
rtc::scoped_refptr<FakeDataChannel> test_dc;

CGO_Channel CGO_getFakeDataChannel() {
  test_dc = new rtc::RefCountedObject<FakeDataChannel>();
  return (void *)test_dc;
}

void CGO_fakeMessage(CGO_Channel channel, void *data, int size) {
  auto bytes = new rtc::Buffer((char*)data, size);
  auto dc = (webrtc::DataChannelInterface*)channel;
  auto buffer = DataBuffer(*bytes, true);
  dc->Send(buffer);
}

void CGO_fakeStateChange(CGO_Channel channel, int state) {
  test_dc->SetState((DataChannelInterface::DataState)state);
}

void CGO_fakeBufferAmount(CGO_Channel channel, int amount) {
  test_dc->SetBufferedAmount(amount);
}

const int CGO_DataStateConnecting = DataChannelInterface::DataState::kConnecting;
const int CGO_DataStateOpen = DataChannelInterface::DataState::kOpen;
const int CGO_DataStateClosing = DataChannelInterface::DataState::kClosing;
const int CGO_DataStateClosed = DataChannelInterface::DataState::kClosed;
