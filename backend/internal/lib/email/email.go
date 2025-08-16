package email

import (
	"github.com/eduartepaiva/go-boilerplate/internal/config"
	"github.com/rs/zerolog"
)

type Client struct{}

func NewClient(config *config.Config, logger *zerolog.Logger) *Client {
	return &Client{}
}

func (e *Client) SendWelcomeEmail(to, firstName string) error {
	return nil
}
