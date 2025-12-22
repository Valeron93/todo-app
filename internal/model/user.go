package model

type User struct {
	Id             int64
	Username       string
	HashedPassword string
}

type UserRepo interface {
	RegisterUser(username string, password string) (User, error)
	DeleteUser(username string) error
	GetByUsername(username string) (User, error)
	Login(username string, password string) (User, error)
}

type SessionManager interface {
	CreateSession(userId int64) (string, error)
	GetUser(token string) (User, error)
	RevokeSession(token string) error
}
