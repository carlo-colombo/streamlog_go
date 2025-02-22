package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/carlo-colombo/streamlog_go/log"
	"github.com/carlo-colombo/streamlog_go/sse"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
)

//go:embed templates
var templates embed.FS

//go:embed templates/log.html
var logTmpl embed.FS

func main() {
	port := flag.String("port", "0", "port")
	flag.Parse()

	logs := make(chan string)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			logs <- line
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		tmpl, _ := template.ParseFS(templates, "templates/index.html")

		var q struct{}
		tmpl.Execute(w, q)
	})

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)

		isSSE := r.URL.Query().Has("sse")
		if isSSE {
			w.Header().Set("Content-Type", "text/event-stream")
		} else {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(http.StatusOK)
		flusher.Flush()

		encoder := newEncoder(w, isSSE)

		for {
			log.Log{Line: <-logs}.Encode(encoder)

			flusher.Flush()
		}
	})

	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", *port))

	fmt.Fprintf(os.Stderr, "Starting on http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, nil))
}

func newEncoder(w io.Writer, isSSE bool) log.Encoder {
	if isSSE {
		data, _ := templates.ReadFile("templates/log.html")
		return sse.NewEncoder(w, string(data))
	}
	return json.NewEncoder(w)
}
