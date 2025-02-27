package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gocolly/colly/v2"
)

type ScrapedData struct {
	URL       string    `xml:"url" toml:"url"`
	Timestamp time.Time `xml:"timestamp" toml:"timestamp"`
	Links     []Link    `xml:"links" toml:"links"`
}

type Link struct {
	Text string `xml:"text" toml:"text"`
	URL  string `xml:"url" toml:"url"`
}

func main() {
	// Get URL from user
	fmt.Print("Enter the URL to scrape: ")
	reader := bufio.NewReader(os.Stdin)
	userURL, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading input:", err)
	}
	userURL = strings.TrimSpace(userURL)

	// Get format preference
	fmt.Print("Choose output format (xml/toml): ")
	format, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading input:", err)
	}
	format = strings.ToLower(strings.TrimSpace(format))

	if format != "xml" && format != "toml" {
		log.Fatal("Invalid format. Please choose 'xml' or 'toml'")
	}

	data := &ScrapedData{
		URL:       userURL,
		Timestamp: time.Now(),
		Links:     make([]Link, 0),
	}

	c := colly.NewCollector()

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		data.Links = append(data.Links, Link{
			Text: e.Text,
			URL:  link,
		})
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error while scraping: %v\n", err)
	})

	err = c.Visit(userURL)
	if err != nil {
		log.Fatal(err)
	}

	// Save the data
	filename := fmt.Sprintf("scrape_result.%s", format)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	if format == "xml" {
		encoder := xml.NewEncoder(file)
		encoder.Indent("", "  ")
		if err := encoder.Encode(data); err != nil {
			log.Fatal("Error encoding XML:", err)
		}
	} else {
		if err := toml.NewEncoder(file).Encode(data); err != nil {
			log.Fatal("Error encoding TOML:", err)
		}
	}

	fmt.Printf("Results saved to %s\n", filename)
}
