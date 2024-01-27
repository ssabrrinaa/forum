package handler

import (
	"forum/internal/handlers/authhandler"
	"forum/internal/handlers/posthandler"
	"forum/internal/services"
)

type Handler struct {
	AuthHandler *authhandler.AuthHandler
	PostHandler *posthandler.PostHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		AuthHandler: authhandler.NewAuthHandler(service.AuthService),
		PostHandler: posthandler.NewPostHandler(service.PostService),
	}
}
