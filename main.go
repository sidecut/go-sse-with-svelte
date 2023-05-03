package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var msgChan chan string

func timeHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	if msgChan != nil {
		msg := time.Now().Format(time.TimeOnly)
		msgChan <- msg
	}
}

func sseHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	msgChan = make(chan string)
	defer func() {
		if msgChan != nil {
			close(msgChan)
			msgChan = nil
			fmt.Println("Client closed connection")
		}
	}()

	c.Writer.Flush()

	for {
		select {
		case message := <-msgChan:
			c.String(http.StatusOK, "data: %s\n\n", message)
			c.Writer.Flush()

		case <-c.Done():
			fmt.Printf("Client %v disconnected from SSE.\n", c.Request.RemoteAddr)
			return
		}
	}
}

func main() {
	r := gin.Default()
	r.GET("/event", sseHandler)
	r.GET("/time", timeHandler)

	log.Fatal(r.Run(":8080"))
}
