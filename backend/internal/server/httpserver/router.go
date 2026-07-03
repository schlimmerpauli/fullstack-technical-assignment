package httpserver

import (
	"assignment-backend/internal/products"
	"net/http"
)

func NewRouter(productHandler *products.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	productHandler.RegisterRoutes(mux)

	return mux
}
