package main

import (
	"bufio"
	"flag"
	"fmt"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		for {
			fmt.Fprintln(w, <-logs)
			flusher.Flush()
		}
	})

	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", *port))

	fmt.Fprintf(os.Stderr, "Starting on http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, nil))
}
