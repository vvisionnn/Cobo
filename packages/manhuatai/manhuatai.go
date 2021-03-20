package manhuatai

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

const base = "https://m.manhuatai.com"
const searchUrl = "https://m.manhuatai.com/sort/all.html"

var headers = map[string]string{
	"user-agent":                "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
	"upgrade-insecure-requests": "1",
}

type definitionInfo struct {
	Low    string `json:"low"`
	Middle string `json:"middle"`
	High   string `json:"high"`
}

type ChapterInfo struct {
	Definition          definitionInfo `json:"definition"`
	ChapterName         string         `json:"chapter_name"`
	ChapterNewId        string         `json:"chapter_newid"`
	ChapterId           uint           `json:"chapter_id"`
	ChapterDomainSuffix string         `json:"chapter_domain_suffix"`
	ChapterDomain       string         `json:"chapter_domain"`
	StartNum            uint           `json:"start_num"`
	EndNum              uint           `json:"end_num"`
	Price               uint           `json:"price"`
	ChapterImageAddr    string         `json:"chapter_image_addr"`
	CreateDate          float64        `json:"create_date"`
	Rule                string         `json:"rule"`
}

type Chapter struct {
	Title string       `json:"title"`
	Url   string       `json:"url"`
	Info  *ChapterInfo `json:"info"`
}

type Comic struct {
	//Title         string     `json:"title"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Cover       string `json:"cover"`
	//InfoCover     string     `json:"info_cover"`
	LatestChapter string     `json:"latest_chapter"`
	Chapters      []*Chapter `json:"chapters"`
}

func getComicDefaultCollector() *colly.Collector {
	var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36"
	const domain = "m.manhuatai.com"
	var collector = colly.NewCollector(
		colly.AllowedDomains(domain),
		colly.UserAgent(userAgent),
	)
	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	return collector
}

func regularJsonStr(Data []byte) []byte {
	reg := regexp.MustCompile("([a-zA-Z]\\w*):")
	regStr := reg.ReplaceAllString(string(Data), `"$1":`)
	return []byte(regStr)
}

func (c *Chapter) getDefinitionFromStr(definition string) string {
	if definition == "middle" {
		return c.Info.Definition.Middle
	}
	if definition == "high" {
		return c.Info.Definition.High
	}
	return c.Info.Definition.Low
}

func (c *Chapter) getMaxDefinition() string {
	if c.Info.Definition.High != "" {
		return c.Info.Definition.High
	}
	if c.Info.Definition.Middle != "" {
		return c.Info.Definition.Middle
	}
	return c.Info.Definition.Low
}

func (c *Chapter) GetChapterInfo() error {
	collector := getComicDefaultCollector()
	var err error
	var _infoStr string
	var chapterInfo ChapterInfo
	var definitionsRe = regexp.MustCompile(`(?m)window\.\$definitions=({.*?}),`)
	var infoRe = regexp.MustCompile(`(?m)current_chapter:(.*?),prev_chapter:`)

	collector.OnHTML("body > script:nth-child(3)", func(e *colly.HTMLElement) {
		_infoStr = e.Text
	})
	if err = collector.Visit(c.Url); err != nil {
		fmt.Println("err:", err)
	}

	// parse info string
	definitionsMatch := definitionsRe.FindStringSubmatch(_infoStr)
	infoMatch := infoRe.FindStringSubmatch(_infoStr)
	if len(definitionsMatch) < 2 || len(infoMatch) < 2 {
		return errors.New("cannot match enough chapter info")
	}

	if err = json.Unmarshal(regularJsonStr([]byte(definitionsMatch[1])), &chapterInfo.Definition); err != nil {
		return err
	}
	if err = json.Unmarshal(regularJsonStr([]byte(infoMatch[1])), &chapterInfo); err != nil {
		return err
	}

	chapterInfo.ChapterDomain = "https://mhpic." + chapterInfo.ChapterDomain
	chapterInfo.Rule = strings.Replace(chapterInfo.Rule, "$$", "%d", 1)
	c.Info = &chapterInfo
	return nil
}

func (c *Chapter) GetAllImageUrl(definition string) ([]string, error) {
	var allImageUrl []string
	var err error

	if err = c.GetChapterInfo(); err != nil {
		return nil, err
	}

	// maxDefinition := c.getMaxDefinition()
	imgDefinition := c.getDefinitionFromStr(definition)
	urlFmt := c.Info.ChapterDomain + c.Info.Rule + imgDefinition + ".webp"
	for page := c.Info.StartNum; page <= c.Info.EndNum; page++ {
		allImageUrl = append(allImageUrl, fmt.Sprintf(urlFmt, page))
	}
	return allImageUrl, nil
}

func downloadFile(URL, fileName string, wg *sync.WaitGroup, errChan chan error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		errChan <- errors.New("received non 200 response code")
		wg.Done()
		return
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}

	wg.Done()
}

func (c *Chapter) Download(folderPath string) error {
	var err error
	folderPath = path.Join(path.Dir(folderPath), path.Base(folderPath))
	// if folderPath not exist, create
	if _, err = os.Stat(folderPath); os.IsExist(err) {
		if err = os.RemoveAll(folderPath); err != nil {
			return err
		}
	}

	if err = os.MkdirAll(folderPath, 0755); err != nil {
		return err
	}

	// get all image url
	allImageUrl, err := c.GetAllImageUrl("high")
	if err != nil {
		return err
	}

	// get image path and save
	wg := sync.WaitGroup{}
	errChan := make(chan error)
	for index := 0; index < len(allImageUrl); index++ {
		// download
		imagePath := path.Join(folderPath, fmt.Sprintf("%d.jpg", index+1))
		wg.Add(1)
		go downloadFile(allImageUrl[index], imagePath, &wg, errChan)
	}

	wgDoneChan := make(chan bool)
	go func() {
		wg.Wait()
		close(wgDoneChan)
	}()

	select {
	case <-wgDoneChan:
		fmt.Println("Done")
	case err = <-errChan:
		return err
	}
	return nil
}

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

func NewComic(partUrl string) *Comic {
	fullUrl := base + partUrl
	return &Comic{
		Url: fullUrl,
	}
}

func NewChapterFromSuffix(suffix string) (*Chapter, error) {
	if len(suffix) < 1 { return nil, errors.New("suffix error") }
	if suffix[0] != '/' { suffix = "/" + suffix }
	return &Chapter{
		Url: base + suffix,
	}, nil
}
