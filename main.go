package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

	listener, _ := net.Listen("tcp", ":0")

	fmt.Fprintf(os.Stderr, "Starting on http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, nil))
}
