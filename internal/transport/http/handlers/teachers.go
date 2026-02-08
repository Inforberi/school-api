package handlers

import (
	"fmt"
	"net/http"
)

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, teachers Path!"))
	fmt.Println("Hello, teachers Path!")

	switch r.Method {
	case http.MethodGet:

	case http.MethodPost:
		w.Write([]byte("Hello, teachers Path!"))
		fmt.Println("Hello, teachers Path!")
	}
}
