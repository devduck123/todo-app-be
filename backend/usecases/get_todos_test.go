package usecases_test

import (
	"fmt"
	"testing"

	"github.com/devduck123/todo-app-be/entities"
	"github.com/devduck123/todo-app-be/usecases"
	"github.com/gomagedon/expectate"
)

var dummyTodos = []entities.Todo{
	{
		Title:       "todo 1",
		Description: "description of todo 1",
		IsCompleted: true,
	},
	{
		Title:       "todo 2",
		Description: "description of todo 2",
		IsCompleted: false,
	},
	{
		Title:       "todo 3",
		Description: "description of todo 3",
		IsCompleted: true,
	},
}

type BadTodosRepo struct{}

func (BadTodosRepo) GetAllTodos() ([]entities.Todo, error) {
	return nil, fmt.Errorf("error occurred")
}

type MockTodosRepo struct{}

func (MockTodosRepo) GetAllTodos() ([]entities.Todo, error) {
	return dummyTodos, nil
}

func TestGetTodos(t *testing.T) {
	// test
	t.Run("returns ErrInternal when TodosRepository returns error", func(t *testing.T) {
		expect := expectate.Expect(t)

		repo := new(BadTodosRepo)

		todos, err := usecases.GetTodos(repo)

		expect(err).ToBe(usecases.ErrInternal)
		if todos != nil {
			t.Fatalf("expected todos to be nil")
		}
	})

	// test
	t.Run("returns todos from TodosRepository", func(t *testing.T) {
		expect := expectate.Expect(t)

		repo := new(MockTodosRepo)

		todos, err := usecases.GetTodos(repo)

		expect(err).ToBe(nil)
		expect(todos).ToEqual(dummyTodos)
	})
}
