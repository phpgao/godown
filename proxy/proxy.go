package proxy

import (
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/url"
	"time"
)

var Logger *logrus.Logger

type Pool struct {
	Proxies []url.URL
}


func Check(proxy url.URL, target string)(ok bool){
	var timeout = time.Duration(5 * time.Second)

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(&proxy)},
		Timeout:   timeout,
	}
	resp, err := client.Get(target)
	if err != nil {
		return
	}

	if resp.StatusCode != 200{
		return
	}

	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return true
}


func (p *Pool) add(proxy string) {
	temp, err := url.Parse(proxy)
	if err != nil {
		Logger.Info(proxy)
	}

	p.Proxies = append(p.Proxies, *temp)
}
