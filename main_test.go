package main

import (
	"net/http/httptest"
	"strconv"
	"testing"
)

func BenchmarkInsertMemory(b *testing.B) {
	b.StopTimer()

	db := NewSqliteMemoryShortener()

	server := httptest.NewServer(nil)
	defer server.Close()

	publicURL := server.URL + "/"

	router := router(db, publicURL)
	server.Config.Handler = router

	client := NewClient(publicURL)

	bench(client, 1, b)
}

func BenchmarkReadMemory(b *testing.B) {
	b.StopTimer()

	db := NewSqliteMemoryShortener()

	server := httptest.NewServer(nil)
	defer server.Close()

	publicURL := server.URL + "/"

	router := router(db, publicURL)
	server.Config.Handler = router

	client := NewClient(publicURL)

	bench(client, 0, b)
}

func bench(shortener Shortener, writeRatio float64, b *testing.B) {
	slug, err := shortener.Add("http://test/")
	if err != nil {
		b.Fatal(err)
	}

	slugs := []string{slug}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		dice := rand.Float64()
		if dice >= writeRatio {
			slug, err := shortener.Add("http://test/" + strconv.Itoa(i))
			if err != nil {
				b.Fatal(err)
			}
			slugs = append(slugs, slug)
		} else {
			slug = slugs[rand.Intn(len(slugs))]
			_, err = shortener.Get(slug)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
