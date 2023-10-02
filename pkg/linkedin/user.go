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

// User represents the structure of a LinkedIn user.
type User struct {
	ConnectionCount int    `json:"connectionCount" csv:"connectionCount"`
	FollowerCount   int    `json:"followerCount"   csv:"followerCount"`
	UserTitle       string `json:"userTitle"       csv:"userTitle"`
	Location        string `json:"location"        csv:"location"`
	Name            string `json:"name"            csv:"name"`
	UserLink        string `json:"userLink"        csv:"userLink"`
}

func (u *User) CsvContent() string {
	if u == nil {
		return ""
	}
	return CsvContent(u)
}

func (u *User) CsvHeader() string {
	if u == nil {
		return ""
	}
	return CsvHeader(u)
}

func (u *User) Json() string {
	if u == nil {
		return ""
	}
	return Json(u)
}

type Users []*User

func (us Users) Len() int {
	return len(us)
}

func (us Users) Get(i int) Serializable {
	return Serializable(us[i])
}

func SearchUsersOnline(keywords []string, interval time.Duration, debug bool) (Users, error) {
	var users []*User
	urls, err := GoogleGetLinkedInUserURLs(keywords, interval, debug)
	if err != nil {
		return nil, fmt.Errorf("error fetching LinkedIn user URLs: %v", err)
	}

	var errs []string
	for _, url := range urls {
		user, err := GetUserFromUrl(url, debug)
		if err != nil {
			errs = append(errs, fmt.Sprintf("error fetching user from URL %s: %v", url, err))
			continue
		}
		users = append(users, user)
	}

	if len(errs) > 0 {
		return users, fmt.Errorf("encountered errors: %s", strings.Join(errs, "; "))
	}

	return users, nil
}

func GetUserFromUrl(url string, debug bool) (*User, error) {
	if debug {
		fmt.Printf("going to fetch user from url %v", url)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &User{}, err
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("User-Agent", uarand.GetRandom())

	user, err := getUserFromRequest(req, debug)
	if err != nil {
		if httpErr, ok := err.(*HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
			return &User{}, err
		}
		return &User{}, err
	}
	return user, nil
}

func getUserFromRequest(req *http.Request, debug bool) (*User, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &User{}, fmt.Errorf("failed to fetch LinkedIn user: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &User{}, &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("received non-2xx response: %d %s", resp.StatusCode, resp.Status),
		}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return &User{}, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var user User
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

	user = User{
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

	// Print the user for testing
	if debug {
		log.Printf("User: %+v", user)
	}

	return &user, nil
}
