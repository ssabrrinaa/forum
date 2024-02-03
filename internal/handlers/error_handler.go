package handler

import (
	"net/http"
	"text/template"
	"forum/internal/exceptions"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {
	errObj := r.Context().Value("error").(exceptions.BaseError)
	fmt.Println(errObj)
	t, parseErr := template.ParseFiles("ui/templates/errors.html")
	if parseErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	t.Execute(w, errObj)
}
