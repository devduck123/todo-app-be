package usecases

import "github.com/devduck123/todo-app-be/entities"

func GetTodos(repo TodosRepository) ([]entities.Todo, error) {
	todos, err := repo.GetAllTodos()
	if err != nil {
		return nil, ErrInternal
	}

	return todos, nil
}
