package mydramalist

import (
	"github.com/gocolly/colly"
	"strconv"
)

const (
	baseAddress = "mydramalist.com"
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
		url := e.ChildAttr("a", "href")
		_ = client.scraper.Visit(url)
	})
	client.searcher.OnHTML("ul.pagination", func(e *colly.HTMLElement) {
		activePageNumber, err := strconv.Atoi(e.ChildText(".active"))
		if err == nil {
			e.ForEach(".page-item", func(i int, child *colly.HTMLElement) {
				pageNumber, err := strconv.Atoi(child.Text)
				if err == nil && pageNumber == activePageNumber+1 {
					_ = client.searcher.Visit(child.Attr("href"))
				}
			})
		}
	})
	return client
}
