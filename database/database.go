package database

import (
	"database/sql"
	"io"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func readSqlFile(filename string) (string, error) {
	file, err := os.Open("database/" + filename + ".sql")
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func createTable(db *sql.DB) error {
	createTableQuery, err := readSqlFile("create-table")
	if err != nil {
		return err
	}
	_, err = db.Exec(createTableQuery)

	if err != nil {
		return err
	}

	return nil

}

func InitDatabase() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Ping to database error", err)
	}

	log.Println("Connected to database")

	err = createTable(db)

	if err != nil {
		log.Fatal("Create table error", err)
	}

	log.Println("Table created")
}
