package main

import (
    "log"
    "net/http"
    "github.com/NiksonGo/REST-API/router"
)

func main() {
    mux := router.SetupRouter()
    log.Println("Сервер запущен на порту: 8080")
    err := http.ListenAndServe(":8080", mux)
    log.Fatal(err)
}
