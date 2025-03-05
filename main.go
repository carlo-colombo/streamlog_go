package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", "0", "port")
	flag.Parse()

	store := NewStore()

	go store.Scan(os.Stdin)

	fsys, _ := fs.Sub(static, "app/dist/app/browser")

	http.Handle("/", http.FileServer(http.FS(fsys)))
	http.HandleFunc("/clients", ClientsHandler(store))
	http.HandleFunc("/logs", LogsHandler(store))
	http.HandleFunc("/filter", FilterHandler(store))

	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", *port))

	_, err := fmt.Fprintf(os.Stderr, "Starting on http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)

	if err != nil {
		log.Fatal(err)
	}

	err = http.Serve(listener, nil)
	log.Fatal(err)
}
