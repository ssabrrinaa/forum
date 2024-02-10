package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {
	errObj := r.Context().Value("error")
	fmt.Printf("Handler %v\n", errObj)
	if errObj == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, ok := errObj.(error)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	t, parseErr := template.ParseFiles("ui/templates/error.html")
	if parseErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err := t.Execute(w, errObj)
	if err != nil {
		return
	}
}
