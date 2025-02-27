package main

import (
	"embed"
	"flag"
	"fmt"
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
	})

	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%d", len(store.Clients()))
	})

	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", *port))

	_, err := fmt.Fprintf(os.Stderr, "Starting on http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)

	if err != nil {
		log.Fatal(err)
	}

	err = http.Serve(listener, nil)
	log.Fatal(err)
}
