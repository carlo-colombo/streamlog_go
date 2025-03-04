package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/carlo-colombo/streamlog_go/sse"
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
			case <-store.FilterChangeFor():
				// Send reset event
				fmt.Fprintf(w, "event: reset\ndata: reset\n\n")
				flusher.Flush()

				// Send current filtered logs
				for _, logItem := range store.List() {
					_ = logItem.Encode(encoder)
				}
				flusher.Flush()
			case line := <-store.LineFor(uid):
				_ = line.Encode(encoder)
				flusher.Flush()
			}
		}
	}
}

func FilterHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var request struct {
			Filter string `json:"filter"`
		}
		if err := json.Unmarshal(body, &request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		store.SetFilter(request.Filter)
		w.WriteHeader(http.StatusOK)
	}
}
