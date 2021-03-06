package request

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var (
	userAgentList []string
	DefaultClient IClient = &Client{}
)

func init() {
	userAgentList = []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
		"Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 6 Build/LYZ28E) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Mobile/14F89;GameHelper",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1",
		"Mozilla/5.0 (iPad; CPU OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:46.0) Gecko/20100101 Firefox/46.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.4",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:46.0) Gecko/20100101 Firefox/46.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/13.10586",
	}
}

type IClient interface {
	Get(url string, qs map[string]string, shouldPretend bool) (resp *http.Response, body []byte, err error)
	Head(url string, qs map[string]string, shouldPretend bool) (resp *http.Response, body []byte, err error)
}

type Client struct {
}

func (c *Client) Get(url string, qs map[string]string, shouldPretend bool) (*http.Response, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("build request failed: %v", err)
	}
	if len(qs) != 0 {
		query := req.URL.Query()
		for k, v := range qs {
			query.Set(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}
	if shouldPretend {
		pretend(req)
	}
	return sendAndGetBody(req)
}

func (c *Client) Head(url string, qs map[string]string, shouldPretend bool) (*http.Response, []byte, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("build request failed: %v", err)
	}
	if len(qs) != 0 {
		query := req.URL.Query()
		for k, v := range qs {
			query.Set(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}
	if shouldPretend {
		pretend(req)
	}
	return sendAndGetBody(req)
}

func userAgent() string {
	return userAgentList[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(userAgentList))]
}

func pretend(req *http.Request) {
	req.Header.Set("User-Agent", userAgent())
	req.Header.Set("Origin", "https://c.y.qq.com")
	req.Header.Set("Referer", "https://c.y.qq.com")
}

func sendAndGetBody(req *http.Request) (*http.Response, []byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do request failed: %s", err)
	}
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, nil, fmt.Errorf("read response body failed: %v", err)
	}
	return resp, buf.Bytes(), nil
}
