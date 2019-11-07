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
}

func NewClient() Client {
	client := Client{
		searcher: colly.NewCollector(colly.AllowedDomains(baseAddress)),
	}
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

func (c *Client) Search(query string) ([]Drama, error) {
	var dramas []Drama
	c.searcher.OnHTML(".text-primary.title", func(e *colly.HTMLElement) {
		titleUrl := "https://" + baseAddress + e.ChildAttr("a", "href")
		drama, err := scrapeDrama(titleUrl)
		if err != nil {
			return
		}
		dramas = append(dramas, drama)
	})
	err := c.searcher.Visit("https://mydramalist.com/search?q=" + url.QueryEscape(query))
	if err != nil {
		return nil, err
	}
	return dramas, nil
}

func scrapeDrama(url string) (Drama, error) {
	drama := Drama{}
	scraper := colly.NewCollector()
	scraper.OnHTML(".box", func(e *colly.HTMLElement) {
		synopsis := e.ChildText(".show-synopsis")
		if len(synopsis) > 0 {
			drama.Summary = synopsis
		}
		e.ForEach("li.list-item", func(i int, child *colly.HTMLElement) {
			if len(child.Text) > 0 {
				if strings.HasPrefix(child.Text, "Drama:") {
					drama.Title = strings.TrimPrefix(child.Text, "Drama: ")
				} else if strings.HasPrefix(child.Text, "Country:") {
					drama.Country = strings.TrimSpace(strings.TrimPrefix(child.Text, "Country: "))
				} else if strings.HasPrefix(child.Text, "Episodes:") {
					num, err := strconv.Atoi(strings.TrimPrefix(child.Text, "Episodes: "))
					if err == nil {
						drama.NumEpisodes = uint(num)
					}
				} else if strings.HasPrefix(child.Text, "Native Title:") {
					drama.NativeTitle = strings.TrimPrefix(child.Text, "Native Title: ")
				} else if strings.HasPrefix(child.Text, "Also Known As:") {
					split := strings.Split(strings.TrimPrefix(child.Text, "Also Known As: "), ",")
					for _, s := range split {
						drama.OtherTitles = append(drama.OtherTitles, strings.TrimSpace(s))
					}
				} else if strings.HasPrefix(child.Text, "Directors:") {
					split := strings.Split(strings.TrimPrefix(child.Text, "Directors: "), ",")
					for _, s := range split {
						drama.Directors = append(drama.Directors, strings.TrimSpace(s))
					}
				} else if strings.HasPrefix(child.Text, "Screenwriters:") {
					split := strings.Split(strings.TrimPrefix(child.Text, "Screenwriters: "), ",")
					for _, s := range split {
						drama.Screenwriters = append(drama.Screenwriters, strings.TrimSpace(s))
					}
				} else if strings.HasPrefix(child.Text, "Genres:") {
					split := strings.Split(strings.TrimPrefix(child.Text, "Genres: "), ",")
					for _, s := range split {
						drama.Genres = append(drama.Genres, strings.TrimSpace(s))
					}
				} else if strings.HasPrefix(child.Text, "Tags:") {
					split := strings.Split(strings.TrimPrefix(child.Text, "Tags: "), ",")
					for _, s := range split {
						drama.Tags = append(drama.Tags, strings.TrimSpace(s))
					}
				} else if strings.HasPrefix(child.Text, "Related Content") {
					child.ForEach(".title", func(k int, c *colly.HTMLElement) {
						drama.RelatedContent = append(drama.RelatedContent, c.ChildText("a.text-primary"))
					})
				} else if strings.HasPrefix(child.Text, "Score:") {
					scoreText := strings.TrimPrefix(child.Text, "Score: ")
					if len(scoreText) == 0 {
						return
					}
					num, err := strconv.ParseFloat(strings.Split(scoreText, " ")[0], 32)
					if err == nil {
						drama.Rating = float32(num)
					}
				}
			}
		})
	})
	err := scraper.Visit(url)
	return drama, err
}
