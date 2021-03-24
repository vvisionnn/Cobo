package manhuatai

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

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

func NewChapterFromSuffix(suffix string) (*Chapter, error) {
	if len(suffix) < 1 { return nil, errors.New("suffix error") }
	if suffix[0] != '/' { suffix = "/" + suffix }
	return &Chapter{
		Url: base + suffix,
	}, nil
}