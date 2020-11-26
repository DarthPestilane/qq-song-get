package main

import (
	"fmt"
	"github.com/DarthPestilane/qq-song-get/api"
	"github.com/DarthPestilane/qq-song-get/logger"
	"github.com/DarthPestilane/qq-song-get/util"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"runtime"
)

// Some build tags
// Should be passed from: go build -ldflags "-s -w -X main.Version=1.0.0 ..."
var (
	// Version represents the app version.
	Version string = "0.0.0"

	// BuildTime represents the time when the app built.
	BuildTime string
)

func main() {
	app := cli.NewApp()
	app.Name = "qq song download"
	app.Version = fmt.Sprintf("%s; build at %s; build by %s", Version, BuildTime, runtime.Version())
	app.UsageText = "./qq-song-get [options] url\n   example: ./qq-song-get --color off https://y.qq.com/n/yqq/album/000dilOO3JYIr4.html"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "color",
			Value: "on",
			Usage: "Display color output. Accept `on` and off",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		// show help txt when there are no args
		if ctx.NArg() == 0 {
			return cli.ShowAppHelp(ctx)
		}

		// determine color output
		if ctx.String("color") == "off" {
			logger.SetFormatter(&logger.Formatter{DisableColor: true})
		}

		// here we start!
		url := ctx.Args().Get(0)
		if url == "" {
			logrus.Fatalf("url is required")
		}
		typ, mid, err := api.FindMid(url)
		if err != nil {
			logrus.Fatalf("find mid failed: %v", err)
		}

		// fetch song info from single song or album
		songs, err := api.Info(typ, mid)
		if err != nil {
			logrus.Fatalf("fetch media info failed: %v", err)
		}

		// prepare to download
		mp3List, err := api.Prepare(songs)
		if err != nil {
			logrus.Fatalf("prepare failed: %v", err)
		}

		if 0 == len(mp3List) {
			logrus.Warn("没有可下载的音乐")
			return nil
		}

		// download now!
		util.DownloadBatch(mp3List)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("start app failed: %v", err)
	}
}
