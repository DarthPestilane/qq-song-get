package util

import (
	"fmt"
	"github.com/DarthPestilane/qq-song-get/api"
	"github.com/DarthPestilane/qq-song-get/request"
	requestMock "github.com/DarthPestilane/qq-song-get/test/mock/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"os"
	"testing"
)

func TestDownloadBatch(t *testing.T) {
	t.Run("when mp3 list is empty", func(t *testing.T) {
		t.Run("should return immediately when mp3 list is empty", func(t *testing.T) {
			assert.NoError(t, os.RemoveAll(downloadPath))
			assert.NotPanics(t, func() { DownloadBatch(nil) })
			assert.NoDirExists(t, downloadPath)
		})
	})
	t.Run("when mp3 list is not empty", func(t *testing.T) {
		mp3List := []api.MP3{
			{
				Singer:      "test1-singer",
				Title:       "test1-title",
				DownloadURL: "test1.com",
			},
			{
				Singer:      "test2-singer",
				Title:       "test2-title",
				DownloadURL: "test2.com",
			},
		}
		content := []byte("test")
		t.Run("should save the content to file", func(t *testing.T) {
			assert.NoError(t, os.RemoveAll(downloadPath))
			client := &requestMock.Client{}
			client.On("Head", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{
				Status:        "ok",
				StatusCode:    200,
				ContentLength: int64(len(content)),
			}, []byte{}, nil)
			client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{}, content, nil)
			request.DefaultClient = client
			assert.NotPanics(t, func() { DownloadBatch(mp3List) })
			for _, mp3 := range mp3List {
				assert.FileExists(t, fmt.Sprintf("%s/%s-%s.mp3", downloadPath, mp3.Singer, mp3.Title))
			}
		})
	})
}
