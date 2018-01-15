/*
WebRTC audio echo demo server.
See echo.html for the client.

To use, `go run echo.go`, then open echo.html in a browser.
You may want to wear headphones or turn your volume low,
because there will be some audio feedback (much of which should
be removed by echo cancellation).
*/
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/leslie-wang/go-webrtc"
	"golang.org/x/net/websocket"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	webrtc.SetLoggingVerbosity(1)

	http.Handle("/ws", websocket.Handler(handleWebSocket))
	log.Fatal(http.ListenAndServe(":49372", nil))
}

func handleWebSocket(ws *websocket.Conn) {
	pc, err := webrtc.NewPeerConnection(webrtc.NewConfiguration())
	if err != nil {
		log.Fatal(err)
	}

	pc.OnIceCandidate = func(c webrtc.IceCandidate) {
		if err := websocket.JSON.Send(ws, c); err != nil {
			log.Println(err)
			return
		}
	}

	pc.OnAddTrack = func(r *webrtc.RtpReceiver, s []*webrtc.MediaStream) {
		echo := &echo{}
		r.Track().(*webrtc.AudioTrack).AddSink(echo)
		pc.AddTrack(webrtc.NewAudioTrack("audio-echo", echo), nil)

		// A much simpler way to echo audio (but less useful
		// for demonstrative purposes) is to just:
		// pc.AddTrack(r.Track(), s)
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

type echo struct {
	sync.Mutex
	sinks []webrtc.AudioSink
}

func (e *echo) AddAudioSink(s webrtc.AudioSink) {
	e.Lock()
	defer e.Unlock()
	e.sinks = append(e.sinks, s)
}

func (e *echo) RemoveAudioSink(s webrtc.AudioSink) {
	e.Lock()
	defer e.Unlock()
	for i, s2 := range e.sinks {
		if s2 == s {
			e.sinks = append(e.sinks[:i], e.sinks[i+1:]...)
		}
	}
}

func (e *echo) OnAudioData(data [][]float64, sampleRate float64) {
	e.Lock()
	defer e.Unlock()
	for _, s := range e.sinks {
		s.OnAudioData(data, sampleRate)
	}
}
