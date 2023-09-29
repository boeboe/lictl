package linkedin

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func SearchJobsOnline(regions []string, keywords []string) ([]Job, error) {
	var allJobs []Job

	for offset := 0; offset <= 975; offset += 25 {
		url := baseURL + "locations=" + strings.Join(regions, ",") + "&keywords=" + strings.Join(keywords, ",") + fmt.Sprintf("&start=%d", offset)
		jobs, err := SearchJobsPerPage(url)
		if err != nil {
			return nil, err
		}
		if len(jobs) == 0 {
			break
		}
		allJobs = append(allJobs, jobs...)
	}
	for _, job := range allJobs {
		log.Println(job)
	}
	return allJobs, nil
}

func SearchJobsPerPage(url string) ([]Job, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LinkedIn jobs: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var jobs []Job
	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		var job Job

		// First structure
		jobTitle := strings.TrimSpace(s.Find(".base-search-card__title").Text())
		if jobTitle != "" {
			job = Job{
				JobTitle:           jobTitle,
				CompanyName:        strings.TrimSpace(s.Find("h4").Text()),
				CompanyLinkedInURL: cleanURL(s.Find("h4 a").AttrOr("href", "")),
				Location:           strings.TrimSpace(s.Find(".job-search-card__location").Text()),
				DatePosted:         strings.TrimSpace(s.Find(".job-search-card__listdate").AttrOr("datetime", "")),
				JobLink:            cleanURL(s.Find(".base-card__full-link").AttrOr("href", "")),
				JobURN:             strings.Split(s.Find("div").AttrOr("data-entity-urn", ""), ":")[3],
			}
		} else {
			// Second structure (fallback)
			jobTitle = strings.TrimSpace(s.Find(".other-job-title-selector").Text()) // Replace with the correct selector
			if jobTitle != "" {
				job = Job{
					JobTitle:           jobTitle,
					CompanyName:        strings.TrimSpace(s.Find("h4").Text()),
					CompanyLinkedInURL: "",
					Location:           strings.TrimSpace(s.Find(".job-search-card__location").Text()),
					DatePosted:         strings.TrimSpace(s.Find(".job-search-card__listdate").AttrOr("datetime", "")),
					JobLink:            cleanURL(s.Find(".other-job-title-selector").AttrOr("href", "")), // Replace with the correct selector
					JobURN:             strings.Split(s.Find("a").AttrOr("data-entity-urn", ""), ":")[3],
				}
			}
		}

		if job.JobTitle != "" { // Only append if we found a job title
			jobs = append(jobs, job)
		}
	})

	// Print the jobs for testing
	for _, job := range jobs {
		log.Println(job)
	}

	return jobs, nil
}
