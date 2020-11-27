package request

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type Client struct {
	mock.Mock
}

func (m *Client) Get(url string, qs map[string]string, shouldPretend bool) (resp *http.Response, body []byte, err error) {
	args := m.Called(url, qs, shouldPretend)
	return args.Get(0).(*http.Response), args.Get(1).([]byte), args.Error(2)
}

func (m *Client) Head(url string, qs map[string]string, shouldPretend bool) (resp *http.Response, body []byte, err error) {
	args := m.Called(url, qs, shouldPretend)
	return args.Get(0).(*http.Response), args.Get(1).([]byte), args.Error(2)
}
