package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteHandler struct {
	db *sql.DB
}

func (s *sqliteHandler) GetTodos() []*Todo {
	return nil
}

func (s *sqliteHandler) AddTodo(name string) *Todo {
	return nil
}

func (s *sqliteHandler) RemoveTodo(id int) bool {
	return false
}

func (s *sqliteHandler) CompleteTodo(id int, complete bool) bool {
	return false
}

func (s *sqliteHandler) Close() {
	s.db.Close()
}

func newSqliteHandler() DBHandler {
	database, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}
	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			completed BOOLEAN,
			createdAt DATETIME
		)`)
	statement.Exec()
	return &sqliteHandler{db: database}
}
