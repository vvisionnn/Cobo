package manhuatai

import "testing"

func TestSearch(t *testing.T) {
	comics, err := Search("记")
	if err != nil { t.Error(err) }

	for _, comic := range comics {
		t.Log(*comic)
	}
}
