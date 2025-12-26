package model

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userRepoSql struct {
	db *sql.DB
}

type sessionManagerSql struct {
	db *sql.DB
}

func NewUserRepoSql(db *sql.DB) UserRepo {
	return &userRepoSql{
		db: db,
	}
}

func NewSessionManagerSql(db *sql.DB) SessionManager {
	return &sessionManagerSql{
		db: db,
	}
}

func (u *userRepoSql) DeleteUser(username string) error {
	panic("unimplemented")
}

func (u *userRepoSql) GetByUsername(username string) (User, error) {
	var user User
	err := u.db.QueryRow(
		`SELECT id, username, hashed_password
		 FROM users
		 WHERE username = ?`,
		username,
	).Scan(&user.Id, &user.Username, &user.HashedPassword)

	return user, err
}

// Login implements [UserRepo].
func (u *userRepoSql) Login(username string, password string) (User, error) {
	var user User
	err := u.db.QueryRow(
		`SELECT id, username, hashed_password
		 FROM users 
		 WHERE username = ?`,
		username,
	).Scan(&user.Id, &user.Username, &user.HashedPassword)

	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrInvalidCredentials
	}

	if err != nil {
		log.Print(err)
		return User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return User{}, ErrInvalidCredentials
	}

	return user, err

}

func (u *userRepoSql) RegisterUser(username string, password string) (User, error) {
	// TODO: validate username and password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	_, err = u.GetByUsername(username)

	if !errors.Is(err, sql.ErrNoRows) {
		return User{}, err
	}

	var newUser User
	err = u.db.QueryRow(
		`INSERT INTO users (username, hashed_password) 
		 VALUES (?, ?) 
		 RETURNING id, username, hashed_password`,
		username,
		hashedPassword,
	).Scan(&newUser.Id, &newUser.Username, &newUser.HashedPassword)

	return newUser, err
}

// CreateSession implements [SessionManager].
func (s *sessionManagerSql) CreateSession(userId int64) (string, error) {
	token := uuid.New().String()

	_, err := s.db.Exec(
		`INSERT INTO sessions (token, user_id)
		VALUES (?, ?)`,
		token, userId,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *sessionManagerSql) GetSession(token string) (Session, error) {

	var session Session
	err := s.db.QueryRow(
		`SELECT
			sessions.token,
    	 	users.id,
    		users.username
		 FROM sessions
		 JOIN users ON sessions.user_id = users.id
		 WHERE sessions.token = ?`,
		token,
	).Scan(&session.Token, &session.User.Id, &session.User.Username)
	return session, err
}

// GetUser implements [SessionManager].
func (s *sessionManagerSql) GetUser(token string) (User, error) {

	var user User
	err := s.db.QueryRow(
		`SELECT 
    	 	users.id,
    		users.username
		 FROM sessions
		 JOIN users ON sessions.user_id = users.id
		 WHERE sessions.token = ?`,
		token,
	).Scan(&user.Id, &user.Username)
	return user, err
}

// RevokeSession implements [SessionManager].
func (s *sessionManagerSql) RevokeSession(token string) error {

	result, err := s.db.Exec(`DELETE FROM sessions WHERE token = ?`, token)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRecord
	}

	return nil
}
