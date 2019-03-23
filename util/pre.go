package util

import (
	"github.com/phpgao/godown/downloader"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path"
)

func Download(target string, c chan int) {
	//p1

	_, _, l, url := Check(target)
	log.Error(url)
	d := &downloader.Downloader{
		Url:      url,
		Filename: path.Base(url),
		Length:   l,
	}

	//err := d.Normal()
	err := d.Threaded()
	if err != nil {
		log.Error(err)
	}
	c <- 1
}

//Check 检查下载地址是否满足需求
func Check(url string) (Len bool, Range bool, Length int64, RUrl string) {
	client := &http.Client{}
	req, _ := http.NewRequest("HEAD", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	log.Info(resp.Request)
	log.Info(resp.Request.Response)
	log.Info(resp.StatusCode)
	log.Info(resp.Request.URL)
	err = resp.Body.Close()
	if err != nil {
		log.Error(err)
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
	//log.Infof("request url = %s", resp.Request.URL.String())
	RUrl = resp.Request.URL.String()
	return
}
