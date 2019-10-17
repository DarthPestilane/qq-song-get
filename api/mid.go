package api

import (
	"fmt"
	"regexp"
)

// FindMid 找到媒体ID
func FindMid(url string) (typ, mid string, err error) {
	re := regexp.MustCompile("/(song|album)/(\\w+)\\.html")
	matched, ok := re.FindStringSubmatch(url), re.MatchString(url)
	if !ok {
		return "", "", fmt.Errorf("invalid qq music address: %s", url)
	}
	return matched[1], matched[2], nil
}
