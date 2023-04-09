package internal

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/net/websocket"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func dataToSend() []byte {
	data := make([]byte, int64(FileSizeToSend*1024*1024))
	_, _ = rand.Read(data)
	return data
}

// send data on server
func trashbinServerSocketFunc(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer func(ws *websocket.Conn) {
			_ = ws.Close()
		}(ws)
		for {
			// Write
			err := websocket.Message.Send(ws, dataToSend())
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

// startHttpServer - start http server with websocket
func (r *RunConfig) startHttpServer() {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.GET("/", trashbinServerSocketFunc)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", r.desiredPort)))
}
