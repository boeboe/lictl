package linkedin

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/boeboe/lictl/pkg/utils"
)

const baseURL = "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search?"

// Job represents the structure of a LinkedIn job.
type Job struct {
	JobTitle           string `json:"jobTitle"`
	CompanyName        string `json:"companyName"`
	CompanyLinkedInURL string `json:"companyLinkedInURL"`
	Location           string `json:"location"`
	DatePosted         string `json:"datePosted"`
	JobLink            string `json:"jobLink"`
	JobURN             string `json:"jobURN"`
}

func cleanURL(link string) string {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return link
	}
	parsedURL.RawQuery = ""
	return parsedURL.String()
}

func SearchJobsOnline(regions []string, keywords []string, interval time.Duration, debug bool) ([]Job, error) {
	var allJobs []Job

	for offset := 0; offset <= 975; offset += 25 {
		url := baseURL + "location=" + strings.Join(regions, ",") + "&keywords=" + strings.Join(keywords, ",") + fmt.Sprintf("&start=%d", offset)
		if debug {
			fmt.Printf("going to fetch search url %v", url)
		}

		jobs, err := SearchJobsPerPage(url, debug)
		if err != nil {
			if httpErr, ok := err.(*utils.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				return allJobs, err // Return the jobs fetched so far along with the error
			}
			return nil, err
		}
		if len(jobs) == 0 {
			break
		}
		allJobs = append(allJobs, jobs...)
		time.Sleep(interval)
	}
	return allJobs, nil
}

func SearchJobsPerPage(url string, debug bool) ([]Job, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LinkedIn jobs: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &utils.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("received non-2xx response: %d %s", resp.StatusCode, resp.Status),
		}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var jobs []Job
	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		var job Job

		jobTitle := strings.TrimSpace(s.Find(".base-search-card__title").Text())
		companyName := strings.TrimSpace(s.Find("h4").Text())
		location := strings.TrimSpace(s.Find(".job-search-card__location").Text())
		datePosted := strings.TrimSpace(s.Find(".job-search-card__listdate").AttrOr("datetime", ""))
		jobLink := cleanURL(s.Find(".base-card__full-link").AttrOr("href", ""))
		jobURN := strings.Split(s.Find("div").AttrOr("data-entity-urn", ""), ":")[3]

		var companyLinkedInURL string
		if href, exists := s.Find("h4 a").Attr("href"); exists {
			companyLinkedInURL = cleanURL(href)
		}

		job = Job{
			JobTitle:           jobTitle,
			CompanyName:        companyName,
			CompanyLinkedInURL: companyLinkedInURL,
			Location:           location,
			DatePosted:         datePosted,
			JobLink:            jobLink,
			JobURN:             jobURN,
		}

		if job.JobTitle != "" { // Only append if we found a job title
			jobs = append(jobs, job)
		}
	})

	// Print the jobs for testing
	if debug {
		for _, job := range jobs {
			log.Println(job)
		}
	}

	return jobs, nil
}
