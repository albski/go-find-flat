package main

import (
	"bytes"
	"log"
	"os"
	"strconv"
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

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if botToken == "" || chatIDStr == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID is not set in environment variables")
	}
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		log.Fatalf("invalid TELEGRAM_CHAT_ID: %v", err)
	}
	tele := NewTelegramBot(botToken, chatID)

	sp := NewScraper()

	m, err := NewEntriesManager("entries.json")
	if err != nil {
		log.Fatalf("entries manager failed to init: %v", err)
	}

	for {
		gistContent, _ := fetchLatestGist(gistID)
		urls := strings.Split(gistContent, "\n")
		m.UpdateEntriesOccurs(urls)

		entries := m.GetEntries()
		for _, entry := range entries {
			rdr := sp.getTextContent(entry.URL)
			if rdr == nil {
				continue
			}
			buf := new(bytes.Buffer)
			buf.ReadFrom(rdr)
			str := buf.String()

			newPrices := getFlatPricesFromSite(str)
			same := compareOccurs(entry.Prices, newPrices)
			if same {
				continue
			}

			tele.SendMessage(`change detected: ` + entry.URL)
			err := entry.SetEntryPrices(newPrices, m.SetEntry)
			if err != nil {
				log.Print(err)
			}
		}
		time.Sleep(20 * time.Second)
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

func compareOccurs[T comparable](slice1 []T, slice2 []T) (same bool) {
	mapOccurs := func(slice []T) map[T]int {
		occurs := make(map[T]int)
		for _, item := range slice {
			occurs[item]++
		}
		return occurs
	}
	map1 := mapOccurs(slice1)
	map2 := mapOccurs(slice2)

	if len(map1) != len(map2) {
		return false
	}

	for key, count1 := range map1 {
		if count2, exists := map2[key]; !exists || count1 != count2 {
			return false
		}
	}

	return true
}
