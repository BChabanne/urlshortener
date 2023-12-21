package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"
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

func postFormURL(shortener Shortener, domain string, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("error when parsing form", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

	url := r.Form.Get("url")
	slug, err := shortener.Add(url)

	if errors.Is(err, InvalidURL) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL"))
		return
	} else if err != nil {
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

func postJsonURL(shortener Shortener, domain string, w http.ResponseWriter, r *http.Request) {
	var body AddReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

	slug, err := shortener.Add(body.Url)

	if errors.Is(err, InvalidURL) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL"))
		return
	} else if err != nil {
		log.Println("error shortening url", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AddResp{
		Slug: slug,
	})
}

func getSlug(shortener Shortener, w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimLeft(r.URL.Path, "/")
	url, err := shortener.Get(slug)
	if errors.Is(err, SlugNotFound) {
		http.NotFound(w, r)
		return
	} else if errors.Is(err, InvalidSlug) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	} else if err != nil {
		log.Println("error while getting slug", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func router(shortener Shortener, domain string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path == "/" {
				home(w, r)
			} else {
				getSlug(shortener, w, r)
			}
			break
		case http.MethodPost:
			contentType := r.Header.Get("Content-Type")
			if contentType == "application/json" {
				postJsonURL(shortener, domain, w, r)
				return
			}
			postFormURL(shortener, domain, w, r)
			break
		default:
			http.Error(w, "Method Not alowed", http.StatusMethodNotAllowed)
		}
	}
}
