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

	"io/ioutil"
	"github.com/leslie-wang/go-webrtc"
	"github.com/leslie-wang/go-webrtc/demo/go2go/chat/common"
	"bytes"
	"time"
)

var username = "Alice"
var wait chan int
var c *common.Common


func startChat() {
	fmt.Println("------- chat enabled! -------")
}

func endChat() {
	fmt.Println("------- chat disabled -------")
}

func receiveChat(data []byte) {
	fmt.Printf("------- receive %s -------\n", string(data))
	fmt.Println("------- sendback pong -------")

	c.SendDC("pong")

	fmt.Println("------- chat quit -------")
	wait <- 1
}

func signalSend(msg string) {
	fmt.Println(msg + "\n")
	for {
		_, err := http.Post("http://localhost:7777/offer", "application/json", bytes.NewBufferString(msg))
		if err != nil {
			time.Sleep(time.Second * 5)
		} else {
			return
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read answer body error")
	}
	c.SignalReceive(string(buffer))
}

func startWebServer() {
	http.HandleFunc("/answer", handler)
	http.ListenAndServe(":6666", nil)
}

func main() {
	webrtc.SetLoggingVerbosity(1)

	wait = make(chan int, 1)
	fmt.Println("=== go-webrtc go2go chat demo ===")
	fmt.Println("Welcome, " + username + "!")

	c = &common.Common{}
	err := c.Start(true, signalSend, startChat, endChat, receiveChat)
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
