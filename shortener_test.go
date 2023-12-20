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
}
