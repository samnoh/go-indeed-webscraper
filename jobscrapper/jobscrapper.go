package jobscrapper

import (
	"log"

	"github.com/samnoh/go-indeed-webscraper/jobscrapper/indeed"
)

var jobWebsites = map[string]int{
	"indeed": 1,
}

func Start(websiteName string, keyword string) {
	switch jobWebsites[websiteName] {
	case 1:
		indeed.Scrap(keyword)
	default:
		log.Fatalln("No such option")
	}
}
