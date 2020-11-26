package util

import (
	"bytes"
	"github.com/DarthPestilane/qq-song-get/api"
	"github.com/DarthPestilane/qq-song-get/request"
	"github.com/cheggaaa/pb"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"sync"
)

const downloadPath = "downloads"

// DownloadBatch 批量下载
func DownloadBatch(mp3List []api.MP3) {
	logrus.Info("begin download.")
	logrus.Info("press Ctrl+C twice to exit.")

	// initiate progress bars
	bars := make([]*pb.ProgressBar, 0, len(mp3List))
	for range mp3List {
		bars = append(bars, pb.New(0))
	}
	progressPool := pb.NewPool(bars...)
	if err := progressPool.Start(); err != nil {
		logrus.Fatalf("start progress_bar pool failed: %v", err)
	}

	// create download directory first
	if err := os.MkdirAll(downloadPath, os.ModePerm); err != nil {
		logrus.Fatalf("mkdir failed: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(mp3List))
	for i, mp3 := range mp3List {
		i, mp3 := i, mp3
		go func() {
			defer wg.Done()

			// prepare single progress bar
			fileName := adjustFileName(mp3.Singer + "-" + mp3.Title + ".mp3")
			bar := bars[i].Prefix(fileName + ":").SetUnits(pb.U_BYTES) // locate individual progress bar
			bar.Start()                                                // progress bar start
			defer bar.Finish()                                         // trigger bar.Finish() before wg.Done()

			// send HEAD request only for content-length
			headResp, _, err := request.DefaultClient.Head(mp3.DownloadURL, nil, false)
			if err != nil {
				logrus.Errorf("request for %s failed: %v", fileName, err)
				return
			}
			if headResp.StatusCode > 299 || headResp.StatusCode < 200 {
				logrus.Errorf("request for '%s' failed: http status: %s", fileName, headResp.Status)
				return
			}
			bar.SetTotal64(headResp.ContentLength) // set progress bar's max length

			// now download mp3!
			_, respBody, err := request.DefaultClient.Get(mp3.DownloadURL, nil, false)
			if err != nil {
				logrus.Errorf("request for %s failed: %v", fileName, err)
				return
			}

			// handle response body and progress bar
			file, err := os.Create(downloadPath + string(os.PathSeparator) + fileName)
			if err != nil {
				logrus.Fatalf("create file failed: %v", err)
			}
			if _, err := io.Copy(file, bar.NewProxyReader(bytes.NewReader(respBody))); err != nil {
				logrus.Fatalf("copy downloaded content failed: %v", err)
			}
		}()
	}
	wg.Wait()

	if err := progressPool.Stop(); err != nil {
		logrus.Fatalf("stop progress_bar pool failed: %v", err)
	}
}

func adjustFileName(fileName string) string {
	replaceMap := map[string]string{
		"/": "",
		" ": ".",
		"'": "",
		`"`: "",
	}
	str := strings.TrimSpace(fileName)
	for find, toReplace := range replaceMap {
		str = strings.ReplaceAll(str, find, toReplace)
	}
	return str
}
