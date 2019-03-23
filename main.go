package main

import "github.com/phpgao/godown/util"

func main() {
	c := make(chan int)
	//p1
	//请求文件尝试用HEAD获取下载文件大小，测试服务器是否支持range，并判断下载类型
	//同时，加载代理列表供后续使用
	urls := []string{
		"http://mirrors.ustc.edu.cn/debian/extrafiles",
		"http://t.cn/ExmR3wK",
	}

	for _, url := range urls {
		go util.Download(url, c)
	}

	//p2
	//下载模式分为，单线程下载，多线程下载单代理，多线程多代理下载
	for range urls {
		<-c
	}
	//p3
	//合并文件
}
