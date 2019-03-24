package downloader

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/EDDYCJY/fake-useragent"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"io"
	"mime"
	"net/http"
	"os"
	"sync"
)

const (
	HTTPPartialContent = 206
)

var Logger *logrus.Logger

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	// fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

type Downloader struct {
	Filename string
	Dir      string
	threads  int
	Url      string
	Length   int64
	Limit    int64
}

func (d *Downloader) Normal() error {
	Logger.Info("start download" + d.Url)
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
			Logger.Fatal(err)
		}
	}()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		Logger.Errorf("bad status: %s", resp.Status)

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
			Logger.Fatal(err)
		}
	}()

	return nil
}

func (d *Downloader) Threaded() (err error) {
	Logger.Infof("Begin %s", d.Url)
	// d.Length = 101
	Logger.Infof("Length %d", d.Length)
	var wg sync.WaitGroup
	// We need a max number of threads

	limit := d.Limit
	// The missing part
	remain := d.Length % limit
	// size each thread
	length := d.Length / limit
	Logger.Infof("remain %d", remain)
	Logger.Infof("length %d", length)
	// init file with given Content Length
	Logger.Infof("Init %s", d.Filename)
	// First Create the file
	f, err := os.Create(d.Filename)
	if err != nil {
		Logger.Fatal(err)
		return
	}
	Logger.Infof("Fill %s with %d * 0", d.Filename, d.Length)
	// Fill it with zero
	if err = f.Truncate(d.Length); err != nil {
		Logger.Fatal(err)
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
		// Logger.Infof("bytes=%d-%d", start, end-1)
		go func(start, end, i int64) {
			// func(start, end, i int64) {
			// func(start, end, i int64) {
			Logger.Infof("Thread %d begin ", i)
			c := &http.Client{}
			// c.Timeout = time.Second * 10
			req, _ := http.NewRequest("GET", d.Url, nil)

			header := fmt.Sprintf("bytes=%d-%d", start, end-1)
			req.Header.Add("Range", header)
			req.Header.Add("Content-Type", "text/html; charset=UTF-8")
			// req.Header.Add("Connection", "chunked")
			req.Header.Add("User-Agent", randomUA())
			// req.Header.Add("Referer", "http:// mirrors.ustc.edu.cn/")
			req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
			req.Header.Add("Accept-Encoding", "gzip")
			req.Header.Add("DNT", "1")
			req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

			resp, err := c.Do(req)
			if err != nil {
				Logger.Fatal(err)
			}
			Logger.Debug(header)
			Logger.Debug(resp.Request)
			Logger.Debug(resp.Header)
			Logger.Debug(resp.StatusCode)

			// Download
			if resp.StatusCode == HTTPPartialContent {

			}
			f, err := os.OpenFile(d.Filename, os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}

			defer func() {
				Logger.Debugf("task %d resp close", i)
				err := resp.Body.Close()
				if err != nil {
					Logger.Fatal(err)
				}
			}()

			defer func() {
				Logger.Debugf("task %d file close", i)
				err := f.Close()
				if err != nil {
					Logger.Fatal(err)
				}
			}()

			if _, err := f.Seek(start, 0); err != nil {
				panic(err)
			}
			var reader io.ReadCloser
			switch resp.Header.Get("Content-Encoding") {
			case "gzip":
				reader, err = gzip.NewReader(resp.Body)
				defer func() {
					err := reader.Close()
					if err != nil {
						Logger.Fatal(err)
					}
				}()
			default:
				reader = resp.Body
			}

			// _, _ = ioutil.ReadAll(resp.Body)
			// n, err := f.Write(body)
			counter := &WriteCounter{}
			n, err := io.Copy(f, io.TeeReader(reader, counter))
			// n, err := copyBuffer(f, resp.Body, nil)
			if err != nil {
				Logger.Error(err)
			}
			Logger.Debugf("task %d written %d bytes", i, n)

			Logger.Infof("file %s,task id = %d done!", d.Filename, i)
			wg.Done()
		}(start, end, i)

	}
	wg.Wait()
	return
}

func getFileNameFrom(u string, header http.Header) (name string, err error) {
	// Header first
	if cd := header.Get("Content-Disposition"); cd != "" {
		_, params, err := mime.ParseMediaType(`attachment;filename="foo.png"`)
		if err != nil {
			Logger.Warn(err)
		}

		if params["filename"] != "" {
			name = params["filename"]
		}

	}

	// Then from url,if url ends with a slash
	//

	return
}

func randomUA() string {
	chrome := browser.Chrome()
	Logger.Debugf("Chrome: %s", chrome)
	return chrome
}
