#ifndef _C_AUDIOTRACK_H
#define _C_AUDIOTRACK_H

#define WEBRTC_POSIX 1

#ifdef __cplusplus
extern "C" {
#endif

  // In order to present an interface cgo is happy with, nothing in this file
  // can directly reference header files from libwebrtc / C++ world. All the
  // casting must be hidden in the .cc file.

  typedef void* CGO_AudioTrack; // webrtc::AudioTrackInterface*
  typedef int CGO_GoAudioSource; // key into Go audioSourceMap
  typedef int CGO_GoAudioSink; // key into Go audioSinkMap
  typedef void* CGO_AudioSink; // webrtc::AudioTrackSinkInterface

  CGO_AudioTrack CGO_NewAudioTrack(const char* label, CGO_GoAudioSource source);

  void CGO_AudioSinkOnData(CGO_AudioSink s, void* data, int bitsPerSample, int sampleRate, int numberOfChannels, int numberOfFrames);

  CGO_AudioSink CGO_AudioTrack_AddSink(CGO_AudioTrack, CGO_GoAudioSink);
  CGO_GoAudioSink CGO_AudioTrack_RemoveSink(CGO_AudioTrack, CGO_AudioSink);

#ifdef __cplusplus
}
#endif

#endif  // _C_AUDIOTRACK_H
