package linkedin

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/boeboe/lictl/pkg/utils"
	"github.com/corpix/uarand"
)

// Job represents the structure of a LinkedIn job.
type Company struct {
	FollowerCount int    `json:"followerCount"`
	FoundedOn     string `json:"foundedOn"`
	Headline      string `json:"headline"`
	Headquarters  string `json:"headquarters"`
	Industry      string `json:"industry"`
	Name          string `json:"name"`
	Size          string `json:"size"`
	Specialties   string `json:"specialties"`
	Type          string `json:"type"`
	Website       string `json:"website"`
}

func GetCompaniesOnline(urls []string, interval time.Duration, debug bool) ([]Company, error) {
	var allCompanies []Company

	for _, url := range urls {
		if debug {
			fmt.Printf("going to fetch comany url %v", url)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept-Encoding", "identity")
		req.Header.Set("User-Agent", uarand.GetRandom())

		company, err := GetCompanyPage(req, debug)
		if err != nil {
			if httpErr, ok := err.(*utils.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				return allCompanies, err // Return the companies fetched so far along with the error
			}
			return nil, err
		}
		allCompanies = append(allCompanies, company)
		time.Sleep(interval)
	}
	return allCompanies, nil
}

func GetCompanyPage(req *http.Request, debug bool) (Company, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Company{}, fmt.Errorf("failed to fetch LinkedIn company: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Company{}, &utils.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("received non-2xx response: %d %s", resp.StatusCode, resp.Status),
		}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Company{}, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var company Company
	followerCount, _ := extractCompanyFollowers(strings.TrimSpace(doc.Find(".top-card-layout__first-subline").Text()))
	foundedOn := strings.TrimSpace(doc.Find("div[data-test-id='about-us__foundedOn'] dd").Text())
	headline := strings.TrimSpace(doc.Find(".top-card-layout__second-subline").Text())
	headquarters := strings.TrimSpace(doc.Find("div[data-test-id='about-us__headquarters'] dd").Text())
	industry := strings.TrimSpace(doc.Find("div[data-test-id='about-us__industry'] dd").Text())
	name := strings.TrimSpace(doc.Find(".top-card-layout__title").Text())
	size := strings.TrimSpace(doc.Find("div[data-test-id='about-us__size'] dd").Text())
	specialties := strings.TrimSpace(doc.Find("div[data-test-id='about-us__specialties'] dd").Text())
	companyType := strings.TrimSpace(doc.Find("div[data-test-id='about-us__organizationType'] dd").Text())
	website := strings.TrimSpace(doc.Find("div[data-test-id='about-us__website'] dd").Text())

	company = Company{
		FollowerCount: followerCount,
		FoundedOn:     foundedOn,
		Headquarters:  headquarters,
		Headline:      headline,
		Industry:      industry,
		Name:          name,
		Size:          size,
		Specialties:   specialties,
		Type:          companyType,
		Website:       website,
	}

	// Print the company for testing
	if debug {
		log.Printf("Company: %+v", company)
	}

	return company, nil
}

func extractCompanyFollowers(s string) (int, error) {
	// Define a regex pattern to match numbers
	re := regexp.MustCompile(`(\d+,?\d*)\s*followers`)
	match := re.FindStringSubmatch(s)

	if len(match) < 2 {
		return 0, fmt.Errorf("no match found")
	}

	// Remove commas and convert to integer
	numStr := strings.ReplaceAll(match[1], ",", "")
	return strconv.Atoi(numStr)
}
