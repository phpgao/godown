package main

import (
	"github.com/phpgao/godown/downloader"
	"github.com/phpgao/godown/util"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var Logger = logrus.New()

// todo
func init() {

}

func setLogger(c *cli.Context) {
	if c.Bool("debug") {
		Logger.SetLevel(logrus.DebugLevel)
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)

	util.Logger = Logger
	downloader.Logger = Logger
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "display more info",
		},
		cli.BoolFlag{
			Name:  "d, debug",
			Usage: "display more debug info",
		},
		cli.StringFlag{
			Name:  "o, out",
			Usage: "file name with or with out path",
		},
		cli.StringFlag{
			Name:  "p, proxy",
			Usage: "proxys",
		},
	}

	app.Action = func(c *cli.Context) error {
		var urls []string
		if c.NArg() > 0 {
			urls = c.Args()
		}

		if len(urls) > 0 {
			setLogger(c)
			Logger.Debug(urls)
			c := make(chan int)

			// p1
			// 请求文件尝试用HEAD获取下载文件大小，测试服务器是否支持range，并判断下载类型
			// 同时，加载代理列表供后续使用
			// urls := []string{
			// 	"http:// mirrors.ustc.edu.cn/debian/extrafiles",
			// 	"http:// mirrors.ustc.edu.cn/centos/filelist.gz",
			// 	"http:// gzm.com:8080/MacFamilyTree_8.3.6_WaitsUn.com.dmg",
			// 	"http:// gzm.com:8080/MacFamilyTree_8.3.6_WaitsUn.com.dmg",
			// 	"http:// gzm.com:8080/321",
			// 	"http:// gzm.com:8080/README.md",
			// 	"https:// mirrors.163.com/debian-cd/9.8.0/amd64/iso-cd/debian-9.8.0-amd64-netinst.iso",
			// 	"http:// iso.mirrors.ustc.edu.cn/debian-cd/9.8.0/amd64/iso-cd/debian-9.8.0-amd64-netinst.iso",
			// }

			for _, url := range urls {
				go util.Download(url, c)
			}

			// p2
			// 下载模式分为，单线程下载，多线程下载单代理，多线程多代理下载
			for range urls {
				<-c
			}
			// p3
			// 合并文件
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		Logger.Fatal(err)
	}

}
