package jobscrapper

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/samnoh/job_scrapper/helpers"
)

type extractedJob struct {
	id          string
	title       string
	location    string
	salary      string
	description string
}

func Scrap(baseURL string) {
	var jobs []extractedJob
	totalPages := getPages(baseURL)

	for i := 0; i < totalPages; i++ {
		jobs = append(jobs, getPage(baseURL, i)...)
	}

	fmt.Println(jobs)
}

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

func getPage(baseURL string, page int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	res, err := http.Get(pageURL)
	helpers.CheckErr(err)
	helpers.CheckStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	helpers.CheckErr(err)

	doc.Find(".jobsearch-SerpJobCard").Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})

	return jobs
}

func getPages(baseURL string) int {
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
