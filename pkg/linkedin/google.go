package linkedin

import (
	"context"
	"strings"
	"time"

	"github.com/corpix/uarand"
	googlesearch "github.com/rocketlaunchr/google-search"
)

const (
	googleLinkedInCompanyPrefix = "site:linkedin.com/company"
	googleLinkedInPostPrefix    = "site:linkedin.com/posts"
	googleLinkedInPulsePrefix   = "site:linkedin.com/pulse"
	googleLinkedInUserPrefix    = "site:linkedin.com/in"
)

func GoogleGetLinkedInCompanyURLs(keywords []string, interval time.Duration, debug bool) ([]string, error) {
	return googleSearch(googleLinkedInCompanyPrefix, keywords, interval, debug)
}

func GoogleGetLinkedInPostURLs(keywords []string, interval time.Duration, debug bool) ([]string, error) {
	return googleSearch(googleLinkedInPostPrefix, keywords, interval, debug)
}

func GoogleGetLinkedInPulseURLs(keywords []string, interval time.Duration, debug bool) ([]string, error) {
	return googleSearch(googleLinkedInPulsePrefix, keywords, interval, debug)
}

func GoogleGetLinkedInUserURLs(keywords []string, interval time.Duration, debug bool) ([]string, error) {
	return googleSearch(googleLinkedInUserPrefix, keywords, interval, debug)
}

func googleSearch(prefix string, keywords []string, interval time.Duration, debug bool) ([]string, error) {
	query := prefix + " " + strings.Join(keywords, " ")
	opts := googlesearch.SearchOptions{
		Limit:          100,
		UserAgent:      uarand.GetRandom(),
		FollowNextPage: true,
	}
	results, err := googlesearch.Search(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, result := range results {
		urls = append(urls, result.URL)
	}
	return urls, nil
}
