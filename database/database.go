package database

import (
	"database/sql"

	"github.com/lib/pq"
)

type Expense struct {
	Id     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type IDatabaseHelper interface {
	ConnectToDatabase(string) (*sql.DB, error)
	CreateTable(*sql.DB) error
}

type DatabaseHelper struct {
}

func (h *DatabaseHelper) ConnectToDatabase(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseUrl)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (h *DatabaseHelper) CreateTable(db *sql.DB) error {

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`

	_, err := db.Exec(createTableQuery)

	if err != nil {
		return err
	}

	return nil
}

func InitDatabase(databaseUrl string, dbh IDatabaseHelper) (*sql.DB, error) {

	db, err := dbh.ConnectToDatabase(databaseUrl)

	if err != nil {
		return nil, err
	}

	err = dbh.CreateTable(db)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertData(db *sql.DB, expense Expense) (int, error) {
	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id",
		expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags))
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil

}

func GetData(db *sql.DB, id int) (Expense, error) {
	row := db.QueryRow("SELECT * FROM expenses WHERE id = $1", id)
	expense := Expense{}
	err := row.Scan(&expense.Id, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
	if err != nil {
		return Expense{}, err
	}
	return expense, nil
}

func UpdateData(db *sql.DB, id int, payloadExpense Expense) (Expense, error) {
	row := db.QueryRow("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING *",
		payloadExpense.Title, payloadExpense.Amount, payloadExpense.Note, pq.Array(payloadExpense.Tags), id)
	updatedExpense := Expense{}
	err := row.Scan(&updatedExpense.Id, &updatedExpense.Title, &updatedExpense.Amount, &updatedExpense.Note, pq.Array(&updatedExpense.Tags))
	if err != nil {
		return Expense{}, err
	}
	return updatedExpense, nil
}

func GetAllData(db *sql.DB) ([]Expense, error) {
	rows, err := db.Query("SELECT * FROM expenses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	expenses := []Expense{}
	for rows.Next() {
		expense := Expense{}
		err := rows.Scan(&expense.Id, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}
