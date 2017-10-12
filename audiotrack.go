package webrtc

/*
#include "mediastreamtrack.h"
#include "audiotrack.h"
*/
import "C"
import (
	"math"
	"reflect"
	"sync"
	"unsafe"
)

/*
An AudioTrack can be obtained via NewAudioTrack (a local track) or via
PeerConnection.OnAddTrack (a remote track).  Audio samples are provided or
retrieved, respectively, via an AudioSource or AudioSink.
*/
type AudioTrack struct {
	*mediaStreamTrack
	t C.CGO_AudioTrack
}

func newAudioTrack(t C.CGO_AudioTrack) *AudioTrack {
	return &AudioTrack{
		mediaStreamTrack: newMediaStreamTrack(C.CGO_MediaStreamTrack(t)),
		t:                t,
	}
}

// NewAudioTrack creates a local audio track.
func NewAudioTrack(label string, source AudioSource) *AudioTrack {
	t := C.CGO_NewAudioTrack(C.CString(label), C.CGO_GoAudioSource(audioSources.Set(source)))
	return newAudioTrack(t)
}

func (t *AudioTrack) AddSink(s AudioSink) {
	cAudioSinksMu.Lock()
	defer cAudioSinksMu.Unlock()

	if _, ok := cAudioSinks[s]; ok {
		panic("AudioSink already added")
	}
	csink := C.CGO_AudioTrack_AddSink(t.t, C.CGO_GoAudioSink(audioSinks.Set(s)))
	cAudioSinks[s] = csink
}

func (t *AudioTrack) RemoveSink(s AudioSink) {
	cAudioSinksMu.Lock()
	defer cAudioSinksMu.Unlock()

	if csink, ok := cAudioSinks[s]; ok {
		audioSinks.Delete(int(C.CGO_AudioTrack_RemoveSink(t.t, csink)))
		delete(cAudioSinks, s)
	}
}

var (
	audioSources = NewCGOMap()
	audioSinks   = NewCGOMap()

	cAudioSinksMu sync.Mutex
	cAudioSinks   = map[AudioSink]C.CGO_AudioSink{}
)

// An AudioSource pushes audio to the AudioSinks that are added to it.
type AudioSource interface {
	AddAudioSink(AudioSink)
	RemoveAudioSink(AudioSink)
}

// An AudioSink receives audio, typically from an AudioSource.
type AudioSink interface {
	/*
		OnAudioData is called when new audio data is available.
		len(data) == numberOfChannels; len(data[i]) is the same for all i.
	*/
	OnAudioData(data [][]float64, sampleRate float64)
}

//export cgoAudioSourceAddSink
func cgoAudioSourceAddSink(source C.CGO_GoAudioSource, sink C.CGO_AudioSink) {
	audioSources.Get(int(source)).(AudioSource).AddAudioSink(cAudioSink{sink})
}

//export cgoAudioSourceRemoveSink
func cgoAudioSourceRemoveSink(source C.CGO_GoAudioSource, sink C.CGO_AudioSink) {
	audioSources.Get(int(source)).(AudioSource).RemoveAudioSink(cAudioSink{sink})
}

//export cgoAudioSourceDestruct
func cgoAudioSourceDestruct(source C.CGO_GoAudioSource) {
	audioSources.Delete(int(source))
}

type cAudioSink struct {
	s C.CGO_AudioSink
}

func (s cAudioSink) OnAudioData(data [][]float64, sampleRate float64) {
	bitsPerSample := 16
	numberOfChannels := len(data)
	numberOfFrames := len(data[0])
	buf := make([]int16, numberOfChannels*numberOfFrames*bitsPerSample/2)
	i := 0
	for j := range data[0] {
		for _, ch := range data {
			buf[i] = int16(ch[j] * float64(math.MaxInt16))
			i++
		}
	}
	C.CGO_AudioSinkOnData(s.s, unsafe.Pointer(&buf[0]), C.int(bitsPerSample), C.int(sampleRate), C.int(numberOfChannels), C.int(numberOfFrames))
}

//export cgoAudioSinkOnData
func cgoAudioSinkOnData(s C.CGO_GoAudioSink, audioData unsafe.Pointer, bitsPerSample, sampleRate, numberOfChannels, numberOfFrames int) {
	if bitsPerSample != 16 {
		panic("expected 16 bits per sample")
	}

	var buf []int16
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Data = uintptr(audioData)
	sh.Len = numberOfChannels * numberOfFrames
	sh.Cap = sh.Len

	data := make([][]float64, numberOfChannels)
	for i := range data {
		data[i] = make([]float64, numberOfFrames)
	}

	i := 0
	for j := range data[0] {
		for _, ch := range data {
			ch[j] = float64(buf[i]) / float64(math.MaxInt16)
			i++
		}
	}

	audioSinks.Get(int(s)).(AudioSink).OnAudioData(data, float64(sampleRate))
}
