package main

import "github.com/samnoh/job_scrapper/jobscrapper"

var baseURL string = "https://nz.indeed.com/jobs?q=python&limit=50"

func main() {
	jobscrapper.Scrap(baseURL)
}
