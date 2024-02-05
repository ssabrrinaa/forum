package handler

import (
	"html/template"
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {
	errObj := r.Context().Value("error")
	if errObj == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, ok := errObj.(error)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	t, parseErr := template.ParseFiles("ui/templates/errors.html")
	if parseErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err := t.Execute(w, errObj)
	if err != nil {
		return
	}
}
