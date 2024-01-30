package handler

import (
	"net/http"
	"text/template"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {
	errObj := r.Context().Value("error")

	t, parseErr := template.ParseFiles("ui/templates/errors.html")
	if parseErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	t.Execute(w, errObj)
}
