package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindTypeAndMid(t *testing.T) {
	t.Run("when url is invalid", func(t *testing.T) {
		url := "http://invalid.com"
		t.Run("it should return error", func(t *testing.T) {
			typ, mid, err := ExtractTypeAndMid(url)
			assert.NotNil(t, err)
			assert.Zero(t, typ)
			assert.Zero(t, mid)
		})
	})
	t.Run("when url is a single song", func(t *testing.T) {
		expectType := "song"
		expectMid := "test"
		url := fmt.Sprintf("https://y.qq.com/n/yqq/%s/%s.html", expectType, expectMid)
		t.Run("it should return without error", func(t *testing.T) {
			typ, mid, err := ExtractTypeAndMid(url)
			assert.NoError(t, err)
			assert.Equal(t, expectType, typ)
			assert.Equal(t, expectMid, mid)
		})
	})
	t.Run("when url is an album", func(t *testing.T) {
		expectType := "album"
		expectMid := "test"
		url := fmt.Sprintf("https://y.qq.com/n/yqq/%s/%s.html", expectType, expectMid)
		t.Run("it should return without error", func(t *testing.T) {
			typ, mid, err := ExtractTypeAndMid(url)
			assert.NoError(t, err)
			assert.Equal(t, expectType, typ)
			assert.Equal(t, expectMid, mid)
		})
	})
}
