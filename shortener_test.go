package main

import (
	"errors"
	"testing"
)

func TestShortener(t *testing.T) {
	shortener := NewSqliteMemoryShortener()

	_, err := shortener.Add("invalid url")
	if !errors.Is(err, InvalidURL) {
		t.Error("invalid url should be rejected")
	}

	_, err = shortener.Get("invalidslug!")
	if !errors.Is(err, InvalidSlug) {
		t.Error("invalid slug should be rejected")
	}

	_, err = shortener.Get("unknownslug")
	if !errors.Is(err, SlugNotFound) {
		t.Error("unknown slug returns error")
	}

	url := "https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74"
	slug, err := shortener.Add(url)
	if err != nil {
		t.Error("slug should have been generated")
	}

	actualUrl, err := shortener.Get(slug)
	if err != nil {
		t.Error("url should have been returned")
	}

	if url != actualUrl {
		t.Error("urls are not identical")
	}

	url2 := "https://other/url"
	slug2, err := shortener.Add(url2)
	if err != nil {
		t.Error("slug2 should have been generated")
	}

	if slug == slug2 {
		t.Error("slug and slug2 should not be same")
	}
}
