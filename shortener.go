package main

import (
	cryptorand "crypto/rand"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	mathrand "math/rand"
	"net/url"
	"regexp"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

var rand *mathrand.Rand
var slugRegexp *regexp.Regexp

func init() {
	var seed int64
	err := binary.Read(cryptorand.Reader, binary.BigEndian, &seed)
	if err != nil {
		panic(err)
	}
	source := mathrand.NewSource(seed)
	rand = mathrand.New(source)

	slugRegexp, err = regexp.Compile("^[a-zA-Z0-9]+$")
	if err != nil {
		panic(err)
	}
}

type Shortener interface {
	Add(url string) (string, error)
	Get(slug string) (string, error)
}

var InvalidURL = errors.New("Invalid URL")
var InvalidSlug = errors.New("Invalid Slug")
var SlugNotFound = errors.New("Slug Not Found")

type SqliteShortener struct {
	write *sql.DB
	read  *sql.DB
}

var _ Shortener = &SqliteShortener{}

//go:embed schema.sql
var sqliteScript string

func NewSqliteShortener(name string) (*SqliteShortener, error) {

	write, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	write.SetMaxIdleConns(1)
	write.SetMaxOpenConns(1)
	write.SetConnMaxLifetime(0)
	write.SetConnMaxIdleTime(0)

	_, err = write.Exec(sqliteScript)
	if err != nil {
		return nil, err
	}

	read, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	return &SqliteShortener{
		read:  read,
		write: write,
	}, nil
}

func NewSqliteMemoryShortener() *SqliteShortener {
	shortener, err := NewSqliteShortener("file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	shortener.read = shortener.write
	return shortener
}

const base62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// slugSize of 10  as an entropy of log2(62**10) = 59.5 bits
const slugSize = 10

func newSlug() string {
	buffer := make([]byte, slugSize)
	for i := 0; i < slugSize; i++ {
		buffer[i] = base62[rand.Intn(len(base62))]
	}
	return string(buffer)
}

func (shortener *SqliteShortener) Add(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("%w : %s", InvalidURL, u)
	}

	slug := newSlug()

	// TODO handle slug collision but with 60 bits of entropy its highly not probable
	_, err = shortener.write.Exec("INSERT INTO urls(slug,url) VALUES (?,?)", slug, u)
	if err != nil {
		return "", err
	}

	return slug, nil
}

func (shortener *SqliteShortener) Get(slug string) (string, error) {
	if !slugRegexp.MatchString(slug) {
		return "", fmt.Errorf("%w : %s", InvalidSlug, slug)
	}

	rows, err := shortener.read.Query("SELECT url FROM urls WHERE slug=?", slug)
	if err != nil {
		return "", err
	}
	if !rows.Next() {
		return "", fmt.Errorf("%w : %s", SlugNotFound, slug)
	}

	var url string
	err = rows.Scan(&url)
	if err != nil {
		return "", err
	}
	err = rows.Close()
	return url, err
}
