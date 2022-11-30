package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kitavrus/tcpc/tcpc"
)

func main() {

	channelLocal, err := tcpc.New[string](":3000", ":4000")
	if err != nil {
		log.Fatal(err)
	}

	channelLocal.Sendchan <- "GO Sendchan"

	channelRemote, err := tcpc.New[string](":4000", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	msg := <-channelRemote.Recvchan

	fmt.Println("Receive message from channel: ", msg)

	// channelRemote.Recvchan <- "GO"

	// sender, err := tcpc.NewSender[int](":3000")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// receiver, err := tcpc.NewReceiver[int](":3000")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// time.Sleep(2 * time.Second)

	// sender.Chan <- 100

	// msg := <-receiver.Chan

	// fmt.Println("received from channel over the wire: ", msg)
}

/*
sender, err := tcpc.NewSender[int](":3000")
receiver, err := tcpc.NewReceiver[int](":3000")

что будет, если у получателя и отправителя будет разние типы переданы в канал?

*/
