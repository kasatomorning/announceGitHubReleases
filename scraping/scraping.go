package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func scrape(url string) (nextPageURL string) {
	getRelease := colly.NewCollector()
	getRelease.OnHTML(".Link--primary", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		fmt.Println("Link found:", href)
	})
	getRelease.OnRequest(func(r *colly.Request) {
		fmt.Println("Scraping URL:", r.URL.String())
	})
	getRelease.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})
	getRelease.Visit(url)

	getNextPage := colly.NewCollector()
	nextPageURL = ""
	getNextPage.OnHTML(".next_page", func(e *colly.HTMLElement) {
		nextPageLink := e.Attr("href")
		// reach the last page
		if nextPageLink == "" {
			fmt.Println("No more pages")
		} else {
			nextPageURL = "https://github.com/" + nextPageLink
		}
	})
	getNextPage.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})
	getNextPage.Visit(url)
	return
}

func main() {
	url := "https://github.com/yt-dlp/yt-dlp/releases"
	for url != "" {
		fmt.Println(url)
		url = scrape(url)
	}
}
