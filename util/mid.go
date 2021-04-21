package util

import (
	"fmt"
	"regexp"
)

// ExtractTypeAndMid 找到类型媒体ID
func ExtractTypeAndMid(url string) (typ, mid string, err error) {
	re := regexp.MustCompile(`/(song|album)/(\w+)\.html`)
	if !re.MatchString(url) {
		return "", "", fmt.Errorf("invalid qq music address: %s", url)
	}
	matched := re.FindStringSubmatch(url)
	return matched[1], matched[2], nil
}
