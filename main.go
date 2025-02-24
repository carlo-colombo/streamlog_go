package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/carlo-colombo/streamlog_go/logentry"
	"github.com/carlo-colombo/streamlog_go/sse"
	"io/fs"
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

	var clients []chan string
	logs := make(chan string)
	var logsDb []string

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()

			logs <- line
		}
	}()

	go func() {
		for {
			line := <-logs
			logsDb = append(logsDb, line)
			for _, client := range clients {
				client <- line
			}
		}
	}()

	fsys, _ := fs.Sub(static, "app/dist/app/browser")
	http.Handle("/", http.FileServer(http.FS(fsys)))

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)

		encoder := newEncoderAndSetContentHeaders(w, r.URL.Query().Has("sse"))
		flusher.Flush()

		client := make(chan string)
		clients = append(clients, client)

		fmt.Println("client count:", len(clients))
		fmt.Println("current log count:", len(logsDb))

		for _, log := range logsDb {
			logentry.Log{Line: log}.Encode(encoder)
		}

		flusher.Flush()

		for {
			logentry.Log{Line: <-client}.Encode(encoder)

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
