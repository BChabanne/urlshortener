package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Client struct {
	cli *http.Client
	url string
}

type AddReq struct {
	Url string `json:"url"`
}

type AddResp struct {
	Slug string `json:"slug"`
}

var _ Shortener = &Client{}

func NewClient(url string) *Client {
	cli := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &Client{
		cli,
		url,
	}
}

func (c *Client) Add(url string) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(AddReq{
		Url: url,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, c.url, buf)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.cli.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return "", fmt.Errorf("%w %s", InvalidURL, url)
	} else if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	var body AddResp
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", err
	}

	return body.Slug, nil
}

func (c *Client) Get(slug string) (string, error) {
	resp, err := c.cli.Get(c.url + slug)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return "", fmt.Errorf("%w %s", InvalidSlug, slug)
	} else if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("%w %s", SlugNotFound, slug)
	} else if resp.StatusCode != http.StatusFound {
		return "", errors.New(resp.Status)
	}

	url := resp.Header.Get("Location")

	return url, nil
}
