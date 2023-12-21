package main

import (
	"net/http/httptest"
	"strconv"
	"sync"
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

	clients := makeClient(publicURL, 100)

	bench(clients, 1, b)
}

func BenchmarkReadMemory(b *testing.B) {
	b.StopTimer()

	db := NewSqliteMemoryShortener()

	server := httptest.NewServer(nil)
	defer server.Close()

	publicURL := server.URL + "/"

	router := router(db, publicURL)
	server.Config.Handler = router

	clients := makeClient(publicURL, 100)

	bench(clients, 0, b)
}

func BenchmarkReadWriteMemory(b *testing.B) {
	b.StopTimer()

	db := NewSqliteMemoryShortener()

	server := httptest.NewServer(nil)
	defer server.Close()

	publicURL := server.URL + "/"

	router := router(db, publicURL)
	server.Config.Handler = router

	clients := makeClient(publicURL, 100)

	bench(clients, 0.1, b)
}

func makeClient(publicURL string, n int) []Shortener {
	clients := make([]Shortener, n)

	for i := 0; i < n; i++ {
		clients[i] = NewClient(publicURL)
	}

	return clients
}

func bench(shorteners []Shortener, writeRatio float64, b *testing.B) {
	slug, err := shorteners[0].Add("http://test/")
	if err != nil {
		b.Fatal(err)
	}

	add := make(chan int)
	get := make(chan struct{})

	m := &sync.RWMutex{}
	wg := &sync.WaitGroup{}
	wg.Add(len(shorteners))
	slugs := []string{slug}

	for i := 0; i < len(shorteners); i++ {
		go func(i int) {
			shortener := shorteners[i]
			for true {
				select {
				case i, ok := <-add:
					if !ok {
						wg.Done()
						return
					}
					slug, err := shortener.Add("http://test/" + strconv.Itoa(i))
					if err != nil {
						b.Fatal(err)
					}
					m.Lock()
					slugs = append(slugs, slug)
					m.Unlock()
					break
				case <-get:
					m.RLock()
					slug = slugs[rand.Intn(len(slugs))]
					m.RUnlock()
					_, err = shortener.Get(slug)
					if err != nil {
						b.Fatal(err)
					}
					break
				}
			}
		}(i)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		dice := rand.Float64()
		if dice < writeRatio {
			add <- i
		} else {
			get <- struct{}{}
		}
	}
	close(get)
	close(add)
	wg.Wait()
}
