package main

import "errors"

type Shortener interface {
	Add(url string) (string, error)
	Get(slug string) (string, error)
}

var InvalidURL = errors.New("Invalid URL")
var InvalidSlug = errors.New("Invalid Slug")
var SlugNotFound = errors.New("Slug Not Found")

type noop struct{}

var _ Shortener = &noop{}

func (*noop) Add(url string) (string, error) {
	return "noop-slug", nil
}

func (*noop) Get(slug string) (string, error) {
	return "noop-url", nil
}
