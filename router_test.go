package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mock struct{}

var _ Shortener = &mock{}

func (*mock) Add(url string) (string, error) {
	return "noop-slug", nil
}

func (*mock) Get(slug string) (string, error) {
	return "noop-url", nil
}

type mockError struct{}

var _ Shortener = &mockError{}

func (*mockError) Add(url string) (string, error) {
	return "", errors.New("shortener add url is not implemented")
}

func (*mockError) Get(slug string) (string, error) {
	return "", errors.New("shortener get slug is not implemented")
}

func TestRouterHome(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	handler := router(&mock{}, "https://tiny.io/")
	handler(resp, req)

	if resp.Result().StatusCode != http.StatusOK {
		t.Error("200 not returned")
	}
}

func TestRouterBadMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodHead, "/", nil)
	resp := httptest.NewRecorder()
	handler := router(&mock{}, "https://tiny.io/")
	handler(resp, req)

	if resp.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Error("HEAD is not a valid method")
	}
}

func TestRouterPostUrl(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	resp := httptest.NewRecorder()
	handler := router(&mock{}, "https://tiny.io/")
	handler(resp, req)
	if resp.Result().StatusCode != http.StatusOK {
		t.Error("mock should return ok")
	}

	req = httptest.NewRequest(http.MethodPost, "/", nil)
	resp = httptest.NewRecorder()
	handler = router(&mockError{}, "https://tiny.io/")
	handler(resp, req)
	if resp.Result().StatusCode != http.StatusInternalServerError {
		t.Error("mock should return error")
	}
}
