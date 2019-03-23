package downloader

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Downloader struct {
	Filename string
	threads  int
	Url      string
	Length   int64
}

func (d *Downloader) Normal() error {
	log.Info("start download" + d.Url)
	req, err := http.NewRequest("GET", d.Url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Errorf("bad status: %s", resp.Status)

		return errors.New(resp.Status)
	}

	file, err := os.OpenFile(d.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func (d *Downloader) Threaded() (err error) {
	log.Infof("Begin %s", d.Url)
	//d.Length = 101
	log.Infof("Length %d", d.Length)
	var wg sync.WaitGroup
	//We need a max number of threads
	var limit int64
	limit = 3
	//The missing part
	remain := d.Length % limit
	//size each thread
	length := d.Length / limit
	log.Infof("remain %d", remain)
	//init file with given Content Length
	log.Infof("Init %s", d.Filename)

	f, err := os.Create(d.Filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Infof("Fill %s with %d * 0", d.Filename, d.Length)
	if err = f.Truncate(d.Length); err != nil {
		log.Fatal(err)
		return
	}

	var i int64
	for i = 0; i < limit; i++ {
		wg.Add(1)
		start := length * i
		end := length * (i + 1)

		if i == limit-1 {
			end += remain
		}
		//log.Infof("bytes=%d-%d", start, end-1)
		go func(start, end, i int64) {
			log.Infof("Thread %d begin ", i)
			c := &http.Client{}
			c.Timeout = time.Second * 10
			req, _ := http.NewRequest("GET", d.Url, nil)

			log.Warn(d.Url)

			header := fmt.Sprintf("bytes=%d-%d", start, end-1)
			req.Header.Add("Range", header)
			req.Header.Add("Content-Type", "text/html; charset=UTF-8")
			req.Header.Add("Connection", "chunked")
			req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36")
			req.Header.Add("Referer", "http://mirrors.ustc.edu.cn/")
			req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
			req.Header.Add("Accept-Encoding", "gzip, deflate")
			req.Header.Add("DNT", "1")
			req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
			//cookie := http.Cookie{Name: "addr", Value: "222.64.193.9"}
			//req.AddCookie(&cookie)
			resp, err := c.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			log.Info(resp.Request)
			log.Info(resp.Request.Response)
			f, err := os.OpenFile(d.Filename, os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					log.Fatal(err)
				}
			}()

			defer func() {
				err := f.Close()
				if err != nil {
					log.Fatal(err)
				}
			}()

			if _, err := f.Seek(start, 0); err != nil {
				panic(err)
			}

			_, err = io.Copy(f, resp.Body)

			log.Infof("file %s,part %d done!", d.Filename, i+1)
			wg.Done()
		}(start, end, i)

	}
	wg.Wait()
	return
}
