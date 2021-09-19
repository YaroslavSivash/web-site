package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/show", showNotes)
	mux.HandleFunc("/create", createNote)

	log.Println("Запуск сервера на http://127.0.0.1:8000")
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}