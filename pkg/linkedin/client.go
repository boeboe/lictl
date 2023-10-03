package linkedin

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ScrapeClient struct {
	client     *http.Client
	rate       time.Duration
	userAgents []string
	proxies    []string
}

func NewScrapeClient(rate time.Duration, userAgents []string, timeout time.Duration) *ScrapeClient {
	if len(userAgents) == 0 {
		log.Fatal("user agents list cannot be empty")
	}

	return &ScrapeClient{
		client:     &http.Client{Timeout: timeout},
		rate:       rate,
		userAgents: userAgents,
	}
}

func (c *ScrapeClient) SetProxies(proxies []string) {
	c.proxies = proxies
}

func (c *ScrapeClient) Do(req *http.Request) (*http.Response, error) {
	c.waitRandomDuration()
	c.setRandomProxy(req)
	req.Header.Set("User-Agent", c.randomUserAgent())
	return c.client.Do(req)
}

func (c *ScrapeClient) randomUserAgent() string {
	randIndex := rand.Intn(len(c.userAgents))
	return c.userAgents[randIndex]
}

func (c *ScrapeClient) DoWithRetry(req *http.Request, maxRetries int) (*http.Response, error) {
	if maxRetries <= 0 {
		return nil, errors.New("maxRetries should be greater than 0")
	}

	var resp *http.Response
	var err error
	for i := 0; i < maxRetries; i++ {
		resp, err = c.Do(req)
		if err == nil {
			return resp, nil
		}
		log.Printf("attempt %d: Error sending request to %s: %v", i+1, req.URL.String(), err)
		time.Sleep(time.Second * time.Duration(i+1)) // exponential backoff
	}
	return nil, fmt.Errorf("failed after %d retries: %v", maxRetries, err)
}

func (c *ScrapeClient) DoDebug(req *http.Request) (*http.Response, error) {
	log.Printf("Sending request to %s\n", req.URL.String())
	c.waitRandomDuration()
	c.setRandomProxy(req)
	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("request to %s failed: %v", req.URL.String(), err)
		return nil, err
	}
	log.Printf("received response with status code: %d", resp.StatusCode)
	return resp, nil
}

func (c *ScrapeClient) waitRandomDuration() {
	randomDuration := c.rate + time.Duration(rand.Int63n(int64(c.rate)))
	log.Printf("waiting for %v before sending request", randomDuration)
	time.Sleep(randomDuration)
}

func (c *ScrapeClient) setRandomProxy(req *http.Request) {
	if len(c.proxies) > 0 {
		randIndex := rand.Intn(len(c.proxies))
		proxyURL, err := url.Parse(c.proxies[randIndex])
		if err != nil {
			log.Printf("failed to parse proxy URL %s: %v", c.proxies[randIndex], err)
			return
		}
		c.client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	} else {
		// Use default transport if no proxies are set
		c.client.Transport = http.DefaultTransport
	}
}
