package internal

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	gorilla_ws "github.com/gorilla/websocket"
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

		if err = websocket.Message.Send(ws, reply); err != nil {
			log.Println("Err: Unable to send messages: ", err.Error())
			break
		}
	}
}

// trashbinClientSocketFunc - send data from client
// taken from here: https://github.com/gorilla/websocket/blob/master/examples/echo/client.go
// this is now deprecated/archive-only so we probably shouldn't use this but it doesn't appear there are
// a lot of alternative packages yet:
// https://www.reddit.com/r/golang/comments/zh0w0p/gorilla_web_toolkit_is_now_in_archive_only_mode/
func (r *RunConfig) trashbinClientSocketFunc() {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("localhost:%d", r.desiredPort), Path: "/"}
	c, _, err := gorilla_ws.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func(c *gorilla_ws.Conn) {
		wsErr := c.Close()
		if wsErr != nil {
			log.Println("Err: unable to close ws connection: ", wsErr.Error())
		}
	}(c)

	done := make(chan struct{})
	interrupt := make(chan os.Signal, 1)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			writeMsgErr := c.WriteMessage(gorilla_ws.TextMessage, r.socketData)
			if writeMsgErr != nil {
				log.Println("Err: unable to send write msg:", writeMsgErr.Error())
				return
			}
		case <-interrupt:
			log.Println("Warn: interrupting message to shut down server")
			writeMsgErr := c.WriteMessage(
				gorilla_ws.CloseMessage,
				gorilla_ws.FormatCloseMessage(gorilla_ws.CloseNormalClosure, ""),
			)
			if writeMsgErr != nil {
				log.Println("Err: unable to close write:", writeMsgErr.Error())
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// divide file into equal chunks to ensure speed

// startHttpServer - start http server with websocket
func (r *RunConfig) startHttpServer() {
	http.Handle("/", websocket.Handler(trashbinServerSocketFunc))
	log.Printf("Starting web server on :%d\n", r.desiredPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", r.desiredPort), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
