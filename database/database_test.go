package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

	err := InitDatabase(dummyUrl, &FakeDatabaseHelper{})

	if err != nil {
		t.Errorf("Error in InitDatabase: %v", err)
	}

}
