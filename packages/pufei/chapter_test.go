package pufei

import "testing"

func TestChapter_GetImageList(t *testing.T) {
	_chapter, err := NewChapterFromSuffixUrl("/manhua/419/647556.html")
	if err != nil {
		t.Error(err)
	}

	imageList, err := _chapter.GetImageList()
	if err != nil { t.Error(err) }

	t.Log(imageList)
}