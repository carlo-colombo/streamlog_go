package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/carlo-colombo/streamlog_go/log"
	"github.com/carlo-colombo/streamlog_go/sse"
	"io"
	"net"
	"net/http"
	"os"
)

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

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)
		w.WriteHeader(http.StatusOK)
		flusher.Flush()

		encoder := newEncoder(w, r.URL.Query().Has("sse"))

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
		return sse.NewEncoder(w)
	}
	return json.NewEncoder(w)
}
