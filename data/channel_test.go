package data

import (
	"fmt"
	"testing"
)

var c *Channel

func TestNewChannel(t *testing.T) {
	c = NewChannel(nil)
	if nil == c {
		t.Fatal("Could not create NewChannel")
	}
}

func TestChannelLabel(t *testing.T) {
	s := c.Label()
	fmt.Println(s)
}

