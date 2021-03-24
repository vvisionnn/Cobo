package manhuatai

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

func (c *Chapter) GetChapterInfoV1() (*Comic, error) {
	var err error
	var comic Comic

	headers := map[string]string{}
	for k := range defaultHeaders { headers[k] = defaultHeaders[k] }
	for k := range chapterHeaders { headers[k] = chapterHeaders[k] }

	htmlContent, err := GetUrlContent(c.Url, headers)
	if err != nil { return nil, err }

	var comicIdRe = regexp.MustCompile(`(?m)window.comicInfo={comic_id:(.*?),`)
	var chapterNewIdRe = regexp.MustCompile(`(?m),chapter_newid:"(.*?)",chapter_id`)

	// parse info string
	comicIdMatch := comicIdRe.FindStringSubmatch(htmlContent)
	chapterNewIdMatch := chapterNewIdRe.FindStringSubmatch(htmlContent)
	if len(comicIdMatch) < 2 || len(chapterNewIdMatch) < 2 {
		return nil, errors.New("cannot match enough chapter info")
	}

	_id, err := strconv.Atoi(comicIdMatch[1])
	if err != nil { return nil, errors.New(fmt.Sprintf("get comic id error: %v", err)) }
	c.ChapterNewId = chapterNewIdMatch[1]

	comic = Comic{
		ComicId: _id,
		CurrentChapter: c,
	}

	return &comic, nil
}

func GetChapterInfoV10(c *Comic) (*Comic, error) {
	params := url.Values{}
	params.Add("product_id", "2")
	params.Add("productname", "mht")
	params.Add("platformname", "wap")
	params.Add("comic_id", strconv.Itoa(c.ComicId))
	params.Add("chapter_newid", c.CurrentChapter.ChapterNewId)
	params.Add("isWebp", "1")
	params.Add("quality", DefinitionLow)

	apiUrl := chapterInfoUrl + "?" + params.Encode()
	headers := map[string]string{}
	headers["referer"] = c.CurrentChapter.Url
	for k := range defaultHeaders { headers[k] = defaultHeaders[k] }
	for k := range chapterInfoXMLHTTPHeaders { headers[k] = chapterInfoXMLHTTPHeaders[k] }

	resp, err := GetUrlContent(apiUrl, headers)
	if err != nil { return nil, errors.New(fmt.Sprintf("get api 10 error: %v", err)) }


	type tempJson struct {
		Data    *Comic `json:"data"`
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
	var tmp tempJson

	if err = json.Unmarshal([]byte(resp), &tmp); err != nil {
		return nil, errors.New(fmt.Sprintf("unmarshal api resp error: %v", err))
	}
	return tmp.Data, nil
}

func (c *Chapter) GetAllImageUrl() ([]string, error) {
	var err error
	var comic *Comic

	if comic, err = c.GetChapterInfoV1(); err != nil {
		return nil, errors.New(fmt.Sprintf("get chapter info v1 error: %v", err))
	}

	if comic, err = GetChapterInfoV10(comic); err != nil {
		return nil, errors.New(fmt.Sprintf("get chapter info v10 error: %v", err))
	}

	for i := range comic.CurrentChapter.ChapterImgList {
		comic.CurrentChapter.ChapterImgList[i], _ = url.QueryUnescape(comic.CurrentChapter.ChapterImgList[i])
	}

	return comic.CurrentChapter.ChapterImgList, nil
}

func NewChapterFromSuffix(suffix string) (*Chapter, error) {
	if len(suffix) < 1 { return nil, errors.New("suffix error") }
	if suffix[0] != '/' { suffix = "/" + suffix }
	return &Chapter{
		Url: base + suffix,
	}, nil
}