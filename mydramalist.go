package mydramalist

import (
	"github.com/gocolly/colly"
)

const (
	baseAddress = "https://mydramalist.com/"
)

// Client is a client for working with MyDramaList; emulates a REST library
// interface
type Client struct {
	searcher      *colly.Collector
	scraper       *colly.Collector
	searchResults []Drama
}

func NewClient() Client {
	client := Client{
		searcher: colly.NewCollector(colly.AllowedDomains(baseAddress)),
		scraper:  colly.NewCollector(colly.AllowedDomains(baseAddress)),
	}
	client.searcher.OnHTML(".text-primary.title", func(e *colly.HTMLElement) {
		title := e.ChildText("a")

	})
	return client
}
