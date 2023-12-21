package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8000", "listen and serve")
	url := flag.String("url", "http://localhost:8000/", "url on which server is listening")
	db := flag.String("db", "file::memory:?cache=shared", "name of the sqlite database")
	flag.Parse()

	log.Println("Server listening on", *addr, "at", *url)
	log.Println("Sqlite database at", *db)

	shortener, err := NewSqliteShortener(*db)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(*addr, router(shortener, *url))
	log.Fatal(err)
}
