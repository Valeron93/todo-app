package model

import (
	"context"
)

type User struct {
	Id             int64
	Username       string
	HashedPassword string
}

type Session struct {
	Token string
	User  User
}

type UserRepo interface {
	RegisterUser(username string, password string) (User, error)
	DeleteUser(username string) error
	GetByUsername(username string) (User, error)
	Login(username string, password string) (User, error)
}

type SessionManager interface {
	CreateSession(userId int64) (string, error)
	GetSession(token string) (Session, error)
	RevokeSession(token string) error
}

// A type to safely store Session in a context.Context
type sessionCtx struct{}

// Extract Session value from context.Context
func SessionFromCtx(ctx context.Context) (Session, bool) {
	s, ok := ctx.Value(sessionCtx{}).(Session)
	return s, ok
}

// Same as SessionFromCtx, but
// panics if session is not found
func SessionFromCtxMust(ctx context.Context) Session {
	s, ok := SessionFromCtx(ctx)
	if !ok {
		panic("failed to get session from context")
	}

	return s
}

// Create new context.Context with Session as value
func CtxWithSession(ctx context.Context, s Session) context.Context {
	return context.WithValue(ctx, sessionCtx{}, s)
}
