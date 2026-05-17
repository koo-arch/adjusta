package middlewares

import (
	"github.com/koo-arch/adjusta-backend/api"
)

type Middleware struct {
	Server 	  *api.Server
}

func NewMiddleware(server *api.Server) *Middleware {
	return &Middleware{
		Server: server,
	}
}