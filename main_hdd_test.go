//go:build hdd
// +build hdd

package main

import (
	"net/http/httptest"
	"os"
	"testing"
)

func BenchmarkInsertHdd(b *testing.B) {
	b.StopTimer()

	file := b.Name() + ".sqlite"
	os.Remove(file)
	db, err := NewSqliteShortener(file)
	if err != nil {
		b.Fatal(err)
	}

	server := httptest.NewServer(nil)
	defer server.Close()

	publicURL := server.URL + "/"

	router := router(db, publicURL)
	server.Config.Handler = router

	clients := makeClient(publicURL, 100)

	bench(clients, 1, b)
}

func BenchmarkReadHdd(b *testing.B) {
	b.StopTimer()

	file := b.Name() + ".sqlite"
	os.Remove(file)
	db, err := NewSqliteShortener(file)
	if err != nil {
		b.Fatal(err)
	}

	server := httptest.NewServer(nil)
	defer server.Close()

	publicURL := server.URL + "/"

	router := router(db, publicURL)
	server.Config.Handler = router

	clients := makeClient(publicURL, 100)

	bench(clients, 0, b)
}

func BenchmarkReadWriteHdd(b *testing.B) {
	b.StopTimer()

	file := b.Name() + ".sqlite"
	os.Remove(file)
	db, err := NewSqliteShortener(file)
	if err != nil {
		b.Fatal(err)
	}

	server := httptest.NewServer(nil)
	defer server.Close()

	publicURL := server.URL + "/"

	router := router(db, publicURL)
	server.Config.Handler = router

	clients := makeClient(publicURL, 100)

	bench(clients, 0.1, b)
}
