package pufei

import "testing"

func TestSearch(t *testing.T) {
	comics, err := Search("斗")
	if err != nil {
		t.Error(err)
	}

	for _, c := range comics {
		t.Log(*c)
	}
}
