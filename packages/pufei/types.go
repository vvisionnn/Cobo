package pufei

const baseUrl = "http://m.pufei8.com/"
const comicBaseUrl = baseUrl + "manhua/"
const searchUrlFmt = baseUrl + "e/search/?searchget=1&tbname=mh&show=title,player,playadmin,bieming,pinyin,playadmin&tempid=4&keyboard=%s"

const ImageHost = "http://res.img.youzipi.net/"

//const IMAGE_HOST = "http://res.img.fffimage.com/"	// slow
//const IMAGE_HOST = "http://res.img.pufei.net/"	// not work

var defaultHeader = map[string]string{
	"Host": "m.pufei8.com",
	//"Referer": "http://www.pufei8.com/",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"Accept-Encoding":           "gzip, deflate",
	"Accept-Language":           "en,zh;q=0.9,zh-CN;q=0.8",
	"Connection":                "keep-alive",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.9.168 Version/11.52",
}

type Chapter struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

type Comic struct {
	//Id            uint       `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Url           string     `json:"url"`
	Chapters      []*Chapter `json:"chapters"`
	Cover         string     `json:"cover"`
	LatestChapter string     `json:"latest_chapter"`
}
