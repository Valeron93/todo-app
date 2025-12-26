package controller

import (
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

func (c *AuthController) HandleRegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := view.Register().Render(r.Context(), w); err != nil {
		log.Println(err)
	}

}
func (c *AuthController) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := view.Login().Render(r.Context(), w); err != nil {
		log.Println(err)
	}
}

func NewAuth(users model.UserRepo, sessions model.SessionManager) *AuthController {
	return &AuthController{
		users:    users,
		sessions: sessions,
	}
}

func (c *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	// TODO: move this validation into model package
	if username == "" || password == "" || password != confirmPassword {
		err := view.RegisterForm(view.AuthFormData{
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

	session := model.SessionFromCtxMust(r.Context())

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
