package handler

import (
	"net/http"
)

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// handlers passed through middleware
	logoutHandler := http.HandlerFunc(h.AuthHandler.LogOut)
	postCreateHandler := http.HandlerFunc(h.PostHandler.PostCreate)
	postUpdateHandler := http.HandlerFunc(h.PostHandler.PostUpdate)
	errorsHandler := http.HandlerFunc(errorHandler)
	mux.HandleFunc("/register", h.AuthHandler.RegisterUser)
	mux.HandleFunc("/signin", h.AuthHandler.SignIn)
	postGetHandler := http.HandlerFunc(h.PostHandler.PostGet)
	postGetAllHandler := http.HandlerFunc(h.PostHandler.PostGetAll)
	postGetMypostsHandler := http.HandlerFunc(h.PostHandler.GetMyPosts)

	mux.HandleFunc("/register", h.AuthHandler.RegisterUser) // +
	mux.HandleFunc("/signin", h.AuthHandler.SignIn)         // +
	mux.Handle("/logout", h.SessionMiddleware(logoutHandler))

	mux.Handle("/post/create", h.SessionMiddleware(postCreateHandler)) // add PostHandler
	mux.Handle("/post/update", h.SessionMiddleware(postUpdateHandler)) //
	mux.Handle("/", h.ErrorMiddleware(errorsHandler))

	_ = postGetAllHandler
	//mux.Handle("/", h.SessionMiddleware(postGetAllHandler))

	mux.Handle("/post/create", h.SessionMiddleware(postCreateHandler))
	mux.Handle("/post/update", h.SessionMiddleware(postUpdateHandler))
	mux.Handle("/post/get", h.SessionMiddleware(postGetHandler))
	mux.Handle("/post/myposts", h.SessionMiddleware(postGetMypostsHandler))

	return mux
}

/*
+ /post/get/all (all posts, filters)
+ /post/update
+ /post/get
+ /post/myposts

/comment/create
/comment/delete
/like/create
/like/delete
/session/refresh

*/
