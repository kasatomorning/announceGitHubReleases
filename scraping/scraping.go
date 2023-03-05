package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func scrape(url string) (releaseArray []string) {
	getRelease := colly.NewCollector()
	getRelease.OnHTML(".Link--primary", func(e *colly.HTMLElement) {
		href := e.Text
		fmt.Println("Release found:", href)
		releaseArray = append(releaseArray, href)
	})
	getRelease.OnRequest(func(r *colly.Request) {
		fmt.Println("Scraping URL:", r.URL.String())
	})
	getRelease.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})
	getRelease.Visit(url)
	return
}

func getNextPage(url string) (nextPageURL string) {
	getNextPage := colly.NewCollector()
	nextPageURL = ""
	getNextPage.OnHTML(".next_page", func(e *colly.HTMLElement) {
		nextPageLink := e.Attr("href")
		// reach the last page
		if nextPageLink != "" {
			nextPageURL = "https://github.com/" + nextPageLink
		}
	})
	getNextPage.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})
	getNextPage.Visit(url)
	if nextPageURL == "" {
		fmt.Println("Reached the last page")
	}
	return
}

func main() {
	url := "https://github.com/yt-dlp/yt-dlp/releases"
	pageURLs := []string{}
	allReleasesArray := []string{}
	for url != "" {
		pageURLs = append(pageURLs, url)
		url = getNextPage(url)
	}
	for i := 0; i < len(pageURLs); i++ {
		allReleasesArray = append(allReleasesArray, scrape(pageURLs[i])...)
	}
}
