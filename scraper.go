package main

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// most popular user agent on Jul 22 2024
const USER_AGENT = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36`

type Scraper struct {
	client    *http.Client
	userAgent string
}

func NewScraper() *Scraper {
	return &Scraper{
		client:    &http.Client{},
		userAgent: USER_AGENT,
	}
}

func (sp *Scraper) getTextContent(url string) io.Reader {
	req := sp.prepareGET(url)
	resp, err := sp.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	doc, _ := html.Parse(strings.NewReader(string(body)))

	var builder strings.Builder

	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode && n.Parent != nil {
			switch n.Parent.Data {
			case "script", "style", "noscript":
				return
			}

			text := strings.TrimSpace(n.Data)
			if text != "" {
				builder.WriteString(text + "\n")
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(doc)

	return strings.NewReader(builder.String())
}

func (sp *Scraper) prepareGET(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", USER_AGENT)
	return req
}
