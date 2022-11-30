package tcpc_test

import (
	"fmt"
	"testing"
)

func TestNewReceiver(t *testing.T) {

	sender, err := NewSender[int](":3000")

	if err != nil {
		t.Error(err)
	}

	


	receiver, err := NewReceiver[int](":3000")
	if err != nil {
		t.Error(err)
	}

	sender.Chan <- 100

	msg := <-receiver.Chan

	fmt.Println(msg)
}
