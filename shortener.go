package main

type Shortener interface {
	Add(url string) (string, error)
	Get(slug string) (string, error)
}

type noop struct{}

var _ Shortener = &noop{}

func (*noop) Add(url string) (string, error) {
	return "noop-slug", nil
}

func (*noop) Get(slug string) (string, error) {
	return "noop-url", nil
}
