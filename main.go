package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/carlo-colombo/streamlog_go/logentry"
	"github.com/carlo-colombo/streamlog_go/sse"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
)

//go:embed templates
var templates embed.FS

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

		encoder := newEncoderAndSetContentHeaders(w, r.URL.Query().Has("sse"))
		flusher.Flush()

		for {
			logentry.Log{Line: <-logs}.Encode(encoder)

			flusher.Flush()
		}
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
