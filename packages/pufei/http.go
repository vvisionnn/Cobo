package pufei

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/vvisionnn/Cobo/utils"
	"io"
	"io/ioutil"
	"net/http"
)

var client = http.Client{}

func GetContent(url string) (string, error) {
	var bodyString string
	// todo: extract request as a single package
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return bodyString, err
	}

	// add headers
	for k := range defaultHeader {
		req.Header.Set(k, defaultHeader[k])
	}

	resp, err := client.Do(req)
	if err != nil {
		return bodyString, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resErr := errors.New(fmt.Sprintf("request status code %d error", resp.StatusCode))
		return bodyString, resErr
	}

	// todo: to blog article
	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return bodyString, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	bodyByte, err := ioutil.ReadAll(reader)
	if err != nil {
		return bodyString, err
	}

	return utils.GBKToUTF8(string(bodyByte))
}
