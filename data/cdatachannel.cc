#include "cdatachannel.h"
#include "talk/app/webrtc/datachannel.h"
#include "webrtc/base/common.h"
#include <iostream>
#include <stdbool.h>

using namespace std;
using namespace webrtc;

const char *CGO_Channel_Label(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  // assert(NULL != dc);
  if (NULL == dc) {
    return "No internal CGO_Channel.";
  }
  return dc->label().c_str();
}

const bool CGO_Channel_Ordered(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  return dc->ordered();
}

int CGO_Channel_ReadyState(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  return dc->state();
}

const int CGO_DataStateConnecting = DataChannelInterface::DataState::kConnecting;
const int CGO_DataStateOpen = DataChannelInterface::DataState::kOpen;
const int CGO_DataStateClosing = DataChannelInterface::DataState::kClosing;
const int CGO_DataStateClosed = DataChannelInterface::DataState::kClosed;
