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
	message := "file:///home/total/forTest/1/"
	log.Println("Sent...->", message)
	if err := stompConn.Send("mi.ReSyncDataProviderService.v.3.1.listDeepSubFoldersByPathByType",
		"text/plain", []byte(message), stomp.SendOpt.Receipt,
		stomp.SendOpt.Header("JMSReplyTo", "RESPONSES")); err != nil {
		errorProcessing(err)
	}
	wg.Done()
}

func Consumer(stompConn *stomp.Conn, wg *sync.WaitGroup) {
	sub, err := stompConn.Subscribe("RESPONSES", stomp.AckAuto)
	errorProcessing(err)
	message := <-sub.C
	log.Println("Received...->", string(message.Body))
	wg.Done()
}
