package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

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

		return c.JSON(http.StatusCreated, expense)

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

	e.GET("/expenses", func(c echo.Context) error {
		expenses, err := database.GetAllData(db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, expenses)
	})

	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
