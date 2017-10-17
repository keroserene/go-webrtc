#include "fauxaudiodevicemodule.hpp"

#include "webrtc/base/refcount.h"

int32_t FauxAudioDeviceModule::RegisterAudioCallback(webrtc::AudioTransport* audio_callback) {
  rtc::CritScope cs(&crit_);
  audio_callback_ = audio_callback;
  return 0;
}

void FauxAudioDeviceModule::Pull(uint16_t* audioSamples, size_t nSamples, uint32_t samplesPerSec) {
  rtc::CritScope cs(&crit_);
  if (!audio_callback_) {
    return;
  }
  size_t nSamplesOut = 0;
  int64_t elapsed_time_ms = 0;
  int64_t ntp_time_ms = 0;
  audio_callback_->NeedMorePlayData(nSamples, 2, 1, samplesPerSec,
                                    audioSamples, nSamplesOut,
                                    &elapsed_time_ms, &ntp_time_ms);
}
