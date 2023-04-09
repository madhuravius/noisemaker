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
)

var upgrader = gorilla_ws.Upgrader{} // use default options

func (r *RunConfig) generateSocketDataToSend() {
	r.socketData = make([]byte, int64(FileSizeToSend*1024*1024))
	_, _ = rand.Read(r.socketData)
}

// trashbinServerSocketFunc - send data on server
func trashbinServerSocketFunc(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Err on conn upgrade:", err)
		return
	}
	defer func(c *gorilla_ws.Conn) {
		closeErr := c.Close()
		if closeErr != nil {
			log.Print("Err on close:", closeErr.Error())
		}
	}(c)
	for {
		mt, message, errMsgRead := c.ReadMessage()
		if errMsgRead != nil {
			log.Println("read:", errMsgRead)
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
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
		log.Fatal("Err on dial: ", err.Error())
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
	http.HandleFunc("/", trashbinServerSocketFunc)
	log.Printf("Starting web server on :%d\n", r.desiredPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", r.desiredPort), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
