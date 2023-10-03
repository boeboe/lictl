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

// Post represents the structure of a LinkedIn post.
type Post struct {
	ActivityURN          string `json:"activityURN"            csv:"activityURN"`
	Author               string `json:"author"                 csv:"author"`
	AuthorLinkedInUrl    string `json:"authorLinkedInUrl"      csv:"authorLinkedInUrl"`
	AuthorTitle          string `json:"authorTitle"            csv:"authorTitle"`
	CommentCount         int    `json:"commmentCount"          csv:"commmentCount"`
	CompanyFollowerCount int    `json:"companyFollowerCount"   csv:"companyFollowerCount"`
	Freshness            string `json:"freshness"              csv:"freshness"`
	LikesCount           int    `json:"likesCount"             csv:"likesCount"`
	PostLink             string `json:"postLink"               csv:"postLink"`
	PublishDate          string `json:"publishDate"            csv:"publishDate"`
	ShareURN             string `json:"shareURN"               csv:"shareURN"`
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
	stopIteration := false
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		if stopIteration {
			return
		}
		header := s.Find("div[data-test-id=main-feed-activity-card__entity-lockup]")
		footer := s.Find(".main-feed-activity-card__social-actions")

		activityURN := strings.TrimSpace(s.AttrOr("data-activity-urn", ""))
		author := strings.TrimSpace(header.Find(".leading-open").Text())
		authorLinkedInUrl := cleanURL(strings.TrimSpace(header.Find(".leading-open").AttrOr("href", "")))
		commmentCount, _ := extractPostComments(strings.TrimSpace(footer.Find("span[data-test-id=social-actions__comments]").Text()))
		freshness := strings.Split(strings.TrimSpace(header.Find("div span time").Text()), "\n")[0]
		likesCount, _ := extractPostLikes(strings.TrimSpace(footer.Find("span[data-test-id=social-actions__reaction-count]").Text()))
		postLink := cleanURL(doc.Find("head link").AttrOr("href", ""))
		shareURN := strings.TrimSpace(s.AttrOr("data-attributed-urn", ""))

		var companyFollowerCount int
		var authorTitle string
		if strings.Contains(authorLinkedInUrl, "linkedin.com/company") {
			authorTitle = ""
			companyFollowerCount, _ = extractCompanyFollowers(strings.TrimSpace(header.Find("div p").Text()))
		} else {
			authorTitle = strings.TrimSpace(header.Find("div p").Text())
			companyFollowerCount = 0
		}

		post = Post{
			ActivityURN:          activityURN,
			Author:               author,
			AuthorLinkedInUrl:    authorLinkedInUrl,
			AuthorTitle:          authorTitle,
			CommentCount:         commmentCount,
			CompanyFollowerCount: companyFollowerCount,
			Freshness:            freshness,
			LikesCount:           likesCount,
			PostLink:             postLink,
			ShareURN:             shareURN,
		}

		// Print the post for testing
		if debug {
			log.Printf("Post: %+v", post)
		}

		stopIteration = true
	})
	return &post, nil
}

func extractPostLikes(s string) (int, error) {
	// Remove commas and convert to integer
	numStr := strings.ReplaceAll(s, ",", "")
	return strconv.Atoi(numStr)
}

func extractPostComments(s string) (int, error) {
	// Define a regex pattern to match numbers
	re := regexp.MustCompile(`(\d+,?\d*)\s*Comments`)
	match := re.FindStringSubmatch(s)

	if len(match) < 2 {
		return 0, fmt.Errorf("no match found")
	}

	// Remove commas and convert to integer
	numStr := strings.ReplaceAll(match[1], ",", "")
	return strconv.Atoi(numStr)
}
