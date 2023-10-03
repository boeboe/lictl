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
	"github.com/corpix/uarand"
)

// Company represents the structure of a LinkedIn company.
type Company struct {
	FollowerCount int    `json:"followerCount"  csv:"followerCount"`
	FoundedOn     string `json:"foundedOn"      csv:"foundedOn"`
	Headline      string `json:"headline"       csv:"headline"`
	Headquarters  string `json:"headquarters"   csv:"headquarters"`
	Industry      string `json:"industry"       csv:"industry"`
	Name          string `json:"name"           csv:"name"`
	Size          string `json:"size"           csv:"size"`
	Specialties   string `json:"specialties"    csv:"specialties"`
	Type          string `json:"type"           csv:"type"`
	Website       string `json:"website"        csv:"website"`
}

func (c *Company) CsvContent() string {
	if c == nil {
		return ""
	}
	return CsvContent(c)
}

func (c *Company) CsvHeader() string {
	if c == nil {
		return ""
	}
	return CsvHeader(c)
}

func (c *Company) Json() string {
	if c == nil {
		return ""
	}
	return Json(c)
}

type Companies []*Company

func (cs Companies) Len() int {
	return len(cs)
}

func (cs Companies) Get(i int) Serializable {
	return Serializable(cs[i])
}

func SearchCompaniesOnline(keywords []string, interval time.Duration, debug bool) (Companies, error) {
	var companies Companies
	urls, err := GoogleGetLinkedInCompanyURLs(keywords, interval, debug)
	if err != nil {
		return nil, fmt.Errorf("error fetching LinkedIn company URLs: %v", err)
	}

	var errs []string
	for _, url := range urls {
		company, err := GetCompanyFromUrl(url, debug)
		if err != nil {
			errs = append(errs, fmt.Sprintf("error fetching company from URL %s: %v", url, err))
			continue
		}
		companies = append(companies, company)
	}

	if len(errs) > 0 {
		return companies, fmt.Errorf("encountered errors: %s", strings.Join(errs, "; "))
	}

	return companies, nil
}

func GetCompanyFromUrl(url string, debug bool) (*Company, error) {
	if debug {
		fmt.Printf("going to fetch company from url %v", url)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &Company{}, err
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("User-Agent", uarand.GetRandom())

	company, err := getCompanyFromRequest(req, debug)
	if err != nil {
		if httpErr, ok := err.(*HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
			return &Company{}, err
		}
		return &Company{}, err
	}
	return company, nil
}

func getCompanyFromRequest(req *http.Request, debug bool) (*Company, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &Company{}, fmt.Errorf("failed to fetch LinkedIn company: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &Company{}, &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("received non-2xx response: %d %s", resp.StatusCode, resp.Status),
		}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return &Company{}, fmt.Errorf("failed to parse HTML: %w", err)
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
	website := strings.Split(strings.TrimSpace(doc.Find("div[data-test-id='about-us__website'] dd").Text()), "\n")[0]

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

	return &company, nil
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
