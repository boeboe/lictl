package linkedin

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
)

// Pulse represents the structure of a LinkedIn pulse.
type Pulse struct {
	Author       int    `json:"author"         csv:"author"`
	AuthorTitle  string `json:"authorTitle"    csv:"authorTitle"`
	CommentCount int    `json:"commmentCount"  csv:"commmentCount"`
	HashTags     string `json:"hashTags"       csv:"hashTags"`
	LikesCount   int    `json:"likesCount"     csv:"likesCount"`
	PublishDate  string `json:"publishDate"    csv:"publishDate"`
	PulseLink    string `json:"pulseLink"      csv:"pulseLink"`
	Title        string `json:"title"          csv:"title"`
}

func (p *Pulse) CsvContent() string {
	if p == nil {
		return ""
	}
	return CsvContent(p)
}

func (p *Pulse) CsvHeader() string {
	if p == nil {
		return ""
	}
	return CsvHeader(p)
}

func (p *Pulse) Json() string {
	if p == nil {
		return ""
	}
	return Json(p)
}

type Pulses []*Pulse

func (ps Pulses) Len() int {
	return len(ps)
}

func (ps Pulses) Get(i int) Serializable {
	return Serializable(ps[i])
}

func SearchPulsesOnline(keywords []string, interval time.Duration, debug bool) (Pulses, error) {
	var pulses []*Pulse
	urls, err := GoogleGetLinkedInPulseURLs(keywords, interval, debug)
	if err != nil {
		return nil, fmt.Errorf("error fetching LinkedIn pulse URLs: %v", err)
	}

	var errs []string
	for _, url := range urls {
		pulse, err := GetPulseFromUrl(url, debug)
		if err != nil {
			errs = append(errs, fmt.Sprintf("error fetching pulse from URL %s: %v", url, err))
			continue
		}
		pulses = append(pulses, pulse)
	}

	if len(errs) > 0 {
		return pulses, fmt.Errorf("encountered errors: %s", strings.Join(errs, "; "))
	}

	return pulses, nil
}

func GetPulseFromUrl(url string, debug bool) (*Pulse, error) {
	if debug {
		fmt.Printf("going to fetch pulse from url %v", url)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &Pulse{}, err
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("User-Agent", uarand.GetRandom())

	pulse, err := getPulseFromRequest(req, debug)
	if err != nil {
		if httpErr, ok := err.(*HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
			return &Pulse{}, err
		}
		return &Pulse{}, err
	}
	return pulse, nil
}

func getPulseFromRequest(req *http.Request, debug bool) (*Pulse, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &Pulse{}, fmt.Errorf("failed to fetch LinkedIn pulse: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &Pulse{}, &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("received non-2xx response: %d %s", resp.StatusCode, resp.Status),
		}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return &Pulse{}, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var pulse Pulse
	extractCompanyFollowers(strings.TrimSpace(doc.Find(".top-card-layout__first-subline").Text()))
	// followerCount, _ := extractCompanyFollowers(strings.TrimSpace(doc.Find(".top-card-layout__first-subline").Text()))
	// foundedOn := strings.TrimSpace(doc.Find("div[data-test-id='about-us__foundedOn'] dd").Text())
	// headline := strings.TrimSpace(doc.Find(".top-card-layout__second-subline").Text())
	// headquarters := strings.TrimSpace(doc.Find("div[data-test-id='about-us__headquarters'] dd").Text())
	// industry := strings.TrimSpace(doc.Find("div[data-test-id='about-us__industry'] dd").Text())
	// name := strings.TrimSpace(doc.Find(".top-card-layout__title").Text())
	// size := strings.TrimSpace(doc.Find("div[data-test-id='about-us__size'] dd").Text())
	// specialties := strings.TrimSpace(doc.Find("div[data-test-id='about-us__specialties'] dd").Text())
	// companyType := strings.TrimSpace(doc.Find("div[data-test-id='about-us__organizationType'] dd").Text())
	// website := strings.TrimSpace(doc.Find("div[data-test-id='about-us__website'] dd").Text())

	pulse = Pulse{
		// FollowerCount: followerCount,
		// FoundedOn:     foundedOn,
		// Headquarters:  headquarters,
		// Headline:      headline,
		// Industry:      industry,
		// Name:          name,
		// Size:          size,
		// Specialties:   specialties,
		// Type:          companyType,
		// Website:       website,
	}

	// Print the pulse for testing
	if debug {
		log.Printf("Pulse: %+v", pulse)
	}

	return &pulse, nil
}
