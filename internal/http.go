package internal

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func (r *RunConfig) generateSocketDataToSend() {
	r.socketData = make([]byte, int64(FileSizeToSend*1024*1024))
	_, _ = rand.Read(r.socketData)
}

// trashbinServerSocketFunc - send data on server
func trashbinServerSocketFunc(ws *websocket.Conn) {
	var err error
	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			log.Println("Err: Unable to receive messages: ", err.Error())
			break
		}

		msg := reply
		if err = websocket.Message.Send(ws, msg); err != nil {
			log.Println("Err: Unable to send messages: ", err.Error())
			break
		}
	}
}

// send data from client

// divide file into equal chunks to ensure speed

// startHttpServer - start http server with websocket
func (r *RunConfig) startHttpServer() {
	http.Handle("/", websocket.Handler(trashbinServerSocketFunc))
	log.Printf("Starting web server on :%d\n", r.desiredPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", r.desiredPort), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
