package handler

import (
	"forum/internal/handlers/authhandler"
	"forum/internal/handlers/posthandler"
	service "forum/internal/services"
)

type Handler struct {
	AuthHandler                *authhandler.AuthHandler
	PostHandler                *posthandler.PostHandler
	ExcludeSessionHandlersPath map[string]struct{}
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		AuthHandler: authhandler.NewAuthHandler(service.AuthService),
		PostHandler: posthandler.NewPostHandler(service.PostService),
		ExcludeSessionHandlersPath: map[string]struct{}{
			"/post/": {},
		},
	}
}
