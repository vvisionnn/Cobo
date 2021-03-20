package pufei

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func NewComicFromUrlSuffix(suffix string) (*Comic, error) {
	if len(suffix) < 1 { return nil, errors.New("suffix error") }
	if suffix[0] == '/' { suffix = suffix[1:] }

	return &Comic{
		Url: baseUrl + suffix,
	}, nil
}

func (c *Comic) GetComicHTMLContent() (string, error) {
	return GetContent(c.Url)
}

func (c *Comic) GetAllChapters() error {
	content, err := c.GetComicHTMLContent()
	if err != nil {
		return nil
	}
	return c.GetAllChaptersFromContent(content)
}

func (c *Comic) GetAllChaptersFromContent(content string) error {
	var err error
	// parse chapter list from content
	comicStrReader := strings.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(comicStrReader)
	if err != nil {
		return err
	}

	// build chapter struct list
	doc.Find("#chapterList2 > ul > li").Each(func(index int, selection *goquery.Selection) {
		title, _ := selection.Find("a").Attr("title")
		link, _ := selection.Find("a").Attr("href")
		c.Chapters = append(c.Chapters, &Chapter{
			Url:   link,
			Title: title,
		})
	})
	// and return
	return nil
}

func (c *Comic) GetDetail() error {
	content, err := c.GetComicHTMLContent()
	if err != nil {
		return err
	}
	return c.GetDetailFromContent(content)
}

func (c *Comic) GetDetailFromContent(content string) error {
	var err error
	// parse chapter list from content
	comicStrReader := strings.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(comicStrReader)
	if err != nil {
		return err
	}

	c.Name = doc.Find("body > div.main-bar.bar-bg1 > h1").Text()
	c.Description = doc.Find("#bookIntro > p").Text()
	c.Cover = doc.Find("body > div.book-detail > div.cont-list > div.thumb > img").AttrOr("src", "")
	return nil
}

func (c *Comic) GetAllComicInfo() error {
	content, err := c.GetComicHTMLContent()
	if err != nil {
		return err
	}

	if err := c.GetDetailFromContent(content); err != nil {
		return err
	}
	if err := c.GetAllChaptersFromContent(content); err != nil {
		return err
	}
	return nil
}
