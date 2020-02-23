package controllers

import (
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/samnoh/go-indeed-webscraper/helpers"
	"github.com/samnoh/go-indeed-webscraper/jobscraper"
)

const fileName string = "jobs.csv"

func handleHome(c echo.Context) error {
	return c.File("public/home.html")
}

func handleScrap(c echo.Context) error {
	defer os.Remove(fileName)

	keyword := strings.ToLower(helpers.CleanString(c.FormValue("keyword")))
	jobscraper.Start("indeed", keyword)
	return c.Attachment(fileName, "jobs-"+keyword+".csv")
}

func Controllers(e *echo.Echo) {
	e.GET("/", handleHome)

	e.POST("/scrap-job", handleScrap)
}
