package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
    
)

//структура Todo
type Todo struct {
    ID    int    `json: "id"`
    Title string `json: "title"`
    Done  bool   `json: "done"`
}

var(
    todos = []Todo{}
    nextID = 1
    mu     sync.Mutex
)
//маршрутизатор 
func main() {
    mux := http.NewServeMux()
//обработка маршрутов
    mux.HandleFunc("/todos", handleTodos)
    mux.HandleFunc("/todos/", handleTodo)
//запуск сервера
    log.Println("Сервер запущен в порту:8080")
    err := http.ListenAndServe(":8080", mux)
    log.Fatal(err) 
}
//функции для обработки HTTP-запросов и напрвления на соответствующие функ-обработчики в зависимости от метода(GET,POST и т.д)
func handleTodos(w http.ResponseWriter, r *http.Request) { //интерфейс и указатель на структуру запроса
    switch r.Method {
    case "GET":
        getTodos(w, r)
    case "POST":
        createTodo(w, r)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}
//функция обрабатывающая запросы,связанные с конкретной задачей(todo),идентифицируемой по ID.
func handleTodo(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        getTodoByID(w, r)
    case "PUT":
        updateTodoByID(w, r)
    case "DELETE":
        deleteTodoByID(w, r)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}
//реализация CRUD-операций(5 функций)
func getTodos(w http.ResponseWriter, r *http.Request) { //функция возвращает клиенту список задач в формате JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
    var todo Todo
    if err := json.NewDecoder(r.Body).Decode(&todo); err != nil { //декодер JSON 
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    } 
    mu.Lock() //блокирует мьютекс 
    defer mu.Unlock() //откладывает разблокировку мьютекса до завершения функции
    //установка уникального идентифик. для новой задачи
    todo.ID = nextID //даем новый идентифик.
    nextID++ //увелич.значение nextID для следующ.создаваемой задачи
    todos = append(todos, todo)//добавление новой задачи в список задач
    //отправка ответа клиенту
    w.Header().Set("Content-Type", "application/json")//установка заголовка HTTP-ответа Content-Type в значение aplication/json
    json.NewEncoder(w).Encode(todo)//кодирует новую задачу todo в JSON и записывает ее в тело HTTP-ответа,отправляя клиенту
}

func getTodoByID(w http.ResponseWriter, r *http.Request) {
    id, err := parseID(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    mu.Lock()
    defer mu.Unlock()
    //поиск задачи по идентификатору 
    for _, todo := range todos {
        if todo.ID == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(todo)
            return
        }
    }
    http.Error(w, "Задача не найдена", http.StatusNotFound)//обработка случая если задача не найдена
}
//функция обновляющая задачу по ее идентификатору
func updateTodoByID(w http.ResponseWriter, r *http.Request) {
    //парсинг идентификатора из запроса
    id, err := parseID(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    //декодирование обновлен.данных задачи из тела запроса
    var updatedTodo Todo
    if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    mu.Lock()
    defer mu.Unlock()
    for i, todo := range todos {
        if todo.ID == id {
            todos[i].Title = updatedTodo.Title
            todos[i].Done = updatedTodo.Done
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(todos[i])
            return
        }
    }
    http.Error(w, "Задача не найдена", http.StatusNotFound)
}
//функция для удаления задачи по идентификатору 
func deleteTodoByID(w http.ResponseWriter, r *http.Request) {
    id, err := parseID(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    mu.Lock()
    defer mu.Unlock()
    //поиск задачи по идентификатору и удаление
    for i, todo := range todos {
        if todo.ID == id {
            todos = append(todos[:i], todos[i + 1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    http.Error(w, "Задача не найдена", http.StatusNotFound)
}
//функция для извлечения идентификатора задачи из URL запроса.
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



