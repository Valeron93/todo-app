package views

import (
	"log"
	"net/http"

	"github.com/Valeron93/todo-app/internal/controllers/auth"
	"github.com/Valeron93/todo-app/internal/model"
	"github.com/Valeron93/todo-app/internal/templates"
)

type ViewHandler struct {
	todoRepo model.TodoRepo
}

func NewViewHandler(todoRepo model.TodoRepo) *ViewHandler {
	return &ViewHandler{
		todoRepo: todoRepo,
	}
}

func (h *ViewHandler) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.UserKey{}).(model.User)
	todos, err := h.todoRepo.GetAllForUser(user.Id)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	if err := templates.Index(user, todos).Render(r.Context(), w); err != nil {
		log.Println(err)
	}
}

func (h *ViewHandler) HandleRegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.Register().Render(r.Context(), w); err != nil {
		log.Println(err)
	}

}
func (h *ViewHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.Login().Render(r.Context(), w); err != nil {
		log.Println(err)
	}
}
