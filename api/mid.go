package api

import (
	"fmt"
	"regexp"
)

// FindTypeAndMid 找到类型媒体ID
func FindTypeAndMid(url string) (typ, mid string, err error) {
	re := regexp.MustCompile("/(song|album)/(\\w+)\\.html")
	if !re.MatchString(url) {
		return "", "", fmt.Errorf("invalid qq music address: %s", url)
	}
	matched := re.FindStringSubmatch(url)
	return matched[1], matched[2], nil
}
