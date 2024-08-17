package services

import (
	"errors"
    "sync"
    "github.com/NiksonGo/REST-API/models"
    
)

var(
todos        = []models.Todo{}
nextID       = 1
availableIDs = []int{}
mu           sync.Mutex
)

func GetAllTodos() []models.Todo {
    mu.Lock()
    defer mu.Unlock()
    return todos
}

func GetTodoByID(id int) (*models.Todo, error) {
    mu.Lock()
    defer mu.Unlock()

    for _, todo := range todos {
        if todo.ID == id {
            return &todo, nil
        }
    }
    return nil, errors.New("Задача не найдена")
}

func CreateTodo(title string, done bool) models.Todo {
    mu.Lock()
    defer mu.Unlock()

    var id int
    if len(availableIDs) > 0 {
        id = availableIDs[0]
        availableIDs = availableIDs[1:]
    } else {
        id = nextID
        nextID++
    }

    todo := models.Todo{ID: id, Title: title, Done: done}
    todos = append(todos, todo)
    return todo
}

func UpdateTodoByID(id int, title string, done bool) (*models.Todo, error) {
    mu.Lock()
    defer mu.Unlock()

    for i, todo := range todos {
        if todo.ID == id {
            todos[i].Title = title
            todos[i].Done = done
            return &todos[i], nil
        }
    }
    return nil, errors.New("Задача не найдена")
}

func DeleteTodoByID(id int) error {
    mu.Lock()
    defer mu.Unlock()

    for i, todo := range todos {
        if todo.ID == id {
            todos = append(todos[:i], todos[i+1:]...)
            availableIDs = append(availableIDs, id)
            return nil
        }
    }
    return errors.New("Задача не найдена")
}



