package mydramalist

import (
	"github.com/gocolly/colly"
	"net/url"
	"strconv"
	"strings"
)

const (
	baseAddress = "mydramalist.com"
)

// Client is a client for working with MyDramaList; emulates a REST library
// interface
type Client struct {
	searcher *colly.Collector
	results  []*Drama
}

func NewClient() Client {
	client := Client{
		searcher: colly.NewCollector(colly.AllowedDomains(baseAddress)),
	}
	client.searcher.OnHTML(".text-primary.title", func(e *colly.HTMLElement) {
		titleUrl := "https://" + baseAddress + e.ChildAttr("a", "href")
		d, err := scrapeDrama(titleUrl)
		if err != nil {
			return
		}
		client.results = append(client.results, d)
	})
	client.searcher.OnHTML("ul.pagination", func(e *colly.HTMLElement) {
		activePageNumber, err := strconv.Atoi(e.ChildText(".active"))
		if err == nil {
			e.ForEachWithBreak(".page-item", func(i int, child *colly.HTMLElement) bool {
				pageNumber, err := strconv.Atoi(child.ChildAttr(".page-link", "href"))
				if err == nil && pageNumber == activePageNumber+1 {
					err = client.searcher.Visit(child.Attr("href"))
					if err != nil {
						return true
					}
					return false
				}
				return true
			})
		}
	})
	return client
}

func (c *Client) Search(query string) ([]*Drama, error) {
	defer func() {
		c.results = nil
	}()
	if err := c.searcher.Visit("https://mydramalist.com/search?q=" + url.QueryEscape(query)); err != nil {
		return nil, err
	}
	var dramas []*Drama
	for _, drama := range c.results {
		dramas = append(dramas, drama)
	}
	return dramas, nil
}

type dramaScraper struct {
	Scraper *colly.Collector
	Drama
}

func scrapeDrama(url string) (*Drama, error) {
	scraper := dramaScraper{
		Scraper: colly.NewCollector(),
		Drama:   Drama{},
	}
	scraper.Scraper.OnHTML(".box", func(e *colly.HTMLElement) {
		synopsis := e.ChildText(".show-synopsis")
		if len(synopsis) > 0 {
			scraper.Drama.Summary = synopsis
		}
		e.ForEach("li.list-item", func(i int, child *colly.HTMLElement) {
			if len(child.Text) > 0 {
				if strings.HasPrefix(child.Text, "Drama:") {
					scraper.Title = strings.TrimPrefix(child.Text, "Drama: ")
				} else if strings.HasPrefix(child.Text, "Country:") {
					scraper.Country = strings.TrimPrefix(child.Text, "Country: ")
				} else if strings.HasPrefix(child.Text, "Episodes:") {
					num, err := strconv.Atoi(strings.TrimPrefix(child.Text, "Episodes: "))
					if err == nil {
						scraper.NumEpisodes = uint(num)
					}
				} else if strings.HasPrefix(child.Text, "Native Title:") {
					scraper.NativeTitle = strings.TrimPrefix(child.Text, "Native Title: ")
				} else if strings.HasPrefix(child.Text, "Also Known As:") {
					split := strings.Split(strings.TrimPrefix(child.Text, "Also Known As: "), ",")
					for _, s := range split {
						scraper.OtherTitles = append(scraper.OtherTitles, strings.TrimSpace(s))
					}
				} else if strings.HasPrefix(child.Text, "Directors:") {
					split := strings.Split(strings.TrimPrefix(child.Text, "Directors: "), ",")
					for _, s := range split {
						scraper.Directors = append(scraper.Directors, strings.TrimSpace(s))
					}
				} else if strings.HasPrefix(child.Text, "Screenwriters:") {
					split := strings.Split(strings.TrimPrefix(child.Text, "Screenwriters: "), ",")
					for _, s := range split {
						scraper.Screenwriters = append(scraper.Screenwriters, strings.TrimSpace(s))
					}
				} else if strings.HasPrefix(child.Text, "Genres:") {
					split := strings.Split(strings.TrimPrefix(child.Text, "Genres: "), ",")
					for _, s := range split {
						scraper.Genres = append(scraper.Genres, strings.TrimSpace(s))
					}
				} else if strings.HasPrefix(child.Text, "Tags:") {
					split := strings.Split(strings.TrimPrefix(child.Text, "Tags: "), ",")
					for _, s := range split {
						scraper.Tags = append(scraper.Tags, strings.TrimSpace(s))
					}
				} else if strings.HasPrefix(child.Text, "Related Content") {
					child.ForEach(".title", func(k int, c *colly.HTMLElement) {
						scraper.RelatedContent = append(scraper.RelatedContent, c.ChildText("a.text-primary"))
					})
				} else if strings.HasPrefix(child.Text, "Score:") {
					scoreText := strings.TrimPrefix(child.Text, "Score: ")
					if len(scoreText) == 0 {
						return
					}
					num, err := strconv.ParseFloat(strings.Split(scoreText, " ")[0], 32)
					if err == nil {
						scraper.Rating = float32(num)
					}
				}
			}

		})
	})
	if err := scraper.Scraper.Visit(url); err != nil {
		return nil, err
	}
	return &scraper.Drama, nil
}
