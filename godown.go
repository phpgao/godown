package main

import (
	"github.com/phpgao/godown/downloader"
	"github.com/phpgao/godown/proxy"
	"github.com/phpgao/godown/util"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
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
	proxy.Logger = Logger
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
			Name:  "dir",
			Value: ".",
			Usage: "download dir",
		},
		cli.StringFlag{
			Name:  "p, proxy",
			Usage: "proxys",
		},
		cli.StringFlag{
			Name:  "md5",
			Value: "",
			Usage: "md5 check, only support single file download",
		},
		cli.StringFlag{
			Name:  "sha1",
			Value: "",
			Usage: "sha1 check, only support single file download",
		},
		cli.StringFlag{
			Name:  "sha2",
			Value: "",
			Usage: "sha2 check, only support single file download",
		},
		cli.Int64Flag{
			Name:  "c, concurrency",
			Value: 10,
			Usage: "concurrency you know",
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		var urls []string
		if c.NArg() > 0 {
			urls = c.Args()
		}

		if num := len(urls); num > 0 {
			// First,set the logger
			setLogger(c)
			Logger.Debug(urls)
			var filePath string
			var dir string
			check := make(map[string]string)

			if num > 1 {
				// Multi file need a dir
				dir = c.String("dir")
				dir, err = filepath.Abs(dir)
				if err != nil {
					Logger.Error(err)
					return err
				}
			} else {
				// Just a file name and parent dir exists
				filePath = c.String("out")
				filePath, err := filepath.Abs(filePath)
				if err != nil {
					Logger.Error(err)
					return err
				}
				dir = filepath.Dir(filePath)
				filePath = filepath.Base(filePath)

				check["md5"] = c.String("md5")
				check["sha1"] = c.String("sha1")
				check["sha2"] = c.String("sha2")
			}

			_, err = util.CheckDir(dir)
			if err != nil {
				Logger.Errorf("Error when check %s", dir)
			}
			Logger.Debug(dir, filePath)
			workingChan := make(chan error)
			// p1
			// 请求文件尝试用HEAD获取下载文件大小，测试服务器是否支持range，并判断下载类型
			// 同时，加载代理列表供后续使用
			// urls := []string{
			// 	"http://mirrors.ustc.edu.cn/debian/extrafiles",
			// 	"http://mirrors.ustc.edu.cn/centos/filelist.gz",
			// 	"http://gzm.com:8080/MacFamilyTree_8.3.6_WaitsUn.com.dmg",
			// 	"http://gzm.com:8080/MacFamilyTree_8.3.6_WaitsUn.com.dmg",
			// 	"http://gzm.com:8080/321",
			// 	"http://gzm.com:8080/README.md",
			// 	"https:// mirrors.163.com/debian-cd/9.8.0/amd64/iso-cd/debian-9.8.0-amd64-netinst.iso",
			// 	"http:// iso.mirrors.ustc.edu.cn/debian-cd/9.8.0/amd64/iso-cd/debian-9.8.0-amd64-netinst.iso",
			// }

			// l means limit in downloader
			l := c.Int64("concurrency")
			for _, url := range urls {
				go util.Download(url, dir, filePath, workingChan, l, check)
			}

			// p2
			// 下载模式分为，单线程下载，多线程下载单代理，多线程多代理下载
			for range urls {
				<-workingChan
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
