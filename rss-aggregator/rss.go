package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) {
	rssFeed := RSSFeed{}

	client := http.Client{
		Timeout: time.Second * 10,
	}

	response, err := client.Get(url)

	if err != nil {
		return rssFeed, err
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)

	if err != nil {
		return rssFeed, err
	}

	err = xml.Unmarshal(bytes, &rssFeed)

	if err != nil {
		return rssFeed, err
	}

	return rssFeed, nil
}
