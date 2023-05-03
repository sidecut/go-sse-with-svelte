package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var msgChan chan string

func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if msgChan != nil {
		msg := time.Now().Format(time.TimeOnly)
		msgChan <- msg
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	msgChan = make(chan string)
	defer func() {
		if msgChan != nil {
			close(msgChan)
			msgChan = nil
			fmt.Println("Client closed connection")
		}
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	flusher.Flush()

	for {
		select {
		case message := <-msgChan:
			fmt.Fprintf(w, "data: %s\n\n", message)
			flusher.Flush()

		case <-r.Context().Done():
			fmt.Printf("Client %v disconnected from SSE.\n", r.RemoteAddr)
			return
		}
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/event", sseHandler)
	router.HandleFunc("/time", timeHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
