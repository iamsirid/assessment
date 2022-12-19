package main

import (
	"os"

	"github.com/iamsirid/assessment/database"

	"github.com/labstack/echo/v4"
)

func main() {

	err := database.InitDatabase(os.Getenv("DATABASE_URL"), &database.DatabaseHelper{})

	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
