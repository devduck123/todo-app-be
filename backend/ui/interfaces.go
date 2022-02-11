package ui

import "github.com/devduck123/todo-app-be/entities"

type Service interface {
	GetAllTodos() ([]entities.Todo, error)
}
