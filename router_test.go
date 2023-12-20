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

type mockError struct {
	err error
}

var _ Shortener = &mockError{}

func (mock *mockError) Add(url string) (string, error) {
	if mock.err != nil {
		return "", mock.err
	}
	return "", errors.New("shortener add url is not implemented")
}

func (mock *mockError) Get(slug string) (string, error) {
	if mock.err != nil {
		return "", mock.err
	}
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

	req = httptest.NewRequest(http.MethodPost, "/", nil)
	resp = httptest.NewRecorder()
	handler = router(&mockError{
		err: InvalidURL,
	}, "https://tiny.io/")
	handler(resp, req)
	if resp.Result().StatusCode != http.StatusBadRequest {
		t.Error("bad request should be returned on invalid url")
	}
}

func TestGetSlug(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/slug", nil)
	resp := httptest.NewRecorder()
	handler := router(&mock{}, "https://tiny.io/")
	handler(resp, req)
	if resp.Result().StatusCode != http.StatusFound {
		t.Error("valid slug should get redirected")
	}

	req = httptest.NewRequest(http.MethodGet, "/slug", nil)
	resp = httptest.NewRecorder()
	handler = router(&mockError{}, "https://tiny.io/")
	handler(resp, req)
	if resp.Result().StatusCode != http.StatusInternalServerError {
		t.Error("internal error expected")
	}

	req = httptest.NewRequest(http.MethodGet, "/slug-not-found", nil)
	resp = httptest.NewRecorder()
	handler = router(&mockError{
		err: SlugNotFound,
	}, "https://tiny.io/")
	handler(resp, req)
	if resp.Result().StatusCode != http.StatusNotFound {
		t.Error("slug not found expected")
	}

	req = httptest.NewRequest(http.MethodGet, "/slug/nested", nil)
	resp = httptest.NewRecorder()
	handler = router(&mockError{
		err: InvalidSlug,
	}, "https://tiny.io/")
	handler(resp, req)
	if resp.Result().StatusCode != http.StatusBadRequest {
		t.Error("nested slug is invalid")
	}

	req = httptest.NewRequest(http.MethodGet, "/slug!", nil)
	resp = httptest.NewRecorder()
	handler = router(&mockError{
		err: InvalidSlug,
	}, "https://tiny.io/")
	handler(resp, req)
	if resp.Result().StatusCode != http.StatusBadRequest {
		t.Error("slug with invalid charater")
	}
}
