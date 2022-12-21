package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/iamsirid/assessment/database"

	"github.com/labstack/echo/v4"
)

type Err struct {
	Message string `json:"message"`
}

func main() {

	db, err := database.InitDatabase(os.Getenv("DATABASE_URL"), &database.DatabaseHelper{})

	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.POST("/expenses", func(c echo.Context) error {
		expense := database.Expense{}
		err := c.Bind(&expense)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
		}

		id, err := database.InsertData(db, expense)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}

		expense.Id = id

		return c.JSON(http.StatusOK, expense)

	})

	e.GET("/expenses/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
		expense, err := database.GetData(db, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, expense)
	})

	e.PUT("/expenses/:id", func(c echo.Context) error {
		payloadExpense := database.Expense{}
		err := c.Bind(&payloadExpense)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
		updatedExpense, err := database.UpdateData(db, id, payloadExpense)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, updatedExpense)

	})

	e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
