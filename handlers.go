package main

import (
	"fmt"
	"github.com/carlo-colombo/streamlog_go/sse"
	"math/rand"
	"net/http"
	"strconv"
)

func ClientsHandler(store Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", len(store.Clients()))
	}
}

func LogsHandler(store Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)

		w.Header().Set("Content-Type", "text/event-stream")

		data, _ := templates.ReadFile("templates/log.html")
		encoder := sse.NewEncoder(w, string(data))

		flusher.Flush()

		for _, logItem := range store.List() {
			_ = logItem.Encode(encoder)
		}
		flusher.Flush()

		uid := strconv.Itoa(rand.Int())

	Response:
		for {
			select {
			case <-r.Context().Done():
				store.Disconnect(uid)
				break Response
			case line := <-store.LineFor(uid):
				_ = line.Encode(encoder)

				flusher.Flush()
			}
		}
	}
}
