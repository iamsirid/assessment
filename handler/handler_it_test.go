//go:build integration

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

	bodyExpense := database.Expense{
		Title:  "test",
		Amount: 100.0,
		Note:   "test",
		Tags:   []string{"test1", "test2"},
	}
	bodyBytes, err := json.Marshal(bodyExpense)

	if err != nil {
		t.Errorf("Error in encoding request body: %v", err)
	}

	body := bytes.NewBuffer(bodyBytes)

	gotExpense := database.Expense{}

	res := request(http.MethodPost, uri("expenses"), body)
	err = res.Decode(&gotExpense)

	if err != nil {
		t.Errorf("Error in decoding response: %v", err)
	}

	assert.Equal(t, res.StatusCode, http.StatusCreated)
	assert.Equal(t, gotExpense.Title, bodyExpense.Title)
	assert.Equal(t, gotExpense.Amount, bodyExpense.Amount)
	assert.Equal(t, gotExpense.Note, bodyExpense.Note)
	assert.Equal(t, len(gotExpense.Tags), len(bodyExpense.Tags))

	for i := range gotExpense.Tags {
		assert.Equal(t, gotExpense.Tags[i], bodyExpense.Tags[i])
	}

}

func TestGetExpenseHandler(t *testing.T) {
	toSeedExpense := database.Expense{
		Title:  "test",
		Amount: 100.0,
		Note:   "test",
		Tags:   []string{"test1", "test2"},
	}

	toSeedBytes, err := json.Marshal(toSeedExpense)

	if err != nil {
		t.Errorf("Error in encoding request body: %v", err)
	}

	origExpense := seedExpense(t, string(toSeedBytes))

	gotExpense := database.Expense{}

	res := request(http.MethodGet, uri("expenses/"+strconv.Itoa(origExpense.Id)), nil)
	err = res.Decode(&gotExpense)

	if err != nil {
		t.Errorf("Error in decoding response: %v", err)
	}

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, gotExpense.Title, toSeedExpense.Title)
	assert.Equal(t, gotExpense.Amount, toSeedExpense.Amount)
	assert.Equal(t, gotExpense.Note, toSeedExpense.Note)
	assert.Equal(t, len(gotExpense.Tags), len(toSeedExpense.Tags))

	for i := range gotExpense.Tags {
		assert.Equal(t, gotExpense.Tags[i], toSeedExpense.Tags[i])
	}

}

func TestUpdateExpenseById(t *testing.T) {

	toSeedExpense := database.Expense{
		Title:  "test",
		Amount: 100.0,
		Note:   "test",
		Tags:   []string{"test1", "test2"},
	}

	toSeedBytes, err := json.Marshal(toSeedExpense)

	if err != nil {
		t.Errorf("Error in encoding request body: %v", err)
	}

	origExpense := seedExpense(t, string(toSeedBytes))

	bodyExpense := database.Expense{
		Title:  "test-edited",
		Amount: 120.0,
		Note:   "test-edited",
		Tags:   []string{"test1", "test2", "test3"},
	}
	bodyBytes, err := json.Marshal(bodyExpense)

	if err != nil {
		t.Errorf("Error in encoding request body: %v", err)
	}

	body := bytes.NewBuffer(bodyBytes)

	gotExpense := database.Expense{}

	res := request(http.MethodPut, uri("expenses/"+strconv.Itoa(origExpense.Id)), body)
	err = res.Decode(&gotExpense)

	if err != nil {
		t.Errorf("Error in decoding response: %v", err)
	}

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, gotExpense.Title, bodyExpense.Title)
	assert.Equal(t, gotExpense.Amount, bodyExpense.Amount)
	assert.Equal(t, gotExpense.Note, bodyExpense.Note)
	assert.Equal(t, len(gotExpense.Tags), len(bodyExpense.Tags))

	for i := range gotExpense.Tags {
		assert.Equal(t, gotExpense.Tags[i], bodyExpense.Tags[i])
	}

}

func TestGetAllExpenses(t *testing.T) {

	toSeedExpenses := []database.Expense{{
		Title:  "test",
		Amount: 100.0,
		Note:   "test",
		Tags:   []string{"test1", "test2"},
	}, {
		Title:  "test2",
		Amount: 120.0,
		Note:   "test2",
		Tags:   []string{"test1", "test2", "test3"},
	}}

	for _, toSeedExpense := range toSeedExpenses {
		toSeedBytes, err := json.Marshal(toSeedExpense)

		if err != nil {
			t.Errorf("Error in encoding request body: %v", err)
		}

		seedExpense(t, string(toSeedBytes))
	}

	gotExpenses := []database.Expense{}

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&gotExpenses)

	expenseLen := len(gotExpenses)

	lastTwoIndexOfGotExpense := []int{expenseLen - 2, expenseLen - 1}

	if err != nil {
		t.Errorf("Error in decoding response: %v", err)
	}

	assert.Equal(t, res.StatusCode, http.StatusOK)

	for i, indexOfGotExpense := range lastTwoIndexOfGotExpense {
		assert.Equal(t, gotExpenses[indexOfGotExpense].Title, toSeedExpenses[i].Title)
		assert.Equal(t, gotExpenses[indexOfGotExpense].Amount, toSeedExpenses[i].Amount)
		assert.Equal(t, gotExpenses[indexOfGotExpense].Note, toSeedExpenses[i].Note)
		assert.Equal(t, len(gotExpenses[indexOfGotExpense].Tags), len(toSeedExpenses[i].Tags))

		for j := range gotExpenses[indexOfGotExpense].Tags {
			assert.Equal(t, gotExpenses[indexOfGotExpense].Tags[j], toSeedExpenses[i].Tags[j])
		}

	}

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
