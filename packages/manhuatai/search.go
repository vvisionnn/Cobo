package manhuatai

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
)

func Search(title string) ([]*Comic, error) {
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("cache", "false")
	params.Set("search_key", title)
	req.URL.RawQuery = params.Encode()

	// set header
	for hKey := range headers {
		req.Header.Set(hKey, headers[hKey])
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var comics []*Comic

	// Find the review items
	doc.Find("#js_comicSortList > li.comic-item").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		Name := s.Find("p.title").Text()
		link, _ := s.Find("a").Attr("href")
		cover, _ := s.Find("a > div > img").Attr("data-src")
		latestChapter := s.Find("a > div > span.chapter").Text()
		comics = append(comics, &Comic{
			Name:          Name,
			Url:           link,
			Cover:         "https:" + cover,
			LatestChapter: latestChapter,
		})
	})

	return comics, nil
}