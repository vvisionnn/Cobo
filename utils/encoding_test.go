package utils

import (
	"fmt"
	"testing"
)

func TestConvertGBKBytes(t *testing.T) {
	s, _ := ConvertGBKBytes([]byte("斗"))
	t.Log(string(s))
}

func TestUTF8ToGBK(t *testing.T) {
	utf8List := []string{
		"斗",
	}

	for _, seq := range utf8List {
		res, err := UTF8ToGBK(seq)
		if err != nil {
			t.Error(err)
		}
		fmt.Println([]byte(res))
		fmt.Println(res)
	}
}
