package usecases

import "github.com/devduck123/todo-app-be/entities"

type TodosRepository interface {
	GetAllTodos() ([]entities.Todo, error)
}
