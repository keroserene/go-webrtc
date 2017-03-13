/*
 *  Copyright 2017 The WebRTC project authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a BSD-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */

#ifndef WEBRTC_ORTC_RTPTRANSPORTADAPTER_H_
#define WEBRTC_ORTC_RTPTRANSPORTADAPTER_H_

#include <memory>
#include <vector>

#include "webrtc/api/ortc/rtptransportinterface.h"
#include "webrtc/api/rtcerror.h"
#include "webrtc/base/constructormagic.h"
#include "webrtc/base/sigslot.h"
#include "webrtc/media/base/streamparams.h"
#include "webrtc/ortc/rtptransportcontrolleradapter.h"
#include "webrtc/pc/channel.h"

namespace webrtc {

// Implementation of RtpTransportInterface to be used with RtpSenderAdapter,
// RtpReceiverAdapter, and RtpTransportControllerAdapter classes.
//
// TODO(deadbeef): When BaseChannel is split apart into separate
// "RtpTransport"/"RtpTransceiver"/"RtpSender"/"RtpReceiver" objects, this
// adapter object can be removed.
class RtpTransportAdapter : public RtpTransportInterface {
 public:
  // |rtp| can't be null. |rtcp| can if RTCP muxing is used immediately (meaning
  // |rtcp_parameters.mux| is also true).
  static RTCErrorOr<std::unique_ptr<RtpTransportInterface>> CreateProxied(
      const RtcpParameters& rtcp_parameters,
      PacketTransportInterface* rtp,
      PacketTransportInterface* rtcp,
      RtpTransportControllerAdapter* rtp_transport_controller);
  ~RtpTransportAdapter() override;

  // RtpTransportInterface implementation.
  PacketTransportInterface* GetRtpPacketTransport() const override;
  PacketTransportInterface* GetRtcpPacketTransport() const override;
  RTCError SetRtcpParameters(const RtcpParameters& parameters) override;
  RtcpParameters GetRtcpParameters() const override { return rtcp_parameters_; }

  // Methods used internally by OrtcFactory.
  RtpTransportControllerAdapter* rtp_transport_controller() {
    return rtp_transport_controller_;
  }
  void TakeOwnershipOfRtpTransportController(
      std::unique_ptr<RtpTransportControllerInterface> controller);

  // Used by RtpTransportControllerAdapter to tell when it should stop
  // returning this transport from GetTransports().
  sigslot::signal1<RtpTransportAdapter*> SignalDestroyed;

 protected:
  RtpTransportAdapter* GetInternal() override { return this; }

 private:
  RtpTransportAdapter(const RtcpParameters& rtcp_parameters,
                      PacketTransportInterface* rtp,
                      PacketTransportInterface* rtcp,
                      RtpTransportControllerAdapter* rtp_transport_controller);

  PacketTransportInterface* rtp_packet_transport_;
  PacketTransportInterface* rtcp_packet_transport_;
  RtpTransportControllerAdapter* rtp_transport_controller_;
  // Non-null if this class owns the transport controller.
  std::unique_ptr<RtpTransportControllerInterface>
      owned_rtp_transport_controller_;
  RtcpParameters rtcp_parameters_;

  RTC_DISALLOW_IMPLICIT_CONSTRUCTORS(RtpTransportAdapter);
};

}  // namespace webrtc

#endif  // WEBRTC_ORTC_RTPTRANSPORTADAPTER_H_
