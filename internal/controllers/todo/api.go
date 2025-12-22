package todo

import (
	"github.com/Valeron93/todo-app/internal/model"
)

type TodoController struct {
	todoRepo model.TodoRepo
}

type IndexTemplateData struct {
	Todos []model.Todo
	User  model.User
}

func New(todoRepo model.TodoRepo) *TodoController {
	return &TodoController{
		todoRepo: todoRepo,
	}
}
