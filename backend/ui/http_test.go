package ui_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devduck123/todo-app-be/entities"
	"github.com/devduck123/todo-app-be/ui"
	"github.com/gomagedon/expectate"
)

type MockService struct {
	err   error
	todos []entities.Todo
}

func (s MockService) GetAllTodos() ([]entities.Todo, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.todos, nil
}

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

type HTTPTest struct {
	name        string
	service     *MockService
	inputMethod string
	inputURL    string

	expectedStatus int
	expectedTodos  []entities.Todo
}

func TestHTTP(t *testing.T) {
	tests := getTests()
	tests = append(tests, getDisallowedMethodTests()...)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testHTTP(t, test)
		})
	}
}

// main test
func testHTTP(t *testing.T, test HTTPTest) {
	expect := expectate.Expect(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(test.inputMethod, test.inputURL, nil)

	server := ui.NewHTTP()
	server.UseService(test.service)

	server.ServeHTTP(w, r)

	var body []entities.Todo
	json.NewDecoder(w.Result().Body).Decode(&body)

	expect(w.Result().StatusCode).ToBe(test.expectedStatus)
	expect(body).ToEqual(test.expectedTodos)
}

func getTests() []HTTPTest {
	return []HTTPTest{
		{
			name:           "random error gives 500 status and no todos",
			service:        &MockService{err: fmt.Errorf("something bad happened")},
			inputMethod:    "GET",
			inputURL:       "http://mywebsite.com/todos/",
			expectedStatus: 500,
		},
		{
			name:           "random error gives 500 status and no todos",
			service:        &MockService{err: fmt.Errorf("something bad happened")},
			inputMethod:    "GET",
			inputURL:       "http://mywebsite.com/todos",
			expectedStatus: 500,
		},
		{
			name:           "wrong path gives 404 status and no todos",
			service:        &MockService{todos: dummyTodos},
			inputMethod:    "GET",
			inputURL:       "http://mywebsite.com/foo",
			expectedStatus: 404,
		},
		{
			name:           "wrong path gives 404 status and no todos",
			service:        &MockService{todos: dummyTodos},
			inputMethod:    "GET",
			inputURL:       "http://mywebsite.com/bar",
			expectedStatus: 404,
		},
		{
			name:           "returns todos from service if no error",
			service:        &MockService{todos: dummyTodos},
			inputMethod:    "GET",
			inputURL:       "http://mywebsite.com/todos/",
			expectedStatus: 200,
			expectedTodos:  dummyTodos,
		},
	}
}

func getDisallowedMethodTests() []HTTPTest {
	tests := []HTTPTest{}

	disallowedMethods := []string{
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
	}

	for _, method := range disallowedMethods {
		tests = append(tests, HTTPTest{
			name:           fmt.Sprintf("Method %s gives 405 status and no todos", method),
			service:        &MockService{todos: dummyTodos},
			inputURL:       "http://mywebsite.com/todos/",
			inputMethod:    method,
			expectedStatus: http.StatusMethodNotAllowed,
		})
	}

	return tests
}
