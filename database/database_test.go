//go:build !integration

package database

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/lib/pq"

	"github.com/google/go-cmp/cmp"
)

type FakeDatabaseHelper struct{}

func (h *FakeDatabaseHelper) ConnectToDatabase(databaseUrl string) (*sql.DB, error) {
	db, _, err := sqlmock.New()

	return db, err
}

func (h *FakeDatabaseHelper) ReadSqlFile(filename string) (string, error) {
	return "", nil
}

func (h *FakeDatabaseHelper) CreateTable(db *sql.DB) error {
	return nil
}

func TestInitDatabase(t *testing.T) {

	dummyUrl := "postgres://localhost:5432/assessment"

	db, err := InitDatabase(dummyUrl, &FakeDatabaseHelper{})

	if err != nil {
		t.Errorf("Error in InitDatabase: %v", err)
	}

	if db == nil {
		t.Errorf("Database is nil")
	}

}

func TestInsertData(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("Error in creating mock database: %v", err)
	}

	insertExpense := Expense{Title: "test", Amount: 100.0, Note: "test", Tags: []string{"test1", "test2"}}

	wantId := 1

	mockReturnRow := sqlmock.NewRows([]string{"id"}).AddRow(wantId)

	mock.ExpectQuery("INSERT INTO expenses").WithArgs(insertExpense.Title, insertExpense.Amount, insertExpense.Note, pq.Array(insertExpense.Tags)).WillReturnRows(mockReturnRow)

	gotId, err := InsertData(db, Expense{Title: "test", Amount: 100.0, Note: "test", Tags: []string{"test1", "test2"}})

	if err != nil {
		t.Errorf("Error in InsertData: %v", err)
	}

	if wantId != gotId {
		t.Errorf("Expect ID to be %v, got %v", wantId, gotId)
	}
}

func TestGetData(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("Error in creating mock database: %v", err)
	}

	wantExpense := Expense{Id: 1, Title: "test", Amount: 100.0, Note: "test", Tags: []string{"test1", "test2"}}

	mockRow := mock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(wantExpense.Id, wantExpense.Title, wantExpense.Amount, wantExpense.Note, pq.Array(wantExpense.Tags))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM expenses")).WithArgs(1).WillReturnRows(mockRow)

	gotExpense, err := GetData(db, 1)
	if err != nil {
		t.Errorf("Error in GetData: %v", err)
	}

	if !cmp.Equal(wantExpense, gotExpense) {
		t.Errorf("Expect expense to be %v, got %v", wantExpense, gotExpense)
	}
}

func TestUpdateData(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("Error in creating mock database: %v", err)
	}
	originalExpense := Expense{Id: 1, Title: "test", Amount: 100.0, Note: "test", Tags: []string{"test1", "test2"}}

	mock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(originalExpense.Id, originalExpense.Title, originalExpense.Amount, originalExpense.Note, pq.Array(originalExpense.Tags))

	updateExpense := Expense{Id: 1, Title: "test-updated", Amount: 200.0, Note: "test-updated", Tags: []string{"test1", "test2", "testt3"}}

	mockUpdateRow := mock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(updateExpense.Id, updateExpense.Title, updateExpense.Amount, updateExpense.Note, pq.Array(updateExpense.Tags))

	mock.ExpectQuery(regexp.QuoteMeta("UPDATE expenses")).WithArgs(updateExpense.Title, updateExpense.Amount, updateExpense.Note, pq.Array(updateExpense.Tags), updateExpense.Id).WillReturnRows(mockUpdateRow)

	gotExpense, err := UpdateData(db, updateExpense.Id, updateExpense)

	if err != nil {
		t.Errorf("Error in UpdateData: %v", err)
	}

	if !cmp.Equal(updateExpense, gotExpense) {
		t.Errorf("Expect return expense to be %v, got %v", updateExpense, gotExpense)
	}
}

func TestGetAllData(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("Error in creating mock database: %v", err)
	}

	wantExpenses := []Expense{
		{Id: 1, Title: "test", Amount: 100.0, Note: "test", Tags: []string{"test1", "test2"}},
		{Id: 2, Title: "test2", Amount: 120.0, Note: "test2", Tags: []string{"test1", "test2", "test3"}}}

	mockRows := mock.NewRows([]string{"id", "title", "amount", "note", "tags"})

	for _, wantExpense := range wantExpenses {
		mockRows = mockRows.AddRow(wantExpense.Id, wantExpense.Title, wantExpense.Amount, wantExpense.Note, pq.Array(wantExpense.Tags))
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM expenses")).WillReturnRows(mockRows)

	gotExpenses, err := GetAllData(db)
	if err != nil {
		t.Errorf("Error in GetAllData: %v", err)
	}

	for i := 0; i < len(gotExpenses); i++ {
		if !cmp.Equal(wantExpenses[i], gotExpenses[i]) {
			t.Errorf("Expect expense to be %v, got %v", wantExpenses[i], gotExpenses[i])
		}
	}
}
