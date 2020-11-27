package api

import (
	"fmt"
	"github.com/DarthPestilane/qq-song-get/model"
	"github.com/DarthPestilane/qq-song-get/request"
	requestMock "github.com/DarthPestilane/qq-song-get/test/mock/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestPrepare(t *testing.T) {
	songs := []model.Song{
		{
			Mid:   "song-1",
			Title: "song-1-title",
			Singer: []model.Singer{
				{
					ID:   1,
					Mid:  "singer-1",
					Name: "singer-1-name",
				},
			},
		},
		{
			Mid:   "song-2",
			Title: "song-2-title",
			Singer: []model.Singer{
				{
					ID:   2,
					Mid:  "singer-2",
					Name: "singer-2-name",
				},
			},
		},
	}
	t.Run("should return the mp3 list", func(t *testing.T) {
		client := &requestMock.Client{}
		respBody := []byte(`{"code": 0, "req0": {"data": {"midurlinfo": [{"purl": "test.com"}]}}}`)
		client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: 200}, respBody, nil)
		request.DefaultClient = client
		mp3List, err := Prepare(songs)
		assert.NoError(t, err)
		assert.Len(t, mp3List, 2)
	})
	t.Run("should filter the forbidden songs", func(t *testing.T) {
		client := &requestMock.Client{}
		respBody := []byte(`{"code": 0, "req0": {"data": {"midurlinfo": [{"purl": ""}]}}}`)
		client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: 200}, respBody, nil)
		request.DefaultClient = client
		mp3List, err := Prepare(songs)
		assert.NoError(t, err)
		assert.Len(t, mp3List, 0)
	})
	t.Run("should return error if response code is non-zero", func(t *testing.T) {
		client := &requestMock.Client{}
		respBody := []byte(`{"code": 1, "req0": {"data": {"midurlinfo": [{"purl": ""}]}}}`)
		client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: 200}, respBody, nil)
		request.DefaultClient = client
		mp3List, err := Prepare(songs)
		assert.Error(t, err)
		assert.Nil(t, mp3List)
	})
	t.Run("should return error if response body is invalid", func(t *testing.T) {
		client := &requestMock.Client{}
		respBody := []byte(`{"code": 1, "req0": []}`)
		client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: 200}, respBody, nil)
		request.DefaultClient = client
		mp3List, err := Prepare(songs)
		assert.Error(t, err)
		assert.Nil(t, mp3List)
	})
	t.Run("should return error if request failed", func(t *testing.T) {
		client := &requestMock.Client{}
		respBody := []byte(`{}`)
		client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: 200}, respBody, fmt.Errorf("some error"))
		request.DefaultClient = client
		mp3List, err := Prepare(songs)
		assert.Error(t, err)
		assert.Nil(t, mp3List)
	})
}
