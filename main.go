package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	gistID := os.Getenv("ENTRIES_URLS_GIST_ID")
	if len(gistID) == 0 {
		log.Fatalf("ENTRIES_URLS_GIST_ID failed to load from .env")
	}

	tele, err := NewTelegramBot()
	if err != nil {
		log.Fatalf("telegram bot failed to load: %v", err)
	}

	sp := NewScraper()

	m, err := NewEntriesManager("entries.json")
	if err != nil {
		log.Fatalf("entries manager failed to init: %v", err)
	}

	for {
		gistContent, _ := fetchLatestGist(gistID)
		urls := strings.Split(gistContent, "\n")
		fmt.Println(urls)
		m.UpdateEntries(urls)
		fmt.Println(m.GetEntries())

		entries := m.GetEntries()
		for _, entry := range entries {
			rdr := sp.getTextContent(entry.URL)
			buf := new(bytes.Buffer)
			buf.ReadFrom(rdr)
			str := buf.String()
			prices := getFlatPricesFromSite(str)
			firstPrice := prices[0]
			tele.SendMessage(firstPrice)
		}
		time.Sleep(300 * time.Second)
	}
}

func getFlatPricesFromSite(textContent string) []string {
	const currNotation = "zł"

	currNotationStartIndexes := startIndexStrOccurs(textContent, currNotation)
	if len(currNotationStartIndexes) == 0 {
		log.Printf("No `%s` has been found on site which textContent[:100] is %s", currNotation, textContent[:100])
	}

	validatePrice := func(s string) (bool, string) {
		price := make([]rune, 0)
		for _, r := range s {
			switch {
			case unicode.IsSpace(r):
				continue
			case unicode.IsDigit(r):
				price = append(price, r)
			default:
				return false, ""
			}
		}
		return true, string(price)
	}

	textContentRunes := []rune(textContent)
	prices := make([]string, 0)
	for _, idx := range currNotationStartIndexes {
		i := 0
		p := ""
		for {
			i++
			if idx-i < 0 {
				break
			}

			tp := textContentRunes[idx-i : idx]
			bp, vp := validatePrice(string(tp))
			if !bp {
				break
			}

			p = vp
		}

		if len(p) == 0 {
			continue
		}

		prices = append(prices, p)
	}

	return prices
}
