package util

import (
	"github.com/phpgao/godown/downloader"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
)

var Logger *logrus.Logger

func Download(target, dir, filePath string, c chan error, l int64, check map[string]string) {
	// p1
	var err error
	Logger.Debugf("Download dir is %s", dir)
	_, _, length, finalUrl := Check(target)
	if filePath == "" {
		filePath = filepath.Base(finalUrl)
	}
	//<-c
	// Init the downloader
	d := &downloader.Downloader{
		Url:      finalUrl,
		Dir:      dir,
		Filename: filePath,
		Length:   length,
		Limit:    l,
		Check:    check,
	}

	// Do the req
	// err := d.Normal()
	err = d.Threaded()
	if err != nil {
		Logger.Error(err)
	}
	c <- err
}

// Check 检查下载地址是否满足需求
func Check(url string) (Len bool, Range bool, Length int64, RUrl string) {
	client := &http.Client{}
	req, _ := http.NewRequest("HEAD", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		Logger.Error(err)
		return
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
