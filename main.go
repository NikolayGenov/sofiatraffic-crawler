package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScapeStations(stationCode string) {
	doc, err := goquery.NewDocument("http://m.sofiatraffic.bg/schedules/vehicle?stop=" + stationCode + "&lid=24&vt=0&rid=873")
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items

	doc.Find(".no-bullets li a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		station := strings.TrimSpace(s.Text())

		fmt.Printf("Station %d - %v\n", i, station)
	})
}

func main() {
	ScapeStations("1099")
}
