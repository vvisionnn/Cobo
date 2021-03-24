package Cobo

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestSearchPufeiComic(t *testing.T) {

	title := "一"
	comics ,err := SearchPufeiComic(title)
	if err != nil { t.Error(err) }

	for _, _c := range comics {
		t.Log(*_c)
	}
}

func TestSearchManhuatai(t *testing.T) {
	title := "ok"
	comics ,err := SearchManhuatai(title)
	if err != nil { t.Error(err) }

	for _, _c := range comics {
		t.Log(*_c)
	}
}

func TestGetPufeiComicDetail(t *testing.T) {
	suffix := "/manhua/419"
	c, err := GetPufeiComicDetail(suffix)
	if err != nil { t.Error(err) }

	_j, _ := json.Marshal(c)
	t.Log(string(_j))
	//t.Log(*c.ComicDetail)
}


func TestGetManhuataiComicDetail(t *testing.T) {
	suffix := "/sjbwclzy/"
	c, err := GetManhuataiComicDetail(suffix)
	if err != nil { t.Error(err) }

	_j, _ := json.Marshal(c)
	t.Log(string(_j))
}


func TestGetPufeiChapterImageList(t *testing.T) {
	suffix := "/manhua/419/647556.html"
	images, err := GetPufeiChapterImageList(suffix)
	if err != nil { t.Error(err) }

	t.Log(images)
}

func TestGetManhuataiChapterImageList(t *testing.T) {
	suffix := "/doupocangqiong/dpcq_1h.html"
	images, err := GetManhuataiChapterImageList(suffix)
	if err != nil { t.Error(err) }

	for i := range images {
		images[i], _ = url.QueryUnescape(images[i])
	}

	content, err := json.Marshal(images)
	if err != nil {
		t.Error(err)
	}

	t.Log("string(content):", string(content))
}