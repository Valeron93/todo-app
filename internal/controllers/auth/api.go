package auth

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Valeron93/todo-app/internal/model"
	"github.com/Valeron93/todo-app/internal/view"
)

type AuthController struct {
	users    model.UserRepo
	sessions model.SessionManager
}

func New(users model.UserRepo, sessions model.SessionManager) *AuthController {
	return &AuthController{
		users:    users,
		sessions: sessions,
	}
}

func (c *AuthController) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user, err := c.sessions.GetUser(sessionCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(), "user", user),
		))
	})
}

func (c *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		if err := view.RenderPage(w, "register", nil); err != nil {
			log.Print(err)
		}
		return
	}

	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	// TODO: validate form
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "username or password is invalid", http.StatusBadRequest)
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

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		if err := view.RenderPage(w, "login", nil); err != nil {
			log.Print(err)
		}
		return
	}

	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)

}
