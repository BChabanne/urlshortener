package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterHome(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	handler := router(&noop{}, "https://tiny.io/")
	handler(resp, req)

	if resp.Result().StatusCode != http.StatusOK {
		t.Error("200 not returned")
	}
}
