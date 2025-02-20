package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	logs := []string{}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			log := scanner.Text()
			logs = append(logs, log)
			fmt.Println(log)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, log := range logs {
			fmt.Fprintln(w, log)
		}
	})

	listener, _ := net.Listen("tcp", ":0")

	fmt.Fprintf(os.Stderr, "Starting on http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, nil))
}
