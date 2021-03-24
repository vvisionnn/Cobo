package manhuatai

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Search(title string) ([]*ComicDetail, error) {
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("cache", "false")
	params.Set("search_key", title)
	req.URL.RawQuery = params.Encode()

	finalUrl := searchUrl + "?" + params.Encode()
	headers := map[string]string{}
	headers["referer"] = base + "/"
	for hKey := range defaultHeaders {
		headers[hKey] = defaultHeaders[hKey]
	}
	for hKey := range chapterHeaders {
		headers[hKey] = chapterHeaders[hKey]
	}

	content, err := GetUrlContent(finalUrl, headers)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	var comics []*ComicDetail
	// Find the review items
	doc.Find("#js_comicSortList > li.comic-item").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		Name := s.Find("a > p.title").Text()
		link, _ := s.Find("a").Attr("href")
		cover, _ := s.Find("a > div > img").Attr("data-src")
		latestChapter := s.Find("a > div > span.chapter").Text()
		comics = append(comics, &ComicDetail{
			Name:          Name,
			Url:           link,
			Cover:         "https:" + cover,
			LatestChapter: latestChapter,
		})
	})

	return comics, nil
}