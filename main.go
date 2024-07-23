package main

import (
	"io"
	"os"
)

func main() {
	sp := NewScraper()
	rdr := sp.getTextContent("https://www.olx.pl/nieruchomosci/mieszkania/poznan/q-polanka/?search%5Bfilter_float_price%3Ato%5D=2400&search%5Border%5D=created_at%3Adesc")
	io.Copy(os.Stdout, rdr)
}
