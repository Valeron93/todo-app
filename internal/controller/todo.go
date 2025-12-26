package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Valeron93/todo-app/internal/model"
	"github.com/Valeron93/todo-app/internal/view"
)

type TodoController struct {
	todoRepo model.TodoRepo
}

func NewTodo(todoRepo model.TodoRepo) *TodoController {
	return &TodoController{
		todoRepo: todoRepo,
	}
}

func (c *TodoController) HandleTodoListPage(w http.ResponseWriter, r *http.Request) {
	session := model.SessionFromCtxMust(r.Context())

	todos, err := c.todoRepo.GetAllForUser(session.User.Id)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	if err := view.Index(session.User, todos).Render(r.Context(), w); err != nil {
		log.Println(err)
	}
}

func (c *TodoController) HandlePostTodo(w http.ResponseWriter, r *http.Request) {

	session := model.SessionFromCtxMust(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	action := strings.TrimSpace(r.FormValue("todo-action"))

	// TODO: move this validation in model package
	if action == "" {
		http.Error(w, "empty todo", http.StatusBadRequest)
		return
	}

	todo, err := c.todoRepo.CreateForUser(session.User.Id, action)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if err := view.TodoItem(todo).Render(r.Context(), w); err != nil {
		log.Print(err)
	}
}

func (c *TodoController) HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid path parameter", http.StatusBadRequest)
		return
	}

	if err := c.todoRepo.Delete(id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
