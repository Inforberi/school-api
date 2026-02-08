package handlers

import (
	"fmt"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, execs Path!"))
	fmt.Println("Hello, execs Path!")
}
