package main

import "github.com/samnoh/go-indeed-webscraper/jobscrapper"

func main() {
	jobscrapper.Scrap("https://nz.indeed.com/jobs?q=python&limit=50")
}
