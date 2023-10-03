package linkedin

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func cleanURL(link string) string {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return link
	}
	parsedURL.RawQuery = ""
	return parsedURL.String()
}

func extractLikesCount(s string) (int, error) {
	numStr := strings.ReplaceAll(s, ",", "")
	return strconv.Atoi(numStr)
}

func extractCommentsCount(s string) (int, error) {
	re := regexp.MustCompile(`(\d+,?\d*)\s*Comments`)
	match := re.FindStringSubmatch(s)
	if len(match) < 2 {
		return 0, fmt.Errorf("no match found")
	}
	numStr := strings.ReplaceAll(match[1], ",", "")
	return strconv.Atoi(numStr)
}

func extractFollowersCount(s string) (int, error) {
	re := regexp.MustCompile(`(\d+,?\d*)\s*followers`)
	match := re.FindStringSubmatch(s)
	if len(match) < 2 {
		return 0, fmt.Errorf("no match found")
	}
	numStr := strings.ReplaceAll(match[1], ",", "")
	return strconv.Atoi(numStr)
}
