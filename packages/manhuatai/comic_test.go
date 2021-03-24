package manhuatai

import "testing"

func TestComicDetail_GetAllChaptersAndExtraInfo(t *testing.T) {
	comic := NewComic("/yuzhouxiyouji/")
	t.Log(comic.GetAllChaptersAndExtraInfo())
	t.Log(*comic)
}
