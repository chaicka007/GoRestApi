package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=iamroot dbname=taskdb sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Ошибка пинга базы:", err)
	}

	createTable()
}

func createTable() {
	query := `
CREATE TABLE IF NOT EXISTS tasks (
	id SERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT,
	status TEXT NOT NULL DEFAULT 'pending'
);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}
}
