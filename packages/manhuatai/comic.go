package manhuatai

import (
	"errors"
	"github.com/gocolly/colly"
)

func (c *Comic) GetAllChaptersAndExtraInfo() error {
	var err error
	var allUrl []string
	var allTitle []string
	collector := getComicDefaultCollector()
	collector.OnHTML("#js_chapter_list", func(element *colly.HTMLElement) {
		allUrl = element.ChildAttrs("li.item > a", "href")
		allTitle = element.ChildAttrs("li.item > a", "title")
	})

	collector.OnHTML("#js_desc_content", func(element *colly.HTMLElement) {
		c.Description = element.Text
	})

	collector.OnHTML("#detail > div > div > h1", func(element *colly.HTMLElement) {
		author := element.ChildText("span")
		c.Name = element.Text[:len(element.Text)-len(author)]
	})
	collector.OnHTML("#detail > img", func(element *colly.HTMLElement) {
		c.Cover = "https:" + element.Attr("data-src")
	})

	if err = collector.Visit(c.Url); err != nil {
		return err
	}
	if len(allUrl) != len(allTitle) {
		return errors.New("the number of titles is different from the number of URLs")
	}

	c.Chapters = make([]*Chapter, len(allUrl))
	for index := 0; index < len(allUrl); index++ {
		c.Chapters[index] = &Chapter{
			Title: allTitle[index],
			Url:   allUrl[index],
			Info:  nil,
		}
	}
	return nil
}

func NewComic(partUrl string) *Comic {
	fullUrl := base + partUrl
	return &Comic{
		Url: fullUrl,
	}
}
