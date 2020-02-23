package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/samnoh/go-indeed-webscraper/controllers"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Static("/static", "public/assets")

	// Route
	controllers.Controllers(e)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
