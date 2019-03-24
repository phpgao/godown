package main

import (
	"github.com/phpgao/godown/util"
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func main() {
	c := make(chan int)
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	//p1
	//请求文件尝试用HEAD获取下载文件大小，测试服务器是否支持range，并判断下载类型
	//同时，加载代理列表供后续使用
	urls := []string{
		//"http://mirrors.ustc.edu.cn/debian/extrafiles",
		//"http://mirrors.ustc.edu.cn/centos/filelist.gz",
		//"http://gzm.com:8080/MacFamilyTree_8.3.6_WaitsUn.com.dmg",
		//"http://gzm.com:8080/MacFamilyTree_8.3.6_WaitsUn.com.dmg",
		//"http://gzm.com:8080/321",
		//"http://gzm.com:8080/README.md",
		"http://gzm.com:8080/Microsoft_Office_16.20.18120801_WaitsUn.com.dmg",
		//"https://mirrors.163.com/debian-cd/9.8.0/amd64/iso-cd/debian-9.8.0-amd64-netinst.iso",
		//"http://iso.mirrors.ustc.edu.cn/debian-cd/9.8.0/amd64/iso-cd/debian-9.8.0-amd64-netinst.iso",
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
