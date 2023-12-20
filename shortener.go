package main

import "errors"

type Shortener interface {
	Add(url string) (string, error)
	Get(slug string) (string, error)
}

type noop struct{}

var _ Shortener = &noopError{}

func (*noop) Add(url string) (string, error) {
	return "noop-slug", nil
}

func (*noop) Get(slug string) (string, error) {
	return "noop-url", nil
}

type noopError struct{}

var _ Shortener = &noopError{}

func (*noopError) Add(url string) (string, error) {
	return "", errors.New("shortener add url is not implemented")
}

func (*noopError) Get(slug string) (string, error) {
	return "", errors.New("shortener get slug is not implemented")
}
