package handler

import (
	"gitlab.com/skeleton-golang/internal/infrastructure/container"
)

type Handler struct {
	userHandler *userHandler
}

func SetupHandler(container *container.Container) *Handler {
	return &Handler{
		userHandler: NewUserHandler(container.UserService),
	}
}

func (h *Handler) Validate() *Handler {
	if h.userHandler == nil {
		panic("userHandler is nil")
	}
	return h
}
