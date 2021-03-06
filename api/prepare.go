package api

import (
	"encoding/json"
	"fmt"
	"github.com/DarthPestilane/qq-song-get/model"
	"github.com/DarthPestilane/qq-song-get/request"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

const (
	songURL     = "https://u.y.qq.com/cgi-bin/musicu.fcg"
	downloadURL = "http://ws.stream.qqmusic.qq.com/%s"

	reqModule    = "vkey.GetVkeyServer"
	reqMethod    = "CgiGetVkey"
	reqGuid      = "1"
	reqLoginFlag = 1
	reqUin       = "0"
	reqPlatform  = "20"
)

type (
	// SongURLResponse 请求歌曲下载要用的URL的响应结构体
	SongURLResponse struct {
		Code int `json:"code"`
		Req0 struct {
			Data struct {
				MidURLInfo []struct {
					Purl string `json:"purl"`
				} `json:"midurlinfo"`
			} `json:"data"`
		} `json:"req0"`
	}

	// MP3 mp3结构体，包含下载链接、文件名
	MP3 struct {
		Singer      string
		Title       string
		DownloadURL string
	}

	respWrap struct {
		resp     *http.Response
		respBody []byte
		err      error
	}
)

// Prepare 为下载做准备，组装下载URL
func Prepare(songs []model.Song) ([]MP3, error) {
	respMap := sync.Map{}
	wg := sync.WaitGroup{}
	wg.Add(len(songs))

	// batch send request
	for _, song := range songs {
		param := map[string]interface{}{
			"req0": map[string]interface{}{
				"module": reqModule,
				"method": reqMethod,
				"param": map[string]interface{}{
					"guid":      reqGuid,
					"loginflag": reqLoginFlag,
					"songmid":   []string{song.Mid},
					"uin":       reqUin,
					"platform":  reqPlatform,
				},
			},
		}
		enc, _ := json.Marshal(param)
		go func(song model.Song) {
			defer wg.Done()
			resp, respBody, err := request.DefaultClient.Get(songURL, map[string]string{"data": string(enc)}, true)
			respMap.Store(song.Mid, &respWrap{
				resp:     resp,
				respBody: respBody,
				err:      err,
			})
		}(song)
	}
	wg.Wait()

	mp3List := make([]MP3, 0, len(songs))
	for _, song := range songs {
		value, _ := respMap.Load(song.Mid)
		wrap := value.(*respWrap)
		if err := wrap.err; err != nil {
			return nil, fmt.Errorf("request song url failed: %v", err)
		}
		var songURLResp SongURLResponse
		if err := json.Unmarshal(wrap.respBody, &songURLResp); err != nil {
			return nil, fmt.Errorf("parse song url response failed: %v", err)
		}
		if songURLResp.Code != 0 {
			return nil, fmt.Errorf("request failed, code: %d", songURLResp.Code)
		}
		purl := songURLResp.Req0.Data.MidURLInfo[0].Purl
		if purl == "" {
			logrus.Errorf("%s 需要VIP或仅支持客户端播放", song.Title)
			continue
		}
		mp3 := MP3{
			Title:       song.Title,
			Singer:      song.Singer[0].Name,
			DownloadURL: fmt.Sprintf(downloadURL, purl),
		}
		mp3List = append(mp3List, mp3)
	}
	return mp3List, nil
}
