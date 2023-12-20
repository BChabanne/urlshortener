package main

import "net/http"

//go:embed home.html
var homeHtml string

func home(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(homeHtml))
}

func postURL(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("URL shortener is not implemented yet"))
}

func router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		home(w, r)
		break
	case http.MethodPost:
		postURL(w, r)
		break
	default:
		http.Error(w, "Method Not alowed", http.StatusMethodNotAllowed)
	}
}
