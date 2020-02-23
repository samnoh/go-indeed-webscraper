package main

import (
	"github.com/labstack/echo"
	"github.com/samnoh/go-indeed-webscraper/controllers"
)

func main() {
	e := echo.New()
	controllers.Controllers(e)
	e.Logger.Fatal(e.Start(":3000"))
}
