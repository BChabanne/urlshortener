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
	write        *sql.DB
	read         *sql.DB
	writeChannel chan sqlInsert
}

var _ Shortener = &SqliteShortener{}

//go:embed schema.sql
var sqliteScript string

type sqlInsert struct {
	slug      string
	url       string
	asyncResp chan error
}

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

	writeChannel := make(chan sqlInsert)

	shortener := &SqliteShortener{
		read:         read,
		write:        write,
		writeChannel: writeChannel,
	}

	go shortener.batchInsert()

	return shortener, nil
}

func (shortener *SqliteShortener) Close() error {
	close(shortener.writeChannel)
	errClose := shortener.write.Close()
	errRead := shortener.read.Close()
	return errors.Join(errClose, errRead)
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

	err = shortener.writeSlug(slug, u)
	if err != nil {
		return "", err
	}

	return slug, nil
}

func (shortener *SqliteShortener) writeSlug(slug string, url string) error {
	asyncResp := make(chan error, 1)
	shortener.writeChannel <- sqlInsert{
		slug:      slug,
		url:       url,
		asyncResp: asyncResp,
	}
	return <-asyncResp
}

var maxBatchInsert = 1000

func (shortener *SqliteShortener) batchInsert() {
	inserts := make([]sqlInsert, 0, maxBatchInsert)
	for insert := range shortener.writeChannel {
		inserts = inserts[:0]
		inserts = append(inserts, insert)

		// batch already queued inserts
	Batch:
		for {
			select {
			case insert, ok := <-shortener.writeChannel:
				if !ok {
					break Batch
				}
				inserts = append(inserts, insert)
				if len(inserts) >= maxBatchInsert {
					break Batch
				}
			default:
				break Batch
			}
		}

		notifyErrs := func(err error) {
			for _, insert := range inserts {
				insert.asyncResp <- err
			}
		}

		tx, err := shortener.write.Begin()
		if err != nil {
			notifyErrs(err)
			continue
		}

		// TODO handle slug collision but with 60 bits of entropy its highly not probable
		// Though it happens on bench. Behaviour is to replace old url in that case
		// to avoir error
		stmt, err := tx.Prepare("INSERT OR REPLACE INTO urls(slug,url) VALUES (?,?)")
		if err != nil {
			notifyErrs(err)
			continue
		}

		for _, insert := range inserts {
			_, err := stmt.Exec(insert.slug, insert.url)
			if err != nil {
				insert.asyncResp <- err
			}
		}
		err = tx.Commit()
		notifyErrs(err)
	}
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
