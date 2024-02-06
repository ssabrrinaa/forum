package handler

import (
	"net/http"
)

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	logoutHandler := http.HandlerFunc(h.AuthHandler.LogOut)
	postCreateHandler := http.HandlerFunc(h.PostHandler.PostCreate)
	postUpdateHandler := http.HandlerFunc(h.PostHandler.PostUpdate)
	errorsHandler := http.HandlerFunc(errorHandler)
	mux.HandleFunc("/register", h.AuthHandler.RegisterUser)
	mux.HandleFunc("/signin", h.AuthHandler.SignIn)
	mux.Handle("/logout", h.SessionMiddleware(logoutHandler))

	mux.Handle("/post/create", h.SessionMiddleware(postCreateHandler)) // add PostHandler
	mux.Handle("/post/update", h.SessionMiddleware(postUpdateHandler)) //
	mux.Handle("/", h.ErrorMiddleware(errorsHandler))
	return mux
}

/*
/ (all posts, filters)
/post/update
/post/get
/comment/create
/comment/delete
/like/create
/like/delete
/session/refresh
*/
