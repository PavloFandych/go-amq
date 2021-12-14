package main

import (
	"github.com/go-stomp/stomp"
	"log"
	"sync"
)

func main() {
	conn, err := stomp.Dial("tcp", "localhost:61613")
	errorProcessing(err)

	var wg sync.WaitGroup
	wg.Add(2)

	go Consumer(conn, &wg)
	go Producer(conn, &wg)

	wg.Wait()
}

func errorProcessing(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Producer(stompConn *stomp.Conn, wg *sync.WaitGroup) {
	message := "TEXT MESSAGE"
	log.Println("Sent...->", message)
	if err := stompConn.Send("target-queue",
		"text/plain", []byte(message), stomp.SendOpt.Receipt,
		stomp.SendOpt.Header("JMSReplyTo", "response-queue")); err != nil {
		errorProcessing(err)
	}
	wg.Done()
}

func Consumer(stompConn *stomp.Conn, wg *sync.WaitGroup) {
	sub, err := stompConn.Subscribe("target-queue", stomp.AckAuto)
	errorProcessing(err)
	message := <-sub.C
	log.Println("Received...->", string(message.Body))
	wg.Done()
}
