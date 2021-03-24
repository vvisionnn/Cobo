package manhuatai


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
