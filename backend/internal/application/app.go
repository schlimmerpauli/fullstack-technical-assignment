package application

import (
	"assignment-backend/internal/config"
	"assignment-backend/internal/server/httpserver"
	"context"
	"log"
)

type App struct {
	cfg        *config.Config
	httpServer *httpserver.Server
}

func New() *App {
	cfg := config.LoadConfig()

	return &App{
		cfg:        &cfg,
		httpServer: httpserver.New(&cfg),
	}
}

func (a *App) Start() error {
	log.Printf("product-service: starting http server on %s:%s", a.cfg.ServerHost, a.cfg.ServerPort)
	return a.httpServer.Start()
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Println("product-service: shutting down...")

	return a.httpServer.Shutdown(ctx)
}
