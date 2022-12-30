package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/iamsirid/assessment/database"
	"github.com/iamsirid/assessment/handler"

	"github.com/labstack/echo/v4"
)

func main() {

	db, err := database.InitDatabase(os.Getenv("DATABASE_URL"), &database.DatabaseHelper{})

	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authKey := c.Request().Header.Get("Authorization")
			if authKey != "November 10, 2009" {
				return c.JSON(http.StatusUnauthorized, "Unauthorized")
			}
			return next(c)
		}
	})

	e.POST("/expenses", handler.CreateExpenseHandler)

	e.GET("/expenses/:id", handler.GetExpenseByIdHandler)

	e.PUT("/expenses/:id", handler.UpdateExpenseByIdHandler)

	e.GET("/expenses", handler.GetAllExpensesHandler)

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
