package main

import (
	"flag"
	"log"
	"net/http"
)

var addr *string
var url *string

func init() {
	addr = flag.String("addr", "127.0.0.1:8000", "listen and serve")
	url = flag.String("url", "http://localhost:8000/", "url on which server is listening")
	flag.Parse()
}

func main() {
	log.Println("Server listening on", *addr, "at", *url)

	shortener := &noop{}

	err := http.ListenAndServe(*addr, router(shortener, *url))
	log.Fatal(err)
}
