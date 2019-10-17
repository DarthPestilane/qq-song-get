package util

import (
	"github.com/cheggaaa/pb"
	"github.com/darthpesitlane/qq-song-get/api"
	"github.com/darthpesitlane/qq-song-get/request"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"sync"
)

// DownloadBatch 批量下载
func DownloadBatch(mp3List []api.MP3) {
	logrus.Info("begin download.")

	wg := sync.WaitGroup{}
	wg.Add(len(mp3List))

	bars := make([]*pb.ProgressBar, 0, len(mp3List))
	for range mp3List {
		bars = append(bars, pb.New(0))
	}
	pool := pb.NewPool(bars...)
	if err := pool.Start(); err != nil {
		logrus.Fatalf("start progress_bar pool failed: %v", err)
	}

	for i, mp3 := range mp3List {
		i, mp3 := i, mp3
		go func() {
			defer wg.Done()
			fileName := adjustFileName(mp3.Singer + "-" + mp3.Title + ".mp3")
			bar := bars[i].Prefix(fileName + ":").SetUnits(pb.U_BYTES)
			bar.Start()
			defer bar.Finish()
			headResp, err := request.HEAD(mp3.DownloadURL, nil, false)
			if err != nil {
				logrus.Errorf("request for %s failed: %v", fileName, err)
				return
			}
			if headResp.StatusCode > 299 || headResp.StatusCode < 200 {
				logrus.Errorf("request for '%s' failed: http status: %s", fileName, headResp.Status)
				return
			}
			bar.SetTotal64(headResp.ContentLength)
			resp, err := request.GET(mp3.DownloadURL, nil, false)
			if err != nil {
				logrus.Errorf("request for %s failed: %v", fileName, err)
				return
			}
			body := bar.NewProxyReader(resp.Body)
			if err := os.MkdirAll("downloads", os.ModePerm); err != nil {
				logrus.Fatalf("mkdir failed: %v", err)
			}
			file, err := os.Create("downloads/" + fileName)
			if err != nil {
				logrus.Fatalf("create file failed: %v", err)
			}
			if _, err := io.Copy(file, body); err != nil {
				logrus.Fatalf("copy failed: %v", err)
			}
		}()
	}
	wg.Wait()
	if err := pool.Stop(); err != nil {
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
