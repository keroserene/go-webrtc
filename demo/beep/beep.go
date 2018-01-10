/*
WebRTC audio bee demo server.
See beep.html for the client.

To use, `go run beep.go`, then open beep.html in a browser.
This server will send a tone to the browser, which will play it.
*/
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/keroserene/go-webrtc"
	"golang.org/x/net/websocket"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	webrtc.SetLoggingVerbosity(1)

	http.Handle("/ws", websocket.Handler(handleWebSocket))
	log.Fatal(http.ListenAndServe(":49372", nil))
}

func handleWebSocket(ws *websocket.Conn) {
	log.Printf("receive connection\n")
	pc, err := webrtc.NewPeerConnection(webrtc.NewConfiguration())
	if err != nil {
		log.Fatal(err)
	}

	beep := &beep{}
	pc.AddTrack(webrtc.NewAudioTrack("beep-audio", beep), nil)
	go beep.run()

	pc.OnIceCandidate = func(c webrtc.IceCandidate) {
		if err := websocket.JSON.Send(ws, c); err != nil {
			log.Println(err)
			return
		}
	}

	for {
		var msg struct {
			Type string
			Body json.RawMessage
		}
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}

		log.Printf("receive message: %v\n", msg)

		switch msg.Type {
		case "offer":
			offer := webrtc.DeserializeSessionDescription(string(msg.Body))
			if err := pc.SetRemoteDescription(offer); err != nil {
				log.Println(err)
				return
			}
			answer, err := pc.CreateAnswer()
			if err != nil {
				log.Println(err)
				return
			}
			if err := pc.SetLocalDescription(answer); err != nil {
				log.Println(err)
				return
			}
			if err := websocket.JSON.Send(ws, answer); err != nil {
				log.Println(err)
				return
			}
		case "icecandidate":
			c := webrtc.DeserializeIceCandidate(string(msg.Body))
			if err := pc.AddIceCandidate(*c); err != nil {
				log.Println(err)
			}
		default:
			log.Println("unexpected message type:", msg.Type)
		}
	}
}

type beep struct {
	sync.Mutex
	sinks []webrtc.AudioSink
}

func (b *beep) AddAudioSink(s webrtc.AudioSink) {
	b.Lock()
	b.sinks = append(b.sinks, s)
	b.Unlock()
}

func (b *beep) RemoveAudioSink(s webrtc.AudioSink) {
	b.Lock()
	defer b.Unlock()
	for i, s2 := range b.sinks {
		if s2 == s {
			b.sinks = append(b.sinks[:i], b.sinks[i+1:]...)
		}
	}
}

func (b *beep) run() {
	const (
		sampleRate     = 48000
		chunkRate      = 100
		numberOfFrames = sampleRate / chunkRate
		toneFrequency  = 256
	)
	data := [][]float64{make([]float64, numberOfFrames)}
	count := 0
	x := 0.04
	for next := time.Now(); ; next = next.Add(time.Second / chunkRate) {
		time.Sleep(next.Sub(time.Now()))

		for i := range data[0] {
			if count%(sampleRate/toneFrequency/2) == 0 {
				x = -x
			}
			data[0][i] = x
			count++
		}

		b.Lock()
		for _, sink := range b.sinks {
			sink.OnAudioData(data, sampleRate)
		}
		b.Unlock()
	}
}
