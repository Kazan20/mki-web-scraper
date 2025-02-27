package main

import (
    "fmt"
    "log"
    "github.com/gocolly/colly/v2"
)

func main() {
    // Create a new collector
    c := colly.NewCollector(
        colly.AllowedDomains("example.com"), // Replace with the domain you want to scrape
    )

    // On every a element which has href attribute call callback
    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        fmt.Printf("Link found: %q -> %s\n", e.Text, link)
    })

    // Before making a request print "Visiting ..."
    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL.String())
    })

    // Print if error occurs
    c.OnError(func(r *colly.Response, err error) {
        log.Printf("Error while scraping: %v\n", err)
    })

    // Start scraping
    err := c.Visit("http://example.com") // Replace with the URL you want to scrape
    if err != nil {
        log.Fatal(err)
    }
}