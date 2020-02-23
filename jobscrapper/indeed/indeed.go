package indeed

import (
	"encoding/csv"
	"fmt"
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

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := helpers.CleanString(card.Find(".title>a").Text())
	location := helpers.CleanString(card.Find(".sjcl").Text())
	salary := helpers.CleanString(card.Find(".salaryText").Text())
	description := helpers.CleanString(card.Find(".summary").Text())

	c <- extractedJob{
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

func getPage(page int, scrapC chan<- []extractedJob) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	res, err := http.Get(pageURL)

	helpers.CheckErr(err)
	helpers.CheckStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	helpers.CheckErr(err)

	var jobs []extractedJob
	c := make(chan extractedJob)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	scrapC <- jobs
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
		wErr := w.Write(jobSlice)
		helpers.CheckErr(wErr)
	}
}

func Scrap(keyword string) {
	baseURL = baseURL + keyword
	var jobs []extractedJob

	totalNumberOfPages := getNumberOfPages()
	c := make(chan []extractedJob)
	for i := 0; i < totalNumberOfPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalNumberOfPages; i++ {
		jobsPerPage := <-c
		jobs = append(jobs, jobsPerPage...)
	}

	fmt.Println(jobs)
	writeJobs(jobs)
}
