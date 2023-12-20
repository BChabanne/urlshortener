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
	addr := "127.0.0.1:8000"
	log.Println("Server listening on", addr)
	err := http.ListenAndServe(addr, http.HandlerFunc(hello))
	log.Fatal(err)
}
