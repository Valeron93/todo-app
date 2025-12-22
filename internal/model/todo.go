package model

import (
	"database/sql"
	"errors"
)

type Todo struct {
	Id     int64
	Action string
	UserId int64
}

type TodoRepo interface {
	GetAllForUser(userId int64) ([]Todo, error)
	CreateForUser(userId int64, action string) (Todo, error)
	Delete(id int64) error
}

type todoRepoSql struct {
	db *sql.DB
}

func NewTodoRepoSql(db *sql.DB) TodoRepo {
	return &todoRepoSql{
		db: db,
	}
}

// CreateForUser implements [TodoRepo].
func (t *todoRepoSql) CreateForUser(userId int64, action string) (Todo, error) {
	var newTodo Todo

	err := t.db.QueryRow(
		`INSERT INTO todos (action, user_id) 
		 VALUES (?, ?) 
		 RETURNING id, action, user_id`,
		action, userId,
	).Scan(&newTodo.Id, &newTodo.Action, &newTodo.UserId)
	if err != nil {
		return Todo{}, err
	}
	return newTodo, nil
}

func (t *todoRepoSql) GetAllForUser(userId int64) ([]Todo, error) {
	result := []Todo{}
	rows, err := t.db.Query(`
		SELECT id, action FROM todos WHERE user_id = ?
	`, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.Action); err != nil {
			return nil, err
		}
		result = append(result, todo)
	}
	return result, nil
}

func (t *todoRepoSql) Delete(id int64) error {

	result, err := t.db.Exec(`DELETE FROM todos WHERE id = ?`, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
