package main

import (
	_ "embed"
	"flag"
	"log"
	"net/http"
)

//go:embed home.html
var home string

var addr *string

func init() {
	addr = flag.String("addr", "127.0.0.1:8000", "listen and serve")
	flag.Parse()
}

func hello(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(home))
}

func main() {
	log.Println("Server listening on", *addr)
	err := http.ListenAndServe(*addr, http.HandlerFunc(hello))
	log.Fatal(err)
}
