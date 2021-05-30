package storage

import (
	"database/sql"
	"log"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(provider *sql.DB) *Storage {
	storage := Storage{db: provider}
	storage.createSchema()
	return &storage
}

func (s *Storage) createSchema() error {
	stmt, err := s.db.Prepare(`
		CREATE TABLE IF NOT EXISTS userinfo (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id VARCHAR(255) NOT NULL,
			access_token VARCHAR(255) NOT NULL,
			username VARCHAR(255) NOT NULL,
			email VARCHAR (255) NOT NULL
		)`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	log.Println("::: Database started...")
	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) Save(user User) (User, error) {
	stmt, err := s.db.Prepare("INSERT INTO userinfo(session_id, access_token, username, email) values(?,?,?,?)")

	if err != nil {
		return User{}, err
	}

	res, err := stmt.Exec(user.SessionId, user.AccessToken, user.Username, user.Email)
	id, err := res.LastInsertId()

	if err != nil {
		return User{}, err
	}

	user.Id = id
	return user, nil
}

func (s *Storage) Find(sessionId string) (User, error) {
	var user User
	row := s.db.QueryRow("SELECT id, session_id, access_token, username, email FROM userinfo WHERE session_id = ?", sessionId)

	if err := row.Scan(&user.Id, &user.SessionId, &user.AccessToken, &user.Username, &user.Email); err != nil {
		return User{}, err
	}

	return user, nil
}
