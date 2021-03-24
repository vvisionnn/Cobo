package manhuatai

const base = "https://m.manhuatai.com"
const searchUrl = "https://m.manhuatai.com/sort/all.html"
const chapterInfoUrl = "https://m.manhuatai.com/api/getchapterinfov2"

const DefinitionLow = "low"
const DefinitionMiddle = "middle"
const DefinitionHigh = "high"

var defaultHeaders = map[string]string{
	"accept-encoding": "gzip, deflate, br",
	"accept-language": "en",
	"cookie":          "kmh-habit=%7B%22mode%22%3A%22scroll%22%7D; user=%7B%22type%22%3A%22device%22%2C%22Cgold%22%3A0%2C%22coins%22%3A0%2C%22Ulevel%22%3A1%7D",
	"user-agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
}

var chapterHeaders = map[string]string{
	"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"upgrade-insecure-requests": "1",
	"cache-control":             "max-age=0",
	"sec-fetch-dest":            "document",
	"sec-fetch-mode":            "navigate",
	"sec-fetch-site":            "same-origin",
	"sec-fetch-user":            "?1",
}

var chapterInfoXMLHTTPHeaders = map[string]string{
	"accept":           "*/*",
	"x-requested-with": "XMLHttpRequest",
	"sec-fetch-dest":   "empty",
	"sec-fetch-mode":   "cors",
	"sec-fetch-site":   "same-origin",
}

type definitionInfo struct {
	Low    string `json:"low"`
	Middle string `json:"middle"`
	High   string `json:"high"`
}

//type ChapterInfo struct {
//	Definition          definitionInfo `json:"definition"`
//	ChapterName         string         `json:"chapter_name"`
//	ChapterNewId        string         `json:"chapter_newid"`
//	ChapterId           uint           `json:"chapter_id"`
//	ChapterDomainSuffix string         `json:"chapter_domain_suffix"`
//	ChapterDomain       string         `json:"chapter_domain"`
//	StartNum            uint           `json:"start_num"`
//	EndNum              uint           `json:"end_num"`
//	Price               uint           `json:"price"`
//	ChapterImageAddr    string         `json:"chapter_image_addr"`
//	CreateDate          float64        `json:"create_date"`
//	//Rule                string         `json:"rule"`
//}

type Chapter struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	//Info  *ChapterInfo `json:"info"`

	ChapterName         string   `json:"chapter_name"`
	ChapterNewId        string   `json:"chapter_newid"`
	ChapterId           uint     `json:"chapter_id"`
	ChapterDomainSuffix string   `json:"chapter_domain_suffix"`
	ChapterDomain       string   `json:"chapter_domain"`
	StartNum            uint     `json:"start_num"`
	EndNum              uint     `json:"end_num"`
	Price               uint     `json:"price"`
	ChapterImageAddr    string   `json:"chapter_image_addr"`
	CreateDate          int64    `json:"create_date"`
	ChapterImgList      []string `json:"chapter_img_list"`
}

type Comic struct {
	ComicId          int    `json:"comic_id"`
	ComicNewid       string `json:"comic_newid"`
	ComicName        string `json:"comic_name"`
	LastChapterId    string `json:"last_chapter_id"`
	LastChapterNewid string `json:"last_chapter_newid"`
	LastChapterName  string `json:"last_chapter_name"`
	ShowType         int    `json:"show_type"`
	Readtype         int    `json:"readtype"`
	ComicStatus      int    `json:"comic_status"`
	ChargePaid       int    `json:"charge_paid"`
	ChargeCoinFree   int    `json:"charge_coin_free"`
	UpdateTime       int64  `json:"update_time"`
	//boo_virtual_coin: !0,
	IsCopyright int `json:"is_copyright"`

	Definition definitionInfo `json:"definition"`

	CurrentChapter *Chapter `json:"current_chapter"`
	PrevChapter    *Chapter `json:"prev_chapter"`
	NextChapter    *Chapter `json:"next_chapter"`
}

type ComicDetail struct {
	//Title         string     `json:"title"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Cover       string `json:"cover"`
	//InfoCover     string     `json:"info_cover"`
	LatestChapter string     `json:"latest_chapter"`
	Chapters      []*Chapter `json:"chapters"`
}
