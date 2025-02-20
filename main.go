package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
)

type Log struct {
	Line string `json:"line"`
}

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

		encoder := json.NewEncoder(w)

		for {
			encoder.Encode(Log{<-logs})
			flusher.Flush()
		}
	})

	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", *port))

	fmt.Fprintf(os.Stderr, "Starting on http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, nil))
}
