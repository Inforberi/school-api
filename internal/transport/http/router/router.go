package router

import (
	"net/http"
	"restapi/internal/transport/http/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/teachers/", handlers.TeachersHandler)
	mux.HandleFunc("/students", handlers.StudentsHandler)
	mux.HandleFunc("/execs", handlers.ExecsHandler)

	return mux
}
