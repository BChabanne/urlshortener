package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed home.html
var homeHtml string

//go:embed url-posted.html
var urlPostedHtml string

var urlPosted *template.Template

type UrlPostedData struct {
	URL string
}

func init() {
	var err error
	urlPosted, err = template.New("url-posted").Parse(urlPostedHtml)
	if err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(homeHtml))
}

func postURL(shortener Shortener, domain string, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("error when parsing form", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

	url := r.Form.Get("url")
	slug, err := shortener.Add(url)

	if err != nil {
		log.Println("error shortening url", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	urlPosted.Execute(w, UrlPostedData{
		URL: domain + slug,
	})
}

func router(shortener Shortener, domain string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			home(w, r)
			break
		case http.MethodPost:
			postURL(shortener, domain, w, r)
			break
		default:
			http.Error(w, "Method Not alowed", http.StatusMethodNotAllowed)
		}
	}
}
