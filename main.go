package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// most popular user agent on Jul 22 2024
const USER_AGENT = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36`

func getSiteTextContent(url string) io.Reader {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", USER_AGENT)

	client := &http.Client{}
	resp, _ := client.Do(req)
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
			text = strings.ReplaceAll(text, " ", "")
			if text != "" {
				builder.WriteString(text)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(doc)

	return strings.NewReader(builder.String())
}

func main() {
	rdr := getSiteTextContent("https://www.olx.pl/nieruchomosci/mieszkania/poznan/q-polanka/?search%5Bfilter_float_price%3Ato%5D=2400&search%5Border%5D=created_at%3Adesc")
	io.Copy(os.Stdout, rdr)
}
