package api

import (
	"encoding/json"
	"fmt"
	"github.com/DarthPestilane/qq-song-get/model"
	"github.com/DarthPestilane/qq-song-get/request"
)

const (
	typeSong  = "song"
	typeAlbum = "album"

	songInfoURL  = "https://c.y.qq.com/v8/fcg-bin/fcg_play_single_song.fcg?platform=yqq&format=json"
	albumInfoURL = "https://c.y.qq.com/v8/fcg-bin/fcg_v8_album_detail_cp.fcg?newsong=1&platform=yqq&format=json"
)

type (
	// SongInfoResponse 歌曲信息响应
	SongInfoResponse struct {
		Code int          `json:"code"`
		Data []model.Song `json:"data"`
	}

	// AlbumResponse 专辑信息响应
	AlbumResponse struct {
		Code int `json:"code"`
		Data struct {
			GetSongInfo []model.Song `json:"getSongInfo"`
		} `json:"data"`
	}
)

// Info 根据类型和媒体ID获取信息
func Info(typ, mid string) ([]model.Song, error) {
	switch typ {
	case typeSong:
		song, err := infoSingleSong(mid)
		if err != nil {
			return nil, err
		}
		return []model.Song{*song}, nil
	case typeAlbum:
		return infoAlbum(mid)
	}
	return nil, fmt.Errorf("invalid type: %s", typ)
}

func infoSingleSong(mid string) (*model.Song, error) {
	_, respBody, err := request.DefaultClient.Get(songInfoURL, map[string]string{"songmid": mid}, true)
	if err != nil {
		return nil, fmt.Errorf("request for song info failed: %v", err)
	}
	var songInfoResp SongInfoResponse
	if err := json.Unmarshal(respBody, &songInfoResp); err != nil {
		return nil, fmt.Errorf("parse song info failed: %v", err)
	}
	if songInfoResp.Code != 0 {
		return nil, fmt.Errorf("request failed, code: %d", songInfoResp.Code)
	}
	return &songInfoResp.Data[0], nil
}

func infoAlbum(mid string) ([]model.Song, error) {
	_, respBody, err := request.DefaultClient.Get(albumInfoURL, map[string]string{"albummid": mid}, true)
	if err != nil {
		return nil, fmt.Errorf("request for song info failed: %v", err)
	}
	var albumResp AlbumResponse
	if err := json.Unmarshal(respBody, &albumResp); err != nil {
		return nil, fmt.Errorf("parse song info failed: %v", err)
	}
	if albumResp.Code != 0 {
		return nil, fmt.Errorf("request failed, code: %d", albumResp.Code)
	}
	return albumResp.Data.GetSongInfo, nil
}
