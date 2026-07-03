package httpserver

import (
	"assignment-backend/internal/config"
	"assignment-backend/internal/products"
	"context"
	"fmt"
	"net/http"
)

// this could easily be its own package
type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config) *Server {
	productRepository := products.NewRepository()
	productService := products.NewService(productRepository)
	productHandler := products.NewHandler(productService)
	router := NewRouter(productHandler)
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}
func (s *Server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
