package controllers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/NiksonGo/REST-API/services"
    "github.com/NiksonGo/REST-API/models"
)
func GetTodos(w http.ResponseWriter, r *http.Request) { //функция возвращает клиенту список задач в формате JSON
    todos := services.GetAllTodos()
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {
    id, err := parseID(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	todo, err := services.GetTodoByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
    var todo models.Todo
    if err := json.NewDecoder(r.Body).Decode(&todo); err != nil { //декодер JSON 
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    } 
	newTodo := services.CreateTodo(todo.Title, todo.Done)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(newTodo)
}

func UpdateTodoByID(w http.ResponseWriter, r *http.Request) {
    //парсинг идентификатора из запроса
    id, err := parseID(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	var updatedTodo models.Todo
    if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	todo, err := services.UpdateTodoByID(id, updatedTodo.Title, updatedTodo.Done)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}
func DeleteTodoByID(w http.ResponseWriter, r *http.Request) {
    id, err := parseID(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	if err := services.DeleteTodoByID(id); err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
	w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request) (int, error) {
    //объявление переменной для хранения извлечнного идентификатора
    var id int 
    //парсинг идентификатора из URL
    _, err := fmt.Sscanf(r.URL.Path, "/todos/%d", &id)
    if err != nil {
        return 0, fmt.Errorf("Неверный ID") 
    }
    //возвращаем идентификатор
    return id, nil

}
