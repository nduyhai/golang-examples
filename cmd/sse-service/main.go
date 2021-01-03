package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/event", handleSSE)
	http.HandleFunc("/message", sendMessage)
	_ = http.ListenAndServe(":9090", nil)
	log.Println("Server listen on 9090")
}

var messageChan chan string

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan = make(chan string)
	defer func() {
		close(messageChan)
		messageChan = nil
	}()

	flusher := w.(http.Flusher)

	for {
		select {

		case message := <-messageChan:
			_, _ = fmt.Fprintf(w, "data: %s\n\n", message)
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

func sendMessage(writer http.ResponseWriter, request *http.Request) {
	if messageChan != nil {
		messageChan <- "Hello " + time.Now().String()
	}
}
