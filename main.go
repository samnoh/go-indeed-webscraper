package main

import "github.com/samnoh/job_scrapper/jobscrapper"

func main() {
	jobscrapper.Scrap("https://nz.indeed.com/jobs?q=python&limit=50")
}
