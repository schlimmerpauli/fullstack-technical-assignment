package httpserver

import (
	"assignment-backend/internal/config"
	"assignment-backend/internal/products"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterAllowsConfiguredCORSOrigin(t *testing.T) {
	cfg := &config.Config{CORSAllowedOrigins: []string{"http://localhost:4173"}}
	handler := products.NewHandler(products.NewService(products.NewRepositoryWithPaths("..\\..\\..\\data\\metadata.json", "..\\..\\..\\data\\details.json")))
	router := NewRouter(cfg, handler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("Origin", "http://localhost:4173")

	router.ServeHTTP(recorder, request)

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:4173" {
		t.Fatalf("expected allowed origin header to be echoed, got %q", got)
	}
}

func TestRouterHandlesPreflight(t *testing.T) {
	cfg := &config.Config{CORSAllowedOrigins: []string{"*"}}
	handler := products.NewHandler(products.NewService(products.NewRepositoryWithPaths("..\\..\\..\\data\\metadata.json", "..\\..\\..\\data\\details.json")))
	router := NewRouter(cfg, handler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/products", nil)
	request.Header.Set("Origin", "http://localhost:4173")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected wildcard allow origin, got %q", got)
	}
}
