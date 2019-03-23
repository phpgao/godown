package proxy

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type Pool struct {
	Proxies []url.URL
}

func (p *Pool) Test(proxy url.URL, target string) (ok bool) {
	var timeout = time.Duration(15 * time.Second)

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(&proxy)},
		Timeout:   timeout,
	}
	resp, err := client.Get(target)
	if err != nil {
		return ok
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
		log.Info(proxy)
	}

	p.Proxies = append(p.Proxies, *temp)
}
