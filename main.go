package main

import (
	"fmt"
	"github.com/DarthPestilane/qq-song-get/api"
	"github.com/DarthPestilane/qq-song-get/logger"
	"github.com/DarthPestilane/qq-song-get/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

// Some build tags
// Should be passed from: go build -ldflags "-s -w -X main.Version=1.0.0 ..."
var (
	// Version represents the app version.
	Version = "0.0.0"

	// BuildTime represents the time when the app built.
	BuildTime string
)

func main() {
	cmd := &cobra.Command{}
	cmd.Use = "qq-song-get"
	cmd.Long = "./qq-song-get [options] url \n\n example: ./qq-song-get --color=off https://y.qq.com/n/yqq/album/000dilOO3JYIr4.html"
	cmd.Version = fmt.Sprintf("%s; build at %s; build by %s", Version, BuildTime, runtime.Version())
	cmd.Args = cobra.MaximumNArgs(1)
	colorFlag := cmd.Flags().String("color", "on", "Display colorful output. Accept `on` and off")
	debugFlag := cmd.Flags().Bool("debug", false, "Print Set log level to 'debug' to print debug logs")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if colorFlag != nil && *colorFlag == "off" {
			logger.SetFormatter(&logger.Formatter{DisableColor: true})
		}
		if debugFlag != nil && *debugFlag {
			logger.SetLevel(logrus.DebugLevel)
		}
		return proceed(args[0])
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func proceed(url string) error {
	typ, mid, err := util.ExtractTypeAndMid(url)
	if err != nil {
		return fmt.Errorf("find mid failed: %w", err)
	}

	// fetch song info from single song or album
	songs, err := api.Info(typ, mid)
	if err != nil {
		return fmt.Errorf("fetch media info failed: %w", err)
	}

	// prepare to download
	mp3List, err := api.Prepare(songs)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}

	if 0 == len(mp3List) {
		logrus.Warn("没有可下载的音乐")
		return nil
	}

	// download now!
	util.DownloadBatch(mp3List)
	return nil
}
