package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/eduartepaiva/go-boilerplate/internal/config"
	"github.com/eduartepaiva/go-boilerplate/internal/database"
	"github.com/eduartepaiva/go-boilerplate/internal/handler"
	"github.com/eduartepaiva/go-boilerplate/internal/logger"
	"github.com/eduartepaiva/go-boilerplate/internal/repository"
	"github.com/eduartepaiva/go-boilerplate/internal/router"
	"github.com/eduartepaiva/go-boilerplate/internal/server"
	"github.com/eduartepaiva/go-boilerplate/internal/service"
)

const DefaultContextTimeout = 30

func main() {
	cfg := config.LoadConfig()

	loggerService := logger.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()

	log := logger.NewLoggerWithService(cfg.Observability, loggerService)

	if cfg.Primary.Env != "local" {
		if err := database.Migrate(context.Background(), &log, cfg); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate database")
		}
	}

	srv, err := server.New(cfg, &log, loggerService)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
	}

	repos := repository.NewRepositories(srv)
	services, serviceErr := service.NewServices(srv, repos)
	if serviceErr != nil {
		log.Fatal().Err(err).Msg("could not create service")
	}
	handlers := handler.NewHandlers(srv, services)

	r := router.NewRouter(srv, handlers, services)

	srv.SetupHTTPServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		if err = srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}
	stop()
	cancel()

	log.Info().Msg("server exited properly")
}
