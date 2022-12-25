package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/iamsirid/assessment/database"
	"github.com/labstack/echo/v4"
)

type Err struct {
	Message string `json:"message"`
}

func CreateExpenseHandler(c echo.Context) error {

	db := c.Get("db").(*sql.DB)

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
}

func GetExpenseByIdHandler(c echo.Context) error {
	db := c.Get("db").(*sql.DB)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	expense, err := database.GetData(db, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, expense)
}

func UpdateExpenseByIdHandler(c echo.Context) error {

	db := c.Get("db").(*sql.DB)

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
}

func GetAllExpensesHandler(c echo.Context) error {
	db := c.Get("db").(*sql.DB)

	expenses, err := database.GetAllData(db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, expenses)

}
