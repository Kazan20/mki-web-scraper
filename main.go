package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

type model struct {
	urlInput    textinput.Model
	formatInput textinput.Model
	step        int
	data        *ScrapedData
	done        bool
	err         error
}

func initialModel() model {
	urlInput := textinput.New()
	urlInput.Placeholder = "Enter URL to scrape"
	urlInput.Focus()

	formatInput := textinput.New()
	formatInput.Placeholder = "xml or toml"

	return model{
		urlInput:    urlInput,
		formatInput: formatInput,
		step:        0,
		data:        nil,
		done:        false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.step == 0 {
				m.step++
				m.urlInput.Blur()
				m.formatInput.Focus()
				return m, nil
			}
			if m.step == 1 {
				format := strings.ToLower(strings.TrimSpace(m.formatInput.Value()))
				if format != "xml" && format != "toml" {
					m.err = fmt.Errorf("invalid format. Please choose 'xml' or 'toml'")
					return m, nil
				}
				m.formatInput.Blur()
				go m.scrapeAndSave()
				m.step++
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	if m.step == 0 {
		m.urlInput, cmd = m.urlInput.Update(msg)
	} else if m.step == 1 {
		m.formatInput, cmd = m.formatInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress Ctrl+C to quit\n", m.err)
	}

	if m.done {
		return "Scraping completed! Results saved to file.\nPress Ctrl+C to quit\n"
	}

	var s string
	switch m.step {
	case 0:
		s = fmt.Sprintf("Enter URL to scrape:\n%s\n", m.urlInput.View())
	case 1:
		s = fmt.Sprintf("URL: %s\nChoose format (xml/toml):\n%s\n", m.urlInput.Value(), m.formatInput.View())
	case 2:
		s = "Scraping in progress...\n"
	}

	s += "\nPress Ctrl+C to quit\n"
	return s
}

func (m *model) scrapeAndSave() {
	url := strings.TrimSpace(m.urlInput.Value())
	format := strings.ToLower(strings.TrimSpace(m.formatInput.Value()))

	m.data = &ScrapedData{
		URL:   url,
		Links: make([]Link, 0),
	}

	c := colly.NewCollector()

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		m.data.Links = append(m.data.Links, Link{
			Text: e.Text,
			URL:  link,
		})
	})

	if err := c.Visit(url); err != nil {
		m.err = err
		return
	}

	filename := fmt.Sprintf("scrape_result.%s", format)
	file, err := os.Create(filename)
	if err != nil {
		m.err = err
		return
	}
	defer file.Close()

	if format == "xml" {
		encoder := xml.NewEncoder(file)
		encoder.Indent("", "  ")
		if err := encoder.Encode(m.data); err != nil {
			m.err = err
			return
		}
	} else {
		if err := toml.NewEncoder(file).Encode(m.data); err != nil {
			m.err = err
			return
		}
	}

	m.done = true
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
