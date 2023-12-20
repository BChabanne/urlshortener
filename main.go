package main

import (
	_ "embed"
	"flag"
	"log"
	"net/http"
)

//go:embed home.html
var homeHtml string

var addr *string

func init() {
	addr = flag.String("addr", "127.0.0.1:8000", "listen and serve")
	flag.Parse()
}

func home(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(homeHtml))
}

func router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		home(w, r)
		break
	default:
		http.Error(w, "Method Not alowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	log.Println("Server listening on", *addr)

	err := http.ListenAndServe(*addr, http.HandlerFunc(router))
	log.Fatal(err)
}
