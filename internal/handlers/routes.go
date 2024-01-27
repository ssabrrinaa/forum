package handler

import (
	"net/http"
)

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	// handlers passed through middleware
	logoutHandler := http.HandlerFunc(h.AuthHandler.LogOut)
	postCreateHandler := http.HandlerFunc(h.PostHandler.PostCreate)
	postUpdateHandler := http.HandlerFunc(h.PostHandler.PostUpdate)

	mux.HandleFunc("/register", h.AuthHandler.RegisterUser)
	mux.HandleFunc("/signin", h.AuthHandler.SignIn)
	mux.Handle("/logout", h.SessionMiddleware(logoutHandler))

	mux.Handle("/post/create", h.SessionMiddleware(postCreateHandler)) // add PostHandler
	mux.Handle("/post/update", h.SessionMiddleware(postUpdateHandler)) //
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
