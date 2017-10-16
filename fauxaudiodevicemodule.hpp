#ifndef FAUXAUDIODEVICEMODULE_H_
#define FAUXAUDIODEVICEMODULE_H_

#include "webrtc/base/criticalsection.h"
#include "webrtc/modules/audio_device/include/fake_audio_device.h"

/*
FauxAudioDeviceModule drives the mechanism for receiving remote audio tracks by
having its Pull method called periodically (from Go, because I couldn't get
threading to compile in C++).  It does not do the other thing that an
AudioDeviceModule is intended to do, which is to push recorded audio, because,
strangely, pushing any such audio causes it to be interleaved with audio data
from other sent tracks, corrupting it on the receiving end.

This feels hackish.  Maybe there is a better way to do it.
*/
class FauxAudioDeviceModule : public webrtc::FakeAudioDeviceModule {
 public:
  explicit FauxAudioDeviceModule() : audio_callback_(nullptr) {}

  int32_t RegisterAudioCallback(webrtc::AudioTransport* audio_callback) override;

  void Pull(uint16_t* audioSamples, size_t nSamples, uint32_t samplesPerSec);

 protected:
  virtual ~FauxAudioDeviceModule() {}

 private:
  webrtc::AudioTransport* audio_callback_;
  rtc::CriticalSection crit_;
};

#endif  // FAUXAUDIODEVICEMODULE_H_
