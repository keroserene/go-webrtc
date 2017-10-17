#include <_cgo_export.h>  // Allow calling certain Go functions.

#include "audiotrack.h"

#include "webrtc/api/mediastreaminterface.h"
#include "webrtc/pc/mediastreamtrack.h"

using namespace webrtc;

class AudioTrack : public MediaStreamTrack<AudioTrackInterface> {
 public:
  explicit AudioTrack(std::string label, CGO_GoAudioSource source)
    : MediaStreamTrack<AudioTrackInterface>(label)
    , source(source) {}

  virtual std::string kind() const override {
    return kAudioKind;
  }

  virtual AudioSourceInterface* GetSource() const override {
    return nullptr;
  }

  virtual void AddSink(AudioTrackSinkInterface* sink) override {
    cgoAudioSourceAddSink(source, (CGO_AudioSink)sink);
  }
  virtual void RemoveSink(AudioTrackSinkInterface* sink) override {
    cgoAudioSourceRemoveSink(source, (CGO_AudioSink)sink);
  }

  virtual bool GetSignalLevel(int* level) override { return false; }

  virtual rtc::scoped_refptr<AudioProcessorInterface> GetAudioProcessor() override {
    return nullptr;
  }

protected:
  ~AudioTrack() {
    cgoAudioSourceDestruct(source);
  }

private:
  CGO_GoAudioSource source;
};

CGO_AudioTrack CGO_NewAudioTrack(const char* label, CGO_GoAudioSource source) {
  return new rtc::RefCountedObject<AudioTrack>(label, source);
}

class AudioSink : public AudioTrackSinkInterface {
public:
	explicit AudioSink(CGO_GoAudioSink s) : s(s) {}

	virtual void OnData(const void* audio_data,
		                int bits_per_sample,
		                int sample_rate,
		                size_t number_of_channels,
		                size_t number_of_frames) {
		cgoAudioSinkOnData(
			s,
			(void*)audio_data,
			bits_per_sample,
			sample_rate,
			(int)number_of_channels,
			(int)number_of_frames
		);
	}

	CGO_GoAudioSink s;
};

void CGO_AudioSinkOnData(CGO_AudioSink s, void* data, int bitsPerSample, int sampleRate, int numberOfChannels, int numberOfFrames) {
	((AudioTrackSinkInterface*)s)->OnData(data, bitsPerSample, sampleRate, (size_t)numberOfChannels, (size_t)numberOfFrames);
}

CGO_AudioSink CGO_AudioTrack_AddSink(CGO_AudioTrack t, CGO_GoAudioSink gs) {
  auto s = new AudioSink(gs);
  ((AudioTrackInterface*)t)->AddSink(s);
  return s;
}

CGO_GoAudioSink CGO_AudioTrack_RemoveSink(CGO_AudioTrack t, CGO_AudioSink cs) {
  auto s = (AudioSink*)cs;
  ((AudioTrackInterface*)t)->RemoveSink(s);
  auto gs = s->s;
  delete s;
  return gs;
}
