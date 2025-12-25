package todo

import (
	"github.com/Valeron93/todo-app/internal/model"
)

type TodoController struct {
	todoRepo model.TodoRepo
}

func New(todoRepo model.TodoRepo) *TodoController {
	return &TodoController{
		todoRepo: todoRepo,
	}
}
