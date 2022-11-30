package tcpc

import (
	"encoding/gob"
	"log"
	"net"
	"time"
)

type TCPC[T any] struct {
	listenAddr string
	remoteAddr string

	Sendchan     chan T
	Recvchan     chan T
	outboundConn net.Conn
	ln           net.Listener
}

func New[T any](listenAddr, remoteAddr string) (*TCPC[T], error) {
	tcpc := &TCPC[T]{
		listenAddr: listenAddr,
		remoteAddr: remoteAddr,
		Sendchan:   make(chan T, 10),
		Recvchan:   make(chan T, 10),
	}

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	tcpc.ln = ln
	go tcpc.loop()
	go tcpc.acceptLoop()
	go tcpc.dialRemoteAndRead()

	return tcpc, nil
}

func (t *TCPC[T]) loop() {
	for {
		msg := <-t.Sendchan
		log.Println("Sending msg over the wire: ", msg)
		if err := gob.NewEncoder(t.outboundConn).Encode(&msg); err != nil {
			log.Println(err)
		}
	}
}

func (t *TCPC[T]) acceptLoop() {

	defer func() {
		t.ln.Close()
	}()

	for {
		conn, err := t.ln.Accept()
		if err != nil {
			log.Println(" r.listener.Accept() Error: ", err)
			return
		}

		log.Printf("sender connected %s", conn.RemoteAddr())
		go t.handlerConn(conn)
	}
}

func (t *TCPC[T]) handlerConn(conn net.Conn) {
	for {
		var msg T

		if err := gob.NewDecoder(conn).Decode(&msg); err != nil {
			log.Println(err)
			continue
		}

		t.Recvchan <- msg
	}
}

func (t *TCPC[T]) dialRemoteAndRead() {

	conn, err := net.Dial("tcp", t.remoteAddr)
	if err != nil {
		log.Printf("dialRemoteAndRead (%s) net.Dial ", err)
		time.Sleep(time.Second * 3)
		t.dialRemoteAndRead()
	}

	// go t.readLoop(conn)

	t.outboundConn = conn
}

// func (t *TCPC[T]) readLoop(conn net.Conn) {
// 	var msg T
// 	for {
// 		if err := gob.NewDecoder(conn).Decode(&msg); err != nil {
// 			log.Println(err)
// 		}
// 	}
// }

var defaultDialInterval = 3 * time.Second

type Sender[T any] struct {
	Chan         chan T
	remoteAddr   string
	outboundConn net.Conn
	dialInterval time.Duration
}

func NewSender[T any](remoteAddr string) (*Sender[T], error) {

	sender := &Sender[T]{
		Chan:         make(chan T),
		remoteAddr:   remoteAddr,
		dialInterval: defaultDialInterval,
	}

	// go sender.dialRemoteAndRead()
	go sender.loop()

	return sender, nil
}

func (s *Sender[T]) loop() {
	for {
		msg := <-s.Chan

		log.Println("streaming msg over this wire", msg)

		if err := gob.NewEncoder(s.outboundConn).Encode(msg); err != nil {
			log.Println(err)
		}
	}
}

type Reciever[T any] struct {
	Chan       chan T
	listenAddr string
	listener   net.Listener
}

func NewReceiver[T any](listenAddr string) (*Reciever[T], error) {
	recv := &Reciever[T]{
		Chan:       make(chan T),
		listenAddr: listenAddr,
	}

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	recv.listener = ln
	// go recv.acceptLoop()

	return recv, nil
}
