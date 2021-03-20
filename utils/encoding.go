package utils

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/axgle/mahonia"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

const (
	utf8s   = "utf-8,utf8,UTF8,UTF-8"
	chinese = "gb2312,Gb2312,GB2312,gbk,GBK,Gbk,gb18030,GB18030,Gb18030"
)

var (
	rune2byte = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
)

func Encode(str string, encoding string) (string, error) {
	if str == "" {
		return "", nil
	}
	isChinese := strings.Contains(chinese, encoding)
	isUtf8 := strings.Contains(utf8s, encoding)

	if !(isChinese || isUtf8) {
		return "", errors.New("Unrecognized encoding")
	}
	return encode(str, isChinese), nil
}

func Decode(str string, encoding string) (string, error) {
	if str == "" {
		return "", nil
	}
	isChinese := strings.Contains(chinese, encoding)
	isUtf8 := strings.Contains(utf8s, encoding)

	if !(isChinese || isUtf8) {
		return "", errors.New("Unrecognized encoding")
	}
	return decode(str, isChinese)
}

func decode(str string, isChinese bool) (string, error) {
	runes := []rune(str)
	length := len(runes)
	buf := new(bytes.Buffer)
	i := 0
	for i < length {
		r := runes[i]
		if r != '%' {
			buf.WriteRune(r)
			i++
		} else {
			b, err := parseEncodeByte(runes[i+1], runes[i+2])
			if err != nil {
				return "", err
			}
			buf.WriteByte(b)
			i += 3
		}
	}
	if isChinese {
		data, _ := ioutil.ReadAll(transform.NewReader(buf, simplifiedchinese.GB18030.NewDecoder()))
		return string(data), nil
	} else {
		return buf.String(), nil
	}
}

func encode(str string, isChinese bool) string {
	runes := []rune(str)
	builder := strings.Builder{}
	for _, r := range runes {

		if !shouldEscape(r) {
			builder.WriteRune(r)
			continue
		}

		p := make([]byte, utf8.UTFMax)
		c := utf8.EncodeRune(p, r)
		var data []byte
		if isChinese {
			data, _ = ioutil.ReadAll(transform.NewReader(bytes.NewReader(p[0:c]), simplifiedchinese.GB18030.NewEncoder()))
		} else {
			data = p[0:c]
		}
		//for _, v := range data {
		//	fmt.Print(v)
		//	fmt.Print(",")
		//}
		//fmt.Println("")

		target := make([]byte, len(data)*2)
		hex.Encode(target, data)

		i := 0
		for ; i < len(target); i += 2 {
			builder.WriteByte('%')
			builder.Write([]byte{target[i], target[i+1]})
		}
	}
	return builder.String()
}

func shouldEscape(r rune) bool {
	return !((r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '-' || r == '_' || r == '.')
}

func parseEncodeByte(high rune, low rune) (byte, error) {
	if high >= 65 && high <= 90 {
		high = high + 32
	}
	if low >= 65 && low <= 90 {
		low = low + 32
	}
	i := byteOfRune(high)
	if i == -1 {
		return 0, errors.New("invalid arguments")
	}
	j := byteOfRune(low)

	if j == -1 {
		return 0, errors.New("invalid arguments")
	}
	return byte((0x0f&i)<<4 + (0x0f & j)), nil
}

func byteOfRune(r rune) int {
	for i, b := range rune2byte {
		if r == rune(b) {
			return i
		}
	}
	return -1
}

func ConvertString(data string, srcEncoding, dstEncoding string) (string, error) {
	var res []byte
	var err error
	srcDecoder := mahonia.NewDecoder(srcEncoding)
	dstDecoder := mahonia.NewDecoder(dstEncoding)

	decodedStr := srcDecoder.ConvertString(data)
	_, res, err = dstDecoder.Translate([]byte(decodedStr), true)

	return string(res), err
}

func ConvertGBKBytes(data []byte) ([]byte, error) {
	var res []byte
	var err error
	decoder := mahonia.NewDecoder("gbk")
	_, res, err = decoder.Translate(data, true)
	return res, err

}

func GBKToUTF8(data string) (string, error) {
	return ConvertString(data, "gbk", "utf-8")
}

func UTF8ToGBK(data string) (string, error) {
	return ConvertString(data, "utf-8", "gbk")
}
