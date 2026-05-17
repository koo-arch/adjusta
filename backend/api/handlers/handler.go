package handlers

import (
	"github.com/koo-arch/adjusta-backend/api"
)

type Handler struct {
	Server *api.Server
}

func NewHandler(server *api.Server) *Handler {
	return &Handler{
		Server: server,
	}
}