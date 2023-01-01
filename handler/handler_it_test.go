package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/iamsirid/assessment/database"

	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
)

func setupServer() {

	db, err := database.InitDatabase(os.Getenv("DATABASE_URL"), &database.DatabaseHelper{})

	if err != nil {
		panic(err)
	}

	eh := echo.New()

	go func(e *echo.Echo) {

		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set("db", db)
				return next(c)
			}
		})

		e.POST("/expenses", CreateExpenseHandler)

		e.GET("/expenses/:id", GetExpenseByIdHandler)

		e.PUT("/expenses/:id", UpdateExpenseByIdHandler)

		e.GET("/expenses", GetAllExpensesHandler)

		e.Start(os.Getenv("PORT"))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost%v", os.Getenv("PORT")), 3*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
}

func init() {
	setupServer()
}

func TestCreateExpenseHandler(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "test",
		"amount": 100.0,
		"note": "test",
		"tags": ["test1", "test2"]
	}`)
	expense := database.Expense{}

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&expense)

	if err != nil {
		t.Errorf("Error in decoding response: %v", err)
	}

	assert.Equal(t, res.StatusCode, http.StatusCreated)
	assert.Equal(t, expense.Title, "test")
	assert.Equal(t, expense.Amount, 100.0)
	assert.Equal(t, expense.Note, "test")
	assert.Equal(t, len(expense.Tags), 2)
	assert.Equal(t, expense.Tags[0], "test1")
	assert.Equal(t, expense.Tags[1], "test2")

}

func TestGetExpenseHandler(t *testing.T) {
	origExpense := seedExpense(t, `{
		"title": "test",
		"amount": 100.0,
		"note": "test",
		"tags": ["test1", "test2"]
	}`)

	expense := database.Expense{}

	res := request(http.MethodGet, uri("expenses/"+strconv.Itoa(origExpense.Id)), nil)
	err := res.Decode(&expense)

	if err != nil {
		t.Errorf("Error in decoding response: %v", err)
	}

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, expense.Title, "test")
	assert.Equal(t, expense.Amount, 100.0)
	assert.Equal(t, expense.Note, "test")
	assert.Equal(t, len(expense.Tags), 2)
	assert.Equal(t, expense.Tags[0], "test1")
	assert.Equal(t, expense.Tags[1], "test2")

}

func TestUpdateExpenseById(t *testing.T) {

	origExpense := seedExpense(t, `{
		"title": "test",
		"amount": 100.0,
		"note": "test",
		"tags": ["test1", "test2"]
	}`)

	body := bytes.NewBufferString(`{
		"title": "test-edited",
		"amount": 120.0,
		"note": "test-edited",
		"tags": ["test1", "test2","test3"]
	}`)

	gotExpense := database.Expense{}

	res := request(http.MethodPut, uri("expenses/"+strconv.Itoa(origExpense.Id)), body)
	err := res.Decode(&gotExpense)

	if err != nil {
		t.Errorf("Error in decoding response: %v", err)
	}

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, gotExpense.Title, "test-edited")
	assert.Equal(t, gotExpense.Amount, 120.0)
	assert.Equal(t, gotExpense.Note, "test-edited")
	assert.Equal(t, len(gotExpense.Tags), 3)
	assert.Equal(t, gotExpense.Tags[0], "test1")
	assert.Equal(t, gotExpense.Tags[1], "test2")
	assert.Equal(t, gotExpense.Tags[2], "test3")

}

func TestGetAllExpenses(t *testing.T) {

	seedExpense(t, `{
		"title": "test",
		"amount": 100.0,
		"note": "test",
		"tags": ["test1", "test2"]
	}`)
	seedExpense(t, `{
		"title": "test2",
		"amount": 120.0,
		"note": "test2",
		"tags": ["test1", "test2" ,"test3"]
	}`)

	expenses := []database.Expense{}

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&expenses)

	expenseLen := len(expenses)

	if err != nil {
		t.Errorf("Error in decoding response: %v", err)
	}

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, expenses[expenseLen-2].Title, "test")
	assert.Equal(t, expenses[expenseLen-2].Amount, 100.0)
	assert.Equal(t, expenses[expenseLen-2].Note, "test")
	assert.Equal(t, len(expenses[expenseLen-2].Tags), 2)
	assert.Equal(t, expenses[expenseLen-2].Tags[0], "test1")
	assert.Equal(t, expenses[expenseLen-2].Tags[1], "test2")

	assert.Equal(t, expenses[expenseLen-1].Title, "test2")
	assert.Equal(t, expenses[expenseLen-1].Amount, 120.0)
	assert.Equal(t, expenses[expenseLen-1].Note, "test2")
	assert.Equal(t, len(expenses[expenseLen-1].Tags), 3)
	assert.Equal(t, expenses[expenseLen-1].Tags[0], "test1")
	assert.Equal(t, expenses[expenseLen-1].Tags[1], "test2")
	assert.Equal(t, expenses[expenseLen-1].Tags[2], "test3")

}

func seedExpense(t *testing.T, bodyString string) database.Expense {
	var expense database.Expense

	body := bytes.NewBufferString(bodyString)

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&expense)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return expense
}

func uri(paths ...string) string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":2565"
	}
	host := fmt.Sprintf("http://localhost%s", port)
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
