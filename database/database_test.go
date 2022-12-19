package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/lib/pq"
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

	mock.ExpectQuery("INSERT INTO expenses").WithArgs("test", 100.0, "test", pq.Array([]string{"test1", "test2"})).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := InsertData(db, Expense{Title: "test", Amount: 100.0, Note: "test", Tags: []string{"test1", "test2"}})

	if err != nil {
		t.Errorf("Error in InsertData: %v", err)
	}

	if id != 1 {
		t.Errorf("Expect ID to be 1, got %v", id)
	}
}
