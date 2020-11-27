package api

import (
	"fmt"
	"github.com/DarthPestilane/qq-song-get/request"
	requestMock "github.com/DarthPestilane/qq-song-get/test/mock/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestInfo(t *testing.T) {
	mid := "test-mid"
	t.Run("when type is a song", func(t *testing.T) {
		typ := typeSong
		t.Run("should return the song info", func(t *testing.T) {
			client := &requestMock.Client{}
			client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{}, []byte(`{"code": 0, "data": [{"id": 1}]}`), nil).Once()
			request.DefaultClient = client
			songs, err := Info(typ, mid)
			assert.NoError(t, err)
			assert.Len(t, songs, 1)
		})
		t.Run("should return error when the request failed", func(t *testing.T) {
			client := &requestMock.Client{}
			client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{}, []byte{}, fmt.Errorf("some error")).Once()
			request.DefaultClient = client
			songs, err := Info(typ, mid)
			assert.Error(t, err)
			assert.Nil(t, songs)
		})
		t.Run("should return error when the response code is non-zero", func(t *testing.T) {
			client := &requestMock.Client{}
			client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{}, []byte(`{"code": 1}`), nil).Once()
			request.DefaultClient = client
			songs, err := Info(typ, mid)
			assert.Error(t, err)
			assert.Nil(t, songs)
		})
		t.Run("should return error when the response body is invalid", func(t *testing.T) {
			client := &requestMock.Client{}
			client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{}, []byte(`{"code": 0, "data": {"wrong": true}}`), nil).Once()
			request.DefaultClient = client
			songs, err := Info(typ, mid)
			assert.Error(t, err)
			assert.Nil(t, songs)
		})
	})
	t.Run("when type is an album", func(t *testing.T) {
		typ := typeAlbum
		t.Run("should return the songs info from the album", func(t *testing.T) {
			client := &requestMock.Client{}
			client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{}, []byte(`{"code": 0, "data": {"getSongInfo": [{"id": 1}, {"id": 2}]}}`), nil).Once()
			request.DefaultClient = client
			songs, err := Info(typ, mid)
			assert.NoError(t, err)
			assert.Len(t, songs, 2)
			assert.EqualValues(t, 1, songs[0].ID)
			assert.EqualValues(t, 2, songs[1].ID)
		})
		t.Run("should return error when the request failed", func(t *testing.T) {
			client := &requestMock.Client{}
			client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{}, []byte{}, fmt.Errorf("some error")).Once()
			request.DefaultClient = client
			songs, err := Info(typ, mid)
			assert.Error(t, err)
			assert.Nil(t, songs)
		})
		t.Run("should return error when the response code is non-zero", func(t *testing.T) {
			client := &requestMock.Client{}
			client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{}, []byte(`{"code": 1}`), nil).Once()
			request.DefaultClient = client
			songs, err := Info(typ, mid)
			assert.Error(t, err)
			assert.Nil(t, songs)
		})
		t.Run("should return error when the response body is invalid", func(t *testing.T) {
			client := &requestMock.Client{}
			client.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{}, []byte(`{"code": 0, "data": {"getSongInfo": true}}`), nil).Once()
			request.DefaultClient = client
			songs, err := Info(typ, mid)
			assert.Error(t, err)
			assert.Nil(t, songs)
		})
	})
	t.Run("should return error when type is invalid", func(t *testing.T) {
		songs, err := Info("wrong type", "test")
		assert.Error(t, err)
		assert.Nil(t, songs)
	})
}
