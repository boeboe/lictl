package linkedin

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search?"

// Job represents the structure of a LinkedIn job.
type Job struct {
	CompanyLinkedInURL string `json:"companyLinkedInURL" csv:"companyLinkedInURL"`
	CompanyName        string `json:"companyName"        csv:"companyName"`
	DatePosted         string `json:"datePosted"         csv:"datePosted"`
	JobLink            string `json:"jobLink"            csv:"jobLink"`
	JobTitle           string `json:"jobTitle"           csv:"jobTitle"`
	JobURN             string `json:"jobURN"             csv:"jobURN"`
	Location           string `json:"location"           csv:"location"`
}

func (j *Job) CsvContent() string {
	if j == nil {
		return ""
	}
	return CsvContent(j)
}

func (j *Job) CsvHeader() string {
	if j == nil {
		return ""
	}
	return CsvHeader(j)
}

func (j *Job) Json() string {
	if j == nil {
		return ""
	}
	return Json(j)
}

type Jobs []*Job

func (js Jobs) Len() int {
	return len(js)
}

func (js Jobs) Get(i int) Serializable {
	return Serializable(js[i])
}

func SearchJobsOnline(regions []string, keywords []string, interval time.Duration, debug bool) (Jobs, error) {
	var allJobs []*Job

	for offset := 0; offset <= 975; offset += 25 {
		params := url.Values{}
		params.Add("location", strings.Join(regions, ","))
		params.Add("keywords", strings.Join(keywords, ","))
		params.Add("start", fmt.Sprintf("%d", offset))
		url := baseURL + params.Encode()
		if debug {
			fmt.Printf("going to fetch search url %v", url)
		}

		jobs, err := GetJobsFromSearchUrl(url, debug)
		if err != nil {
			if httpErr, ok := err.(*HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
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

func GetJobsFromSearchUrl(url string, debug bool) (Jobs, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LinkedIn jobs: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("received non-2xx response: %d %s", resp.StatusCode, resp.Status),
		}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var jobs []*Job
	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		var job Job

		companyName := strings.TrimSpace(s.Find("h4").Text())
		datePosted := strings.TrimSpace(s.Find(".job-search-card__listdate").AttrOr("datetime", ""))
		jobLink := cleanURL(s.Find(".base-card__full-link").AttrOr("href", ""))
		jobTitle := strings.TrimSpace(s.Find(".base-search-card__title").Text())
		jobURN := strings.Split(s.Find("div").AttrOr("data-entity-urn", ""), ":")[3]
		location := strings.TrimSpace(s.Find(".job-search-card__location").Text())

		var companyLinkedInURL string
		if href, exists := s.Find("h4 a").Attr("href"); exists {
			companyLinkedInURL = cleanURL(href)
		}

		job = Job{
			CompanyLinkedInURL: companyLinkedInURL,
			CompanyName:        companyName,
			DatePosted:         datePosted,
			JobLink:            jobLink,
			JobTitle:           jobTitle,
			JobURN:             jobURN,
			Location:           location,
		}

		if job.JobTitle != "" { // Only append if we found a job title
			jobs = append(jobs, &job)
		}
	})

	// Print the jobs for testing
	if debug {
		for _, job := range jobs {
			log.Printf("Job: %+v", job)
		}
	}

	return jobs, nil
}
