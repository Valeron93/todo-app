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

	originalAction := r.FormValue("todo-action")
	action := strings.TrimSpace(originalAction)

	// TODO: move this validation in model package
	if action == "" {
		view.TodoForm(originalAction, "Item cannot be empty").Render(r.Context(), w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, err := c.todoRepo.CreateForUser(session.User.Id, action)
	if err != nil {
		view.TodoForm(originalAction, "Internal server error").Render(r.Context(), w)
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}
	if err := view.TodoForm("", "").Render(r.Context(), w); err != nil {
		log.Print(err)
	}

	if err := view.TodoItemOutOfBand(todo).Render(r.Context(), w); err != nil {
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

	w.WriteHeader(http.StatusOK)
}
