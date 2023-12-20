package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Shortener interface {
	Add(url string) (string, error)
	Get(slug string) (string, error)
}

var InvalidURL = errors.New("Invalid URL")
var InvalidSlug = errors.New("Invalid Slug")
var SlugNotFound = errors.New("Slug Not Found")

type SqliteShortener struct {
	db *sql.DB
}

var _ Shortener = &SqliteShortener{}

func NewSqliteShortener(name string) (*SqliteShortener, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	return &SqliteShortener{db: db}, nil
}

func NewSqliteMemoryShortener() *SqliteShortener {
	shortener, err := NewSqliteShortener(":memory:")
	if err != nil {
		panic(err)
	}
	return shortener
}

func (*SqliteShortener) Add(url string) (string, error) {
	return "", fmt.Errorf("%w : %s", InvalidURL, url)
}

func (*SqliteShortener) Get(slug string) (string, error) {
	return "", fmt.Errorf("%w : %s", SlugNotFound, slug)
}
