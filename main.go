package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	logs := make(chan string)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
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

	listener, _ := net.Listen("tcp", ":0")

	fmt.Fprintf(os.Stderr, "Starting on http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, nil))
}
