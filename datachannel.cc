#include "datachannel.h"
#include "datachannel.hpp"

#include <stdbool.h>

#include "talk/app/webrtc/test/fakeconstraints.h"
#include "webrtc/base/common.h"

using namespace webrtc;

// Create and register a new DataChannelObserver.
CGO_Channel CGO_Channel_RegisterObserver(void *o, void *goChannel) {
  auto obs = (CGoDataChannelObserver*)o;
  obs->goChannel = goChannel;
  return obs->dc.get();
}

void CGO_Channel_Send(CGO_Channel channel, void *data, int size) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  auto bytes = rtc::Buffer((uint8_t*)data, size);
  auto buffer = DataBuffer(bytes, true);
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
  char *cstr = (char *)malloc(dc->label().length()+1 * sizeof(char));
  std::strcpy(cstr, dc->label().c_str());
  return cstr;
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
  char *cstr = (char *)malloc(dc->protocol().length()+1 * sizeof(char));
  std::strcpy(cstr, dc->protocol().c_str());
  return cstr;
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
// PeerConnection, which is outside the scope of this subpackage.
// However, we can still create a fake DataChannelInterface for testing.

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

rtc::scoped_refptr<CGoDataChannelObserver> test_observer;

void* CGO_getFakeDataChannel() {
  rtc::scoped_refptr<FakeDataChannel> test_dc = new rtc::RefCountedObject<FakeDataChannel>();
  test_observer = new rtc::RefCountedObject<CGoDataChannelObserver>(test_dc);
  auto o = test_observer.get();
  test_dc->RegisterObserver(o);
  return (void *)o;
}

void CGO_fakeMessage(CGO_Channel channel, void *data, int size) {
  auto dc = (FakeDataChannel*)channel;
  auto bytes = rtc::Buffer((char*)data, size);
  auto buffer = DataBuffer(bytes, true);
  dc->Send(buffer);
}

void CGO_fakeStateChange(CGO_Channel channel, int state) {
  auto dc = (FakeDataChannel*)channel;
  dc->SetState((DataChannelInterface::DataState)state);
}

void CGO_fakeBufferAmount(CGO_Channel channel, int amount) {
  auto dc = (FakeDataChannel*)channel;
  dc->SetBufferedAmount(amount);
}

const int CGO_DataStateConnecting = DataChannelInterface::DataState::kConnecting;
const int CGO_DataStateOpen = DataChannelInterface::DataState::kOpen;
const int CGO_DataStateClosing = DataChannelInterface::DataState::kClosing;
const int CGO_DataStateClosed = DataChannelInterface::DataState::kClosed;
