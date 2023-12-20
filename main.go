package main

import (
	"flag"
	"log"
	"net/http"
)

var addr *string

func init() {
	addr = flag.String("addr", "127.0.0.1:8000", "listen and serve")
	flag.Parse()
}

func main() {
	log.Println("Server listening on", *addr)

	shortener := &noop{}

	err := http.ListenAndServe(*addr, router(shortener))
	log.Fatal(err)
}
