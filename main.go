package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/carlo-colombo/streamlog_go/logentry"
	"github.com/carlo-colombo/streamlog_go/sse"
	"io/fs"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
)

//go:embed templates
var templates embed.FS

func main() {
	port := flag.String("port", "0", "port")
	flag.Parse()

	store := NewStore()

	go store.Scan(os.Stdin)

	fsys, _ := fs.Sub(static, "app/dist/app/browser")

	http.Handle("/", http.FileServer(http.FS(fsys)))

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)

		encoder := newEncoderAndSetContentHeaders(w, r.URL.Query().Has("sse"))
		flusher.Flush()

		for _, logItem := range store.List() {
			logItem.Encode(encoder)
		}

		flusher.Flush()
		uid := strconv.Itoa(rand.Int())

	Response:
		for {
			select {
			case <-r.Context().Done():
				store.Unsubscribe(uid)
				break Response
			case line := <-store.LineFor(uid):
				line.Encode(encoder)

				flusher.Flush()
			}
		}
	})

	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", len(store.Clients()))
	})

	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", *port))

	_, err := fmt.Fprintf(os.Stderr, "Starting on http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)

	if err != nil {
		log.Fatal(err)
	}

	err = http.Serve(listener, nil)
	log.Fatal(err)
}

func newEncoderAndSetContentHeaders(w http.ResponseWriter, isSSE bool) logentry.Encoder {
	if isSSE {
		data, _ := templates.ReadFile("templates/log.html")
		w.Header().Set("Content-Type", "text/event-stream")
		return sse.NewEncoder(w, string(data))
	}

	w.Header().Set("Content-Type", "application/jsonl")
	return json.NewEncoder(w)
}
