package todo

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Valeron93/todo-app/internal/model"
	"github.com/Valeron93/todo-app/internal/view"
)

func (a *TodoController) HandlePostTodo(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(model.User)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	action := strings.TrimSpace(r.FormValue("todo-action"))

	if action == "" {
		http.Error(w, "empty todo", http.StatusBadRequest)
		return
	}

	todo, err := a.todoRepo.CreateForUser(user.Id, action)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if err := view.RenderPartial(w, "todo-item", todo); err != nil {
		log.Print(err)
	}

}

func (a *TodoController) HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid path parameter", http.StatusBadRequest)
		return
	}

	if err := a.todoRepo.Delete(id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
