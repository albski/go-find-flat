package main

import (
	"bytes"
	"log"
	"unicode"

	"github.com/joho/godotenv"
)

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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	tele, err := NewTelegramBot()
	if err != nil {
		log.Fatalf("telegram bot failed to load: %v", err)
	}

	sp := NewScraper()
	rdr := sp.getTextContent("https://www.otodom.pl/pl/oferta/kawalerka-polanka-ul-katowicka-bezposrednio-ID4mDk3.html")
	buf := new(bytes.Buffer)
	buf.ReadFrom(rdr)

	str := buf.String()

	prices := getFlatPricesFromSite(str)
	firstPrice := string(prices[0])
	tele.SendMessage(firstPrice)
}
