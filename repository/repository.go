package repository

import (
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(provider *sql.DB) *UserRepository {
	repo := UserRepository{db: provider}
	return &repo
}

func (s *UserRepository) Close() {
	s.db.Close()
}

func (s *UserRepository) Save(user User) (User, error) {
	stmt, err := s.db.Prepare("INSERT INTO gh_user(session_id, access_token, username, email, repos) values($1,$2,$3,$4,$5)")

	if err != nil {
		return User{}, err
	}

	res, err := stmt.Exec(user.SessionId, user.AccessToken, user.Username, user.Email, user.Repos)
	id, err := res.LastInsertId()

	if err != nil {
		return User{}, err
	}

	user.Id = id
	return user, nil
}

func (s *UserRepository) Find(sessionId string) (User, error) {
	var user User
	row := s.db.QueryRow("SELECT id, session_id, access_token, username, email, repos FROM gh_user WHERE session_id = $1", sessionId)

	if err := row.Scan(&user.Id, &user.SessionId, &user.AccessToken,
		&user.Username, &user.Email, &user.Repos); err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *UserRepository) UpdateGitRepositories(sessionId string, reposUrl string) (User, error) {
	stmt, err := s.db.Prepare("UPDATE gh_user SET repos = $1 WHERE session_id = $2")

	if err != nil {
		return User{}, err
	}

	res, err := stmt.Exec(reposUrl, sessionId)
	rowsAffected, err := res.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return User{}, err
	}

	return s.Find(sessionId)
}
