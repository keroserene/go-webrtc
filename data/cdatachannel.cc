#include "cdatachannel.h"
#include "talk/app/webrtc/datachannel.h"
#include "webrtc/base/common.h"
#include <iostream>

using namespace std;

const char *CGO_Channel_Label(CGO_Channel channel) {
  auto dc = (webrtc::DataChannelInterface*)channel;
  // assert(NULL != dc);
  if (NULL == dc) {
    return "No internal CGO_Channel.";
  }
  return dc->label().c_str();
}
