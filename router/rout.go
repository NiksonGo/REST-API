package router

import (
    "net/http"
    "github.com/NiksonGo/REST-API/controllers"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			controllers.GetTodos(w, r)
		case "POST":
			controllers.CreateTodo(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
           controllers.GetTodoByID(w, r)
        case "PUT":
           controllers.UpdateTodoByID(w, r)
        case "DELETE":
            controllers.DeleteTodoByID(w, r)
        default:
            http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        }
    })

    return mux
}
