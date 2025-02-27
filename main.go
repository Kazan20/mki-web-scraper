package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Get URL from user
	fmt.Print("Enter the URL to scrape: ")
	reader := bufio.NewReader(os.Stdin)
	userURL, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading input:", err)
	}
	userURL = strings.TrimSpace(userURL)

	// Create a new collector
	c := colly.NewCollector()

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
	err = c.Visit(userURL)
	if err != nil {
		log.Fatal(err)
	}
}
