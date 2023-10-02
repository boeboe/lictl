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

// Post represents the structure of a LinkedIn post.
type Post struct {
	Author       int    `json:"author"         csv:"author"`
	AuthorTitle  string `json:"authorTitle"    csv:"authorTitle"`
	CommentCount int    `json:"commmentCount"  csv:"commmentCount"`
	HashTags     string `json:"hashTags"       csv:"hashTags"`
	LikesCount   int    `json:"likesCount"     csv:"likesCount"`
	PostLink     string `json:"postLink"       csv:"postLink"`
	PublishDate  string `json:"publishDate"    csv:"publishDate"`
}

func (p *Post) CsvContent() string {
	if p == nil {
		return ""
	}
	return CsvContent(p)
}

func (p *Post) CsvHeader() string {
	if p == nil {
		return ""
	}
	return CsvHeader(p)
}

func (p *Post) Json() string {
	if p == nil {
		return ""
	}
	return Json(p)
}

type Posts []*Post

func (ps Posts) Len() int {
	return len(ps)
}

func (ps Posts) Get(i int) Serializable {
	return Serializable(ps[i])
}

func SearchPostsOnline(keywords []string, interval time.Duration, debug bool) (Posts, error) {
	var posts []*Post
	urls, err := GoogleGetLinkedInPostURLs(keywords, interval, debug)
	if err != nil {
		return nil, fmt.Errorf("error fetching LinkedIn post URLs: %v", err)
	}

	var errs []string
	for _, url := range urls {
		post, err := GetPostFromUrl(url, debug)
		if err != nil {
			errs = append(errs, fmt.Sprintf("error fetching post from URL %s: %v", url, err))
			continue
		}
		posts = append(posts, post)
	}

	if len(errs) > 0 {
		return posts, fmt.Errorf("encountered errors: %s", strings.Join(errs, "; "))
	}

	return posts, nil
}

func GetPostFromUrl(url string, debug bool) (*Post, error) {
	if debug {
		fmt.Printf("going to fetch post from url %v", url)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &Post{}, err
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("User-Agent", uarand.GetRandom())

	post, err := getPostFromRequest(req, debug)
	if err != nil {
		if httpErr, ok := err.(*HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
			return &Post{}, err
		}
		return &Post{}, err
	}
	return post, nil
}

func getPostFromRequest(req *http.Request, debug bool) (*Post, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &Post{}, fmt.Errorf("failed to fetch LinkedIn post: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &Post{}, &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("received non-2xx response: %d %s", resp.StatusCode, resp.Status),
		}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return &Post{}, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var post Post
	extractCompanyFollowers(strings.TrimSpace(doc.Find(".top-card-layout__first-subline").Text()))
	// followerCount, _ := extractPostFollowers(strings.TrimSpace(doc.Find(".top-card-layout__first-subline").Text()))
	// foundedOn := strings.TrimSpace(doc.Find("div[data-test-id='about-us__foundedOn'] dd").Text())
	// headline := strings.TrimSpace(doc.Find(".top-card-layout__second-subline").Text())
	// headquarters := strings.TrimSpace(doc.Find("div[data-test-id='about-us__headquarters'] dd").Text())
	// industry := strings.TrimSpace(doc.Find("div[data-test-id='about-us__industry'] dd").Text())
	// name := strings.TrimSpace(doc.Find(".top-card-layout__title").Text())
	// size := strings.TrimSpace(doc.Find("div[data-test-id='about-us__size'] dd").Text())
	// specialties := strings.TrimSpace(doc.Find("div[data-test-id='about-us__specialties'] dd").Text())
	// postType := strings.TrimSpace(doc.Find("div[data-test-id='about-us__organizationType'] dd").Text())
	// website := strings.TrimSpace(doc.Find("div[data-test-id='about-us__website'] dd").Text())

	post = Post{
		// FollowerCount: followerCount,
		// FoundedOn:     foundedOn,
		// Headquarters:  headquarters,
		// Headline:      headline,
		// Industry:      industry,
		// Name:          name,
		// Size:          size,
		// Specialties:   specialties,
		// Type:          postType,
		// Website:       website,
	}

	// Print the post for testing
	if debug {
		log.Printf("Post: %+v", post)
	}

	return &post, nil
}
