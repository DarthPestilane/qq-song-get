package api

import (
	"fmt"
	"regexp"
)

// FindMid 找到媒体ID
func FindMid(url string) (typ, mid string, err error) {
	re := regexp.MustCompile("/(song|album)/(\\w+)\\.html")
	if !re.MatchString(url) {
		return "", "", fmt.Errorf("invalid qq music address: %s", url)
	}
	matched := re.FindStringSubmatch(url)
	return matched[1], matched[2], nil
}
