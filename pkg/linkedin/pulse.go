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
	Author               string `json:"author"                csv:"author"`
	AuthorLinkedInUrl    string `json:"authorLinkedInUrl"     csv:"authorLinkedInUrl"`
	AuthorTitle          string `json:"authorTitle"           csv:"authorTitle"`
	CommentCount         int    `json:"commmentCount"         csv:"commmentCount"`
	AuthorFollowingCount int    `json:"authorFollowingCount"  csv:"authorFollowingCount"`
	LikesCount           int    `json:"likesCount"            csv:"likesCount"`
	PublishDate          string `json:"publishDate"           csv:"publishDate"`
	PulseLink            string `json:"pulseLink"             csv:"pulseLink"`
	Title                string `json:"title"                 csv:"title"`
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
	header := doc.Find(".base-main-card--link")
	footer := doc.Find(".main-publisher-card")
	social := doc.Find(".x-social-activity")

	author := strings.TrimSpace(header.Find(".base-card__full-link").Text())
	authorLinkedInUrl := cleanURL(strings.TrimSpace(header.Find(".base-card__full-link").AttrOr("href", "")))
	authorTitle := strings.TrimSpace(header.Find(".base-main-card__subtitle").Text())
	commmentCount, _ := extractCommentsCount(strings.TrimSpace(social.Find("a[data-test-id='social-actions__comments']").Text()))
	authorFollowingCount, _ := extractFollowersCount(strings.TrimSpace(footer.Find(".base-main-card__subtitle").Text()))
	likesCount, _ := extractLikesCount(strings.TrimSpace(social.Find("span[data-test-id='social-actions__reaction-count']").Text()))
	publishDate := strings.TrimSpace(header.Find(".base-main-card__metadata").Text())
	pulseLink := cleanURL(doc.Find("head link").AttrOr("href", ""))
	title := strings.TrimSpace(doc.Find(".pulse-title").Text())

	pulse = Pulse{
		Author:               author,
		AuthorLinkedInUrl:    authorLinkedInUrl,
		AuthorTitle:          authorTitle,
		CommentCount:         commmentCount,
		AuthorFollowingCount: authorFollowingCount,
		LikesCount:           likesCount,
		PublishDate:          publishDate,
		PulseLink:            pulseLink,
		Title:                title,
	}

	// Print the pulse for testing
	if debug {
		log.Printf("Pulse: %+v", pulse)
	}

	return &pulse, nil
}
