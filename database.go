package main

import (
	"database/sql"
	"os"
)

type Model interface {
	OpenDB() error
	CloseDB()
	New(string, string)
	FetchDetail(string) (string, string, string)
	FetchPassword(string) string
	Update(string, string, string)
	Delete(string)
}

type UserModel struct {
	db *sql.DB
}

func (user *UserModel) OpenDB() error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	user.db = db
	return err
}

func (user *UserModel) CloseDB() {
	_ = user.db.Close()
}

func (user *UserModel) New(userId, password string) {
	_, _ = user.db.Exec("INSERT INTO users (user_id, password) VALUES ($1, $2)", userId, password)
}

func (user *UserModel) FetchDetail(id string) (userId, nickname, comment string) {
	user.db.QueryRow("SELECT user_id, nickname, comment FROM users WHERE user_id = $1", id).Scan(&userId, &nickname, &comment)
	return
}

func (user *UserModel) FetchPassword(id string) (password string) {
	user.db.QueryRow("SELECT password FROM users WHERE user_id = $1", id).Scan(&password)
	return
}

func (user *UserModel) Update(userId, nickname, comment string) {
	_, _ = user.db.Exec("UPDATE users SET nickname = $1, comment = $2 WHERE user_id = $3", nickname, comment, userId)
}

func (user *UserModel) Delete(id string) {
	_, _ = user.db.Exec("DELETE FROM users WHERE user_id = $1", id)
}
