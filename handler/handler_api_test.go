package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/iamsirid/assessment/database"
)

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

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expect status code to be %v, got %v", http.StatusCreated, res.StatusCode)
	}

	if expense.Title != "test" {
		t.Errorf("Expect title to be test, got %v", expense.Title)
	}

	if expense.Amount != 100.0 {
		t.Errorf("Expect amount to be 100.0, got %v", expense.Amount)
	}

	if expense.Note != "test" {
		t.Errorf("Expect note to be test, got %v", expense.Note)
	}

	if len(expense.Tags) != 2 {
		t.Errorf("Expect tags length to be 2, got %v", len(expense.Tags))
	}

	if expense.Tags[0] != "test1" {
		t.Errorf("Expect tag to be test1, got %v", expense.Tags[0])
	}

	if expense.Tags[1] != "test2" {
		t.Errorf("Expect tag to be test2, got %v", expense.Tags[1])
	}

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

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expect status code to be %v, got %v", http.StatusOK, res.StatusCode)
	}

	if expense.Title != "test" {
		t.Errorf("Expect title to be test, got %v", expense.Title)
	}

	if expense.Amount != 100.0 {
		t.Errorf("Expect amount to be 100.0, got %v", expense.Amount)
	}

	if expense.Note != "test" {
		t.Errorf("Expect note to be test, got %v", expense.Note)
	}

	if len(expense.Tags) != 2 {
		t.Errorf("Expect tags length to be 2, got %v", len(expense.Tags))
	}

	if expense.Tags[0] != "test1" {
		t.Errorf("Expect tag to be test1, got %v", expense.Tags[0])
	}

	if expense.Tags[1] != "test2" {
		t.Errorf("Expect tag to be test2, got %v", expense.Tags[1])
	}
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

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expect status code to be %v, got %v", http.StatusCreated, res.StatusCode)
	}

	if gotExpense.Title != "test-edited" {
		t.Errorf("Expect title to be test-edited, got %v", gotExpense.Title)
	}

	if gotExpense.Amount != 120.0 {
		t.Errorf("Expect amount to be 120.0, got %v", gotExpense.Amount)
	}

	if gotExpense.Note != "test-edited" {
		t.Errorf("Expect note to be test-edited, got %v", gotExpense.Note)
	}

	if len(gotExpense.Tags) != 3 {
		t.Errorf("Expect tags length to be 3, got %v", len(gotExpense.Tags))
	}

	if gotExpense.Tags[0] != "test1" {
		t.Errorf("Expect tag to be test1, got %v", gotExpense.Tags[0])
	}

	if gotExpense.Tags[1] != "test2" {
		t.Errorf("Expect tag to be test2, got %v", gotExpense.Tags[1])
	}

	if gotExpense.Tags[2] != "test3" {
		t.Errorf("Expect tag to be test3, got %v", gotExpense.Tags[1])
	}
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
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expect status code to be %v, got %v", http.StatusCreated, res.StatusCode)
	}

	if expenses[expenseLen-2].Title != "test" {
		t.Errorf("Expect title to be test, got %v", expenses[expenseLen-2].Title)
	}

	if expenses[expenseLen-2].Amount != 100.0 {
		t.Errorf("Expect amount to be 100.0, got %v", expenses[expenseLen-2].Amount)
	}

	if expenses[expenseLen-2].Note != "test" {
		t.Errorf("Expect note to be test, got %v", expenses[expenseLen-2].Note)
	}

	if len(expenses[expenseLen-2].Tags) != 2 {
		t.Errorf("Expect tags length to be 2, got %v", len(expenses[expenseLen-2].Tags))
	}

	if expenses[expenseLen-2].Tags[0] != "test1" {
		t.Errorf("Expect tag to be test1, got %v", expenses[expenseLen-2].Tags[0])
	}

	if expenses[expenseLen-2].Tags[1] != "test2" {
		t.Errorf("Expect tag to be test2, got %v", expenses[expenseLen-2].Tags[1])
	}

	if expenses[expenseLen-1].Title != "test2" {
		t.Errorf("Expect title to be test2, got %v", expenses[expenseLen-1].Title)
	}

	if expenses[expenseLen-1].Amount != 120.0 {
		t.Errorf("Expect amount to be 120.0, got %v", expenses[expenseLen-1].Amount)
	}

	if expenses[expenseLen-1].Note != "test2" {
		t.Errorf("Expect note to be test2, got %v", expenses[expenseLen-1].Note)
	}

	if len(expenses[expenseLen-1].Tags) != 3 {
		t.Errorf("Expect tags length to be 3, got %v", len(expenses[expenseLen-1].Tags))
	}

	if expenses[expenseLen-1].Tags[0] != "test1" {
		t.Errorf("Expect tag to be test1, got %v", expenses[expenseLen-1].Tags[0])
	}

	if expenses[expenseLen-1].Tags[1] != "test2" {
		t.Errorf("Expect tag to be test2, got %v", expenses[expenseLen-1].Tags[1])
	}

	if expenses[expenseLen-1].Tags[2] != "test3" {
		t.Errorf("Expect tag to be test3, got %v", expenses[expenseLen-1].Tags[1])
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
