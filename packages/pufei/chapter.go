package pufei

import (
	"encoding/base64"
	"errors"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
)

func (c *Chapter) getHTMLContentString() (string, error) {
	return GetContent(c.Url)
}

func (c *Chapter) getEncodedString() (string, error) {
	var encodedImageListString string
	chapterStr, err := c.getHTMLContentString()
	if err != nil {
		return encodedImageListString, nil
	}

	var re = regexp.MustCompile(`(?m)cp="(.*?)"`)
	groups := re.FindStringSubmatch(chapterStr)

	if len(groups) < 2 {
		tempErr := errors.New("cannot find encoded image list")
		return encodedImageListString, tempErr
	}

	encodedImageListString = groups[1]
	return encodedImageListString, nil
}

func (c *Chapter) GetImageSuffixList() ([]string, error) {
	var imageSuffixList []string
	encodedImageSuffixString, err := c.getEncodedString()
	if err != nil {
		return imageSuffixList, err
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedImageSuffixString)
	if err != nil {
		return imageSuffixList, err
	}

	vm := otto.New()
	val, err := vm.Run(string(decodedBytes))
	if err != nil {
		return imageSuffixList, err
	}

	decodeImageSuffixString, err := val.ToString()
	if err != nil {
		return imageSuffixList, err
	}

	imageSuffixList = strings.Split(decodeImageSuffixString, ",")
	return imageSuffixList, nil
}

func (c *Chapter) GetImageList() ([]string, error) {
	var imageList []string
	imageSuffixList, err := c.GetImageSuffixList()
	if err != nil {
		return imageList, err
	}
	if len(imageSuffixList) == 0 {
		tempErr := errors.New("no image list found")
		return imageList, tempErr
	}

	if imageSuffixList[0][:7] == "http://" || imageSuffixList[0][:8] == "https://" {
		return imageSuffixList, nil
	}

	// combine host and suffix
	imageList = make([]string, len(imageSuffixList))
	for i, suf := range imageSuffixList {
		imageList[i] = ImageHost + suf
	}

	return imageList, nil
}

func NewChapterFromSuffixUrl(suffix string) (*Chapter, error) {
	if len(suffix) < 1 { return nil, errors.New("suffix error") }
	if suffix[0] == '/' { suffix = suffix[1:] }

	return &Chapter{
		Url: baseUrl + suffix,
	}, nil
}