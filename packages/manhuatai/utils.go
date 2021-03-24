package manhuatai

import (
	"errors"
	"github.com/gocolly/colly"
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"
)

func downloadFile(URL, fileName string, wg *sync.WaitGroup, errChan chan error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		errChan <- errors.New("received non 200 response code")
		wg.Done()
		return
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}

	wg.Done()
}

func getComicDefaultCollector() *colly.Collector {
	var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36"
	const domain = "m.manhuatai.com"
	var collector = colly.NewCollector(
		colly.AllowedDomains(domain),
		colly.UserAgent(userAgent),
	)
	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	return collector
}

func regularJsonStr(Data []byte) []byte {
	reg := regexp.MustCompile("([a-zA-Z]\\w*):")
	regStr := reg.ReplaceAllString(string(Data), `"$1":`)
	return []byte(regStr)
}
