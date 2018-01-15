/*
 * Webrtc chat demo.
 * Send chat messages via webrtc, over go.
 * Can interop with the JS client. (Open chat.html in a browser)
 *
 * To use: `go run chat.go`
 */
package main

import (
	"fmt"
	"os"
	"os/signal"
	"net/http"

	"github.com/leslie-wang/go-webrtc"
	"io/ioutil"
	"bytes"
	"time"
	"github.com/leslie-wang/go-webrtc/demo/go2go/chat/common"
)

var username = "Bob"
var wait chan int
var c *common.Common

func signalSend(msg string) {
	fmt.Println(msg + "\n")
	for {
		_, err := http.Post("http://localhost:6666/answer", "application/json", bytes.NewBufferString(msg))
		if err != nil {
			time.Sleep(time.Second * 5)
		} else {
			return 
		}
	}
}

func startChat() {
	fmt.Println("------- chat enabled! -------")
	fmt.Println("------- start sending ping -------")
	c.SendDC("ping")
}

func endChat() {
	fmt.Println("------- chat disabled -------")
}

func receiveChat(data []byte) {
	fmt.Printf("------- receive %s -------\n", string(data))

	fmt.Println("------- chat quit -------")

	wait <-1
}

func handler(w http.ResponseWriter, r *http.Request) {
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read answer body error")
	}
	c.SignalReceive(string(buffer))
}

func startWebServer() {
	http.HandleFunc("/offer", handler)
	http.ListenAndServe(":7777", nil)
}

func main() {
	webrtc.SetLoggingVerbosity(1)

	wait = make(chan int, 1)
	fmt.Println("=== go-webrtc go2go chat demo ===")
	fmt.Println("Welcome, " + username + "!")

	// start as the "answerer."
	c = &common.Common{}
	err := c.Start(false, signalSend, startChat, endChat, receiveChat)
	if err != nil {
		fmt.Print(err)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Println("Demo interrupted. Disconnecting...")
		c.Close()
		os.Exit(1)
	}()

	go startWebServer()

	<-wait
	fmt.Println("done")
}
