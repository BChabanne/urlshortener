package main

import "errors"

type Shortener interface {
	Add(url string) (string, error)
	Get(slug string) (string, error)
}

type notImplemented struct{}

var _ Shortener = &notImplemented{}

func (*notImplemented) Add(url string) (string, error) {
	return "", errors.New("shortener add url is not implemented")
}

func (*notImplemented) Get(slug string) (string, error) {
	return "", errors.New("shortener get slug is not implemented")
}
