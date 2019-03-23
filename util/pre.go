package util

import (
	"github.com/phpgao/godown/downloader"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path"
)

func Download(url string, c chan int) {
	//p1

	_, _ = Check(&url)
	d := &downloader.Downloader{
		Url:      url,
		Filename: path.Base(url),
	}

	err := d.Normal()

	log.Info(err)
	c <- 1
}

//Check 检查下载地址是否满足需求
func Check(url *string) (Len bool, Range bool) {
	resp, err := http.Get(*url)

	if err != nil {
		panic(err)
	}

	_ = resp.Body.Close()

	if resp.ContentLength > 0 {
		Len = true
	}

	if _, ok := resp.Header["Accept-Ranges"]; ok {
		if resp.Header["Accept-Ranges"][0] == "bytes" {
			Range = true
		}
	}
	log.Info(resp.Request.URL.String())
	*url = resp.Request.URL.String()

	return
}
