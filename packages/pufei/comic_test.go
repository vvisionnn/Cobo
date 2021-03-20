package pufei

import (
	"testing"
)

func TestComic_GetAllComicInfo(t *testing.T) {
	comic, err := NewComicFromUrlSuffix("/manhua/31275")
	if err != nil { t.Error(err) }
	if err := comic.GetAllComicInfo(); err != nil {
		t.Error(err)
	}

	t.Log(*comic)
}