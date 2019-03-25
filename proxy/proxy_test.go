package proxy

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_Proxy(t *testing.T) {
	p := &Pool{
		Proxies:   []url.URL{
			{
				Host : "37.194.50.174:8080",
			},
			{
				Host : "182.52.51.59:38238",
			},
		},
	}


	c := make(chan int)
	for _,prox :=range p.Proxies{

		go func(prox url.URL, p *Pool,c chan int) {
			fmt.Println("go")
			if ok := Check(prox, "https://www.baidu.com");ok{
				t.Log("pass")
			}else{
				t.Error("proxy error")
			}

			fmt.Println("go2")
			c<-1
		}(prox, p,c)

	}
	for range p.Proxies{
		<-c
	}


}
//
//func Test_Division_2(t *testing.T) {
//	t.Error("就是不通过")
//}