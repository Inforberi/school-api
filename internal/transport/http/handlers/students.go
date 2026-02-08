package handlers

import (
	"fmt"
	"net/http"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, students Path!"))
	fmt.Println("Hello, students Path!")
}
