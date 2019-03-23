package downloader

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

type Downloader struct {
	Filename string
	threads  int
	Url      string
}

func (this *Downloader) Normal() error {
	log.Info("start download" + this.Url)
	req, err := http.NewRequest("GET", this.Url, nil)
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

		return nil
	}

	file, err := os.OpenFile(this.Filename, os.O_WRONLY|os.O_CREATE, 0666)
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
