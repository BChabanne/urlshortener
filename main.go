package main

import (
	_ "embed"
	"log"
	"net/http"
)

//go:embed home.html
var home string

func hello(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(home))
}

func main() {
	err := http.ListenAndServe("127.0.0.1:8000", http.HandlerFunc(hello))
	log.Fatal(err)
}
