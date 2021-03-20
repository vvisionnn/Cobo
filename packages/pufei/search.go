package pufei

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/vvisionnn/Cobo/utils"
	"strings"
)

func buildSearchUrl(title string) (string, error) {
	gbkTitle, err := utils.Encode(title, "gbk")
	if err != nil {
		return "", err
	}
	searchUrl := fmt.Sprintf(searchUrlFmt, gbkTitle)
	return searchUrl, nil
}

func getAllComicFromContent(content string) ([]*Comic, error) {
	var comics []*Comic

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	doc.Find("#detail > li").Each(func(index int, s *goquery.Selection) {
		Name := s.Find("a > h3").Text()
		link, _ := s.Find("a").Attr("href")
		cover, _ := s.Find("a > div > img").Attr("data-src")
		latestChapter := s.Find("dl:nth-child(4) > dd").Text()
		comics = append(comics, &Comic{
			Name:         Name,
			Url:           link,
			Cover:         cover,
			LatestChapter: latestChapter,
		})
	})

	return comics, nil
}

func Search(title string) ([]*Comic, error) {
	searchUrl, err := buildSearchUrl(strings.Trim(title, " "))
	if err != nil {
		return nil, err
	}
	content, err := GetContent(searchUrl)
	if err != nil {
		return nil, err
	}

	return getAllComicFromContent(content)
}
