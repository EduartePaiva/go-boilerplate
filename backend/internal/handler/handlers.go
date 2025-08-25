package handler

import (
	"github.com/eduartepaiva/go-boilerplate/internal/server"
	"github.com/eduartepaiva/go-boilerplate/internal/service"
)

type Handlers struct{}

func NewHandlers(s *server.Server, services service.Services) *Handlers {
	return &Handlers{}
}
