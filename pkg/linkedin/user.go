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
	ConnectionCount string `json:"connectionCount" csv:"connectionCount"`
	FollowerCount   string `json:"followerCount"   csv:"followerCount"`
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
	extractFollowersCount(strings.TrimSpace(doc.Find(".top-card-layout__first-subline").Text()))
	connectionCount := strings.TrimSpace(doc.Find(".top-card-layout__first-subline span").Eq(0).Text())
	followerCount := strings.TrimSpace(doc.Find(".top-card-layout__first-subline span").Eq(1).Text())
	userTitle := strings.TrimSpace(doc.Find(".top-card-layout__headline").Text())
	location := strings.TrimSpace(doc.Find(".top-card-layout__first-subline div").Text())
	name := strings.TrimSpace(doc.Find(".top-card-layout__title").Text())
	userLink := cleanURL(doc.Find("head link").AttrOr("href", ""))

	user = User{
		ConnectionCount: connectionCount,
		FollowerCount:   followerCount,
		UserTitle:       userTitle,
		Location:        location,
		Name:            name,
		UserLink:        userLink,
	}

	// Print the user for testing
	if debug {
		log.Printf("User: %+v", user)
	}

	return &user, nil
}
