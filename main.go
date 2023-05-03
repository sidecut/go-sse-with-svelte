package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var msgChan chan string

func timeHandler(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	if msgChan != nil {
		msg := time.Now().Format(time.TimeOnly)
		msgChan <- msg
	}

	return nil
}

func sseHandler(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	msgChan = make(chan string)
	defer func() {
		if msgChan != nil {
			close(msgChan)
			msgChan = nil
			fmt.Println("Client closed connection")
		}
	}()

	c.Response().Flush()

	for {
		select {
		case message := <-msgChan:
			c.String(http.StatusOK, fmt.Sprintf("data: %s\n\n", message))
			c.Response().Flush()

		case <-c.Request().Context().Done():
			fmt.Printf("Client %v disconnected from SSE.\n", c.Request().RemoteAddr)
			return nil
		}
	}
}

func main() {
	e := echo.New()
	e.GET("/event", sseHandler)
	e.GET("/time", timeHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
