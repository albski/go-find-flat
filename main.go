package main

import (
	"io"
	"log"
	"os"
	"unicode"
)

func getFlatPricesFromSite(textContent string) {
	const currNotation = "zł"

	currNotationStartIndexes := startIndexStrOccurs(textContent, currNotation)
	if len(currNotationStartIndexes) == 0 {
		log.Printf("No `%s` has been found on site with textContent[:100] of %s", currNotation, textContent[:100])
	}

	prices := make([]string, 0)

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
		return len(price) > 0, string(price)
	}

}

func main() {
	sp := NewScraper()
	rdr := sp.getTextContent("https://www.olx.pl/nieruchomosci/mieszkania/poznan/q-polanka/?search%5Bfilter_float_price%3Ato%5D=2400&search%5Border%5D=created_at%3Adesc")
	io.Copy(os.Stdout, rdr)
}
