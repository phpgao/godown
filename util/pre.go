package util

import (
	"github.com/phpgao/godown/downloader"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"path"
)

var Logger *logrus.Logger

func Download(target string, c chan int) {
	// p1
	// Check target finalUrl meta
	_, _, l, finalUrl := Check(target)

	// Init the downloader
	d := &downloader.Downloader{
		Url:      finalUrl,
		Filename: url.QueryEscape(path.Base(finalUrl)),
		Length:   l,
	}

	// Do the req
	// err := d.Normal()
	err := d.Threaded()
	if err != nil {
		Logger.Error(err)
	}
	c <- 1
}

// Check 检查下载地址是否满足需求
func Check(url string) (Len bool, Range bool, Length int64, RUrl string) {
	client := &http.Client{}
	req, _ := http.NewRequest("HEAD", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	Logger.Debug(resp.Request)
	Logger.Debug(resp.Request.Response)
	Logger.Debug(resp.StatusCode)
	Logger.Debug(resp.Request.URL)
	err = resp.Body.Close()
	if err != nil {
		Logger.Error(err)
	}

	if resp.ContentLength > 0 {
		Len = true
		Length = resp.ContentLength
	}

	if _, ok := resp.Header["Accept-Ranges"]; ok {
		if resp.Header["Accept-Ranges"][0] == "bytes" {
			Range = true
		}
	}

	RUrl = resp.Request.URL.String()
	return
}
