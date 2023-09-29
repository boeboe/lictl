#!/usr/bin/env python3

import argparse
import datetime
import json
import sys

from locations import *

from bs4 import BeautifulSoup
from pathlib import Path
from urllib.parse import quote_plus, urljoin, urlparse
from torpy.http.requests import TorRequests

JOBSBASEURL='https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search'
AUTHWALLURL='https://www.linkedin.com/authwall'


def save_html(*, html, file):
  with open(file, "w") as f:
      f.write(html)


def save_parsed_results(*, parsed_results, outputfolder):
  json_object = json.dumps(parsed_results, sort_keys=True, indent = 2)
  with open(Path(outputfolder, "results.json"), "w") as jsonfile:
    jsonfile.write(json_object)

  with open(Path(outputfolder, "results.csv"), "w") as csvfile:
    csvfile.write("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s\n" %("company_name","company_linkedin_url","job_date",
      "job_location","job_title","job_url","job_urn","search_keyword","search_location","search_region"))
    for json_job in parsed_results:
      csvfile.write("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s\n" %(json_job['company_name'],json_job['company_linkedin_url'],
        json_job['job_date'],json_job['job_location'],json_job['job_title'],json_job['job_url'],
        json_job['job_urn'],json_job['search_keyword'],json_job['search_location'],json_job['search_region']))


def tor_request(*, tor_requests, url):
  while True:
    try:
      with tor_requests.get_session() as sess:
        return sess.get(url).text
    except:
      print("Tor requests failed... trying again")
      continue


def search(*, tor_requests, location, keyword, outputfolder):
  offset = 0
  while True:
    keyword_safe = quote_plus(keyword)
    search_url = f"{JOBSBASEURL}?location={location}&keywords={keyword_safe}&start={offset}"
    print(f"Going to fetch {search_url}")
    html_text = tor_request(tor_requests=tor_requests, url=search_url)
    soup = BeautifulSoup(html_text, 'html.parser')
    num_jobs = len(soup.find_all('li'))
    print(f"Result: {num_jobs}")

    if num_jobs != 25: 
      if AUTHWALLURL in html_text:
        continue
      elif num_jobs > 0:
        save_html(html=soup.prettify(), file=Path(outputfolder, f"results-{offset}.html"))
        break
      else:
        break
    else:
      save_html(html=soup.prettify(), file=Path(outputfolder, f"results-{offset}.html"))
      offset += 25
      continue


def parse(*, location, region, keyword, outputfolder):
  json_results = []
  job_urn_results = []
  
  for html_file in Path(outputfolder).rglob(f'results-*.html'):

    with open(html_file) as fp:
      soup = BeautifulSoup(fp, "lxml")

    for job in soup.find_all("li"):
      try:
        company_name = job.find("h4").get_text().strip().replace("|", ",")
        company_linkedin_url = urljoin(job.div.h4.a["href"], urlparse(job.div.h4.a["href"]).path)
        job_date = job.find("time", class_="job-search-card__listdate").get("datetime").strip()
        job_location = job.find("span", class_="job-search-card__location").get_text().strip()
        job_title = job.find("h3", class_="base-search-card__title").get_text().strip().replace("|", ",")
        job_url = urljoin(job.div.a["href"], urlparse(job.div.a["href"]).path)
        job_urn = job.find("div")["data-entity-urn"].split(":")[-1]
      except:
        try:
          company_name = job.find("h4").get_text().strip().replace("|", ",")
          company_linkedin_url = ""
          job_date = job.find("time", class_="job-search-card__listdate").get("datetime").strip()
          job_location = job.find("span", class_="job-search-card__location").get_text().strip()
          job_title = job.find("h3", class_="base-search-card__title").get_text().strip().replace("|", ",")
          job_url = urljoin(job.a["href"], urlparse(job.a["href"]).path)
          job_urn = job.find("a")["data-entity-urn"].split(":")[-1]
        except:
          continue

      if (job_urn in job_urn_results):
        continue
      else:
        job_urn_results.append(job_urn)
        json_results.append({ "company_name":company_name,
                              "company_linkedin_url":company_linkedin_url,
                              "job_date":job_date,
                              "job_location":job_location,
                              "job_title":job_title,
                              "job_url":job_url,
                              "job_urn":job_urn,
                              "search_keyword":keyword,
                              "search_location":location,
                              "search_region":region } )
  return json_results


def main(argv):
  parser = argparse.ArgumentParser()

  parser.add_argument('--regions', type=str, required=True)
  parser.add_argument('--keywords', type=str, required=True)

  args = parser.parse_args()
  regs = args.regions
  kws = args.keywords
  outf = Path().absolute().joinpath(f"output-{datetime.datetime.now().strftime('%Y-%m-%d_%H-%M-%S')}")

  print('Going to seach for jobs with the following parameters')
  print(f"  regions: \"{regs}\"")
  print(f"  keywords: \"{kws}\"")
  print(f"  outputfolder: \"{outf}\"")

  Path(outf).mkdir(parents=True, exist_ok=True)
  parsed_results_total = []
  parsed_results_total_unique_jobs_urn = []
  parsed_results_total_unique_jobs_count = 0

  for reg in regs.split(","):
    if not has_region(reg):
      print(f"Skipping unknown region \"{reg}\", try one of {get_regions()}")
      continue
    for loc in get_locations_for_region(reg):
      for kw in kws.split(","):
        resultfolder = Path(outf).joinpath(reg).joinpath(loc).joinpath(kw)
        resultfolder.mkdir(parents=True, exist_ok=True)

        with TorRequests() as tor_reqs:
          search(tor_requests=tor_reqs, location=loc, keyword=kw, outputfolder=resultfolder)

        parsed_results = parse(location=loc, region=reg, keyword=kw, outputfolder=resultfolder)

        for job in parsed_results:
          if job["job_urn"] not in parsed_results_total_unique_jobs_urn:
            parsed_results_total_unique_jobs_count += 1
            parsed_results_total.append(job)
            parsed_results_total_unique_jobs_urn.append(job["job_urn"])

        if len(parsed_results) > 0:
          save_parsed_results(parsed_results=parsed_results, outputfolder=resultfolder)
          print(f"Successfully parsed {len(parsed_results)} jobs in region \"{reg}\", location \"{loc}\" with keyword \"{kw}\"")
        else:
          print(f"Nothing to parse for jobs in region \"{reg}\", location \"{loc}\" with keyword \"{kw}\"")

        # Save all intermediate resuls to not loose agregated results in case of intermittent failure.
        save_parsed_results(parsed_results=parsed_results_total, outputfolder=outf)
        print(f"Successfully parsed {len(parsed_results_total)} jobs in region \"{reg}\"")

  print(f"Successfully finished scraping and found {parsed_results_total_unique_jobs_count} unique jobs")
  print(f"Results are available at results.json/results.csv in the following folder: {outf}")

if __name__ == "__main__":
   main(sys.argv[1:])