package auth

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Valeron93/todo-app/internal/model"
	"github.com/Valeron93/todo-app/internal/templates"
)

type AuthController struct {
	users    model.UserRepo
	sessions model.SessionManager
}

type SessionKey struct{}

func New(users model.UserRepo, sessions model.SessionManager) *AuthController {
	return &AuthController{
		users:    users,
		sessions: sessions,
	}
}

func (c *AuthController) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, ok := r.Context().Value(SessionKey{}).(model.Session)

		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (c *AuthController) AuthRedirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, ok := r.Context().Value(SessionKey{}).(model.Session)

		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (c *AuthController) InjectSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if sessionCookie, err := r.Cookie("session_token"); err == nil {
			if session, err := c.sessions.GetSession(sessionCookie.Value); err == nil {
				ctx = context.WithValue(r.Context(), SessionKey{}, session)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	// TODO: validate form
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if username == "" || password == "" || password != confirmPassword {
		err := templates.RegisterForm(templates.AuthFormData{
			Username:        username,
			Password:        password,
			ConfirmPassword: confirmPassword,
			Error:           "Username or password is invalid",
		}).Render(r.Context(), w)
		if err != nil {
			log.Print(err)
		}
		return
	}

	user, err := c.users.RegisterUser(username, password)

	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Print(err)
	}

	token, err := c.sessions.CreateSession(user.Id)

	if err != nil {
		http.Error(w, "internal server error: failed to create session", http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	// TODO: validate form
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := c.users.Login(username, password)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	token, err := c.sessions.CreateSession(user.Id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	w.Header().Add("HX-Redirect", "/")
}

func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {

	session := r.Context().Value(SessionKey{}).(model.Session)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})

	if err := c.sessions.RevokeSession(session.Token); err != nil {
		log.Printf("failed to revoke session: %v", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}
