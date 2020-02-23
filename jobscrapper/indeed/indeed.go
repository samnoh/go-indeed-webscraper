package indeed

import (
	"encoding/csv"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/samnoh/go-indeed-webscraper/helpers"
)

type extractedJob struct {
	id          string
	title       string
	location    string
	salary      string
	description string
}

var baseURL string = "https://nz.indeed.com/jobs?limit=50&q="

func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("data-jk")
	title := helpers.CleanString(card.Find(".title>a").Text())
	location := helpers.CleanString(card.Find(".sjcl").Text())
	salary := helpers.CleanString(card.Find(".salaryText").Text())
	description := helpers.CleanString(card.Find(".summary").Text())

	return extractedJob{
		id:          id,
		title:       title,
		location:    location,
		salary:      salary,
		description: description,
	}
}

func getNumberOfPages() int {
	res, err := http.Get(baseURL)

	helpers.CheckErr(err)
	helpers.CheckStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	helpers.CheckErr(err)

	nPages := 0
	doc.Find(".pagination").Each(func(i int, page *goquery.Selection) {
		nPages = page.Find("a").Length()
	})

	return nPages
}

func getPage(page int) []extractedJob {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	res, err := http.Get(pageURL)

	helpers.CheckErr(err)
	helpers.CheckStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	helpers.CheckErr(err)

	var jobs []extractedJob
	doc.Find(".jobsearch-SerpJobCard").Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})

	return jobs
}

func writeJobs(jobs []extractedJob) {
	// create a file
	file, err := os.Create("jobs.csv")
	helpers.CheckErr(err)

	// csv
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Description"}

	wErr := w.Write(headers)
	helpers.CheckErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://nz.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.description}
		jwErr := w.Write(jobSlice)
		helpers.CheckErr(jwErr)
	}
}

func Scrap(keyword string) {
	baseURL = baseURL + keyword
	var jobs []extractedJob

	totalNumberOfPages := getNumberOfPages()
	for i := 0; i < totalNumberOfPages; i++ {
		jobs = append(jobs, getPage(i)...)
	}

	writeJobs(jobs)
}
