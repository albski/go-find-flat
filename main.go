package main

import (
	"bytes"
	"fmt"
	"log"
	"unicode"
)

func getFlatPricesFromSite(textContent string) []string {
	const currNotation = "zł"

	currNotationStartIndexes := startIndexStrOccurs(textContent, currNotation)
	if len(currNotationStartIndexes) == 0 {
		log.Printf("No `%s` has been found on site with textContent[:100] of %s", currNotation, textContent[:100])
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
	sp := NewScraper()
	rdr := sp.getTextContent("https://www.olx.pl/nieruchomosci/mieszkania/poznan/q-polanka/?search%5Bfilter_float_price%3Ato%5D=2400&search%5Border%5D=created_at%3Adesc")
	buf := new(bytes.Buffer)
	buf.ReadFrom(rdr)

	str := buf.String()

	prices := getFlatPricesFromSite(str)
	fmt.Println(prices)
}
