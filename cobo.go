package Cobo

import (
	"github.com/vvisionnn/Cobo/packages/manhuatai"
	"github.com/vvisionnn/Cobo/packages/pufei"
)

type ComicPreview struct {
	Name          string `json:"name"`
	Url           string `json:"url"`
	Cover         string `json:"cover"`
	LatestChapter string `json:"latest_chapter"`
}

type ComicDetail struct {
	Description string            `json:"description"`
	Chapters    []*ChapterGeneral `json:"chapters"`
}

type Comic struct {
	ComicPreview
	ComicDetail
}

type ChapterGeneral struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}


func SearchPufeiComic(title string) ([]*ComicPreview, error) {
	_cs, err := pufei.Search(title)
	if err != nil {
		return nil, err
	}

	comics := make([]*ComicPreview, len(_cs))
	for index, c := range _cs {
		chapters := make([]*ChapterGeneral, len(c.Chapters))
		for _index, _c := range c.Chapters {
			chapters[_index] = &ChapterGeneral{
				Url:   _c.Url,
				Title: _c.Title,
			}
		}
		comics[index] = &ComicPreview{
			Name:          c.Name,
			Url:           c.Url,
			Cover:         c.Cover,
			LatestChapter: c.LatestChapter,
		}
	}
	return comics, nil
}

func SearchManhuatai(title string) ([]*ComicPreview, error) {
	_cs, err := manhuatai.Search(title)
	if err != nil {
		return nil, err
	}

	comics := make([]*ComicPreview, len(_cs))
	for index, c := range _cs {
		chapters := make([]*ChapterGeneral, len(c.Chapters))
		for _index, _c := range c.Chapters {
			chapters[_index] = &ChapterGeneral{
				Url:   _c.Url,
				Title: _c.Title,
			}
		}
		comics[index] = &ComicPreview{
			Name:          c.Name,
			Url:           c.Url,
			Cover:         c.Cover,
			LatestChapter: c.LatestChapter,
		}
	}
	return comics, nil
}

func GetPufeiComicDetail(suffix string) (*Comic, error) {
	c, err := pufei.NewComicFromUrlSuffix(suffix)
	if err != nil { return nil, err }
	if err := c.GetAllComicInfo(); err != nil {
		return nil, err
	}

	chapters := make([]*ChapterGeneral, len(c.Chapters))
	for _index, _c := range c.Chapters {
		chapters[_index] = &ChapterGeneral{
			Url:   _c.Url,
			Title: _c.Title,
		}
	}
	result := Comic{}
	result.Name = c.Name
	result.Url = c.Url
	result.Cover = c.Cover
	result.LatestChapter = c.LatestChapter
	result.Description = c.Description
	result.Chapters = chapters

	return &result, nil
}

func GetManhuataiComicDetail(suffix string) (*Comic, error) {
	c := manhuatai.NewComic(suffix)
	if err := c.GetAllChaptersAndExtraInfo(); err != nil {
		return nil, err
	}

	chapters := make([]*ChapterGeneral, len(c.Chapters))
	for _index, _c := range c.Chapters {
		chapters[_index] = &ChapterGeneral{
			Url:   _c.Url,
			Title: _c.Title,
		}
	}
	result := Comic{}
	result.Name = c.Name
	result.Url = c.Url
	result.Cover = c.Cover
	result.LatestChapter = c.LatestChapter
	result.Description = c.Description
	result.Chapters = chapters

	return &result, nil
}

func GetPufeiChapterImageList(suffix string) ([]string, error) {
	_chapter, err := pufei.NewChapterFromSuffixUrl(suffix)
	if err != nil { return nil, err }

	return _chapter.GetImageList()
}

func GetManhuataiChapterImageList(suffix string) ([]string, error) {
	_chapter, err := manhuatai.NewChapterFromSuffix(suffix)
	if err != nil { return nil, err }

	return _chapter.GetAllImageUrl()
}