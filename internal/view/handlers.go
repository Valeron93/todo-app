package view

import (
	"log"
	"net/http"

	"github.com/Valeron93/todo-app/internal/model"
)

type ViewHandler struct {
	todoRepo model.TodoRepo
}

type TemplateData map[string]any

func NewViewHandler(todoRepo model.TodoRepo) *ViewHandler {
	return &ViewHandler{
		todoRepo: todoRepo,
	}
}

func (h *ViewHandler) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(model.User)
	todos, err := h.todoRepo.GetAllForUser(user.Id)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		"Todos": todos,
		"User":  user,
	}

	if err := renderPage(w, "index", data); err != nil {
		log.Print(err)
	}
}

func (h *ViewHandler) HandleRegisterPage(w http.ResponseWriter, r *http.Request) {
	h.renderOrRedirectIfLoggedIn("register", w, r)
}

func (h *ViewHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	h.renderOrRedirectIfLoggedIn("login", w, r)
}

func (h *ViewHandler) renderOrRedirectIfLoggedIn(page string, w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Context().Value("user").(model.User); ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := renderPage(w, page, nil); err != nil {
		log.Print(err)
	}
}
