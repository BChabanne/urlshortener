package main

import (
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func main() {
	err := http.ListenAndServe("127.0.0.1:8000", http.HandlerFunc(hello))
	log.Fatal(err)
}
