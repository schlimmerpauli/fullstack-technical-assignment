package products

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerRejectsUnsupportedMethod(t *testing.T) {
	handler := newTestHandler()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/products", nil)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, recorder.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	if response.Error != "method not allowed" {
		t.Fatalf("expected method not allowed error, got %q", response.Error)
	}
}

func TestHandlerRejectsInvalidQueryParameters(t *testing.T) {
	handler := newTestHandler()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/products?bestseller=maybe", nil)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	if response.Error != "invalid bestseller: must be true or false" {
		t.Fatalf("expected bestseller validation error, got %q", response.Error)
	}
}

func TestHandlerRejectsInvalidSortQueryParameter(t *testing.T) {
	handler := newTestHandler()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/products?sort=price", nil)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	if response.Error != "invalid sort: must be popularity" {
		t.Fatalf("expected sort validation error, got %q", response.Error)
	}
}

func TestHandlerReturnsProductsResponseShape(t *testing.T) {
	handler := newTestHandler()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/products?page=1&pageSize=2", nil)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response payload: %v", err)
	}

	products, ok := payload["products"].([]any)
	if !ok {
		t.Fatal("expected products array in response")
	}
	if len(products) != 2 {
		t.Fatalf("expected 2 products on first page, got %d", len(products))
	}

	firstProduct, ok := products[0].(map[string]any)
	if !ok {
		t.Fatal("expected first product to be an object")
	}
	for _, field := range []string{"id", "name", "category", "brand", "condition", "price", "discount_percent", "bestseller", "colors", "image_url", "stock"} {
		if _, exists := firstProduct[field]; !exists {
			t.Fatalf("expected product field %q in response", field)
		}
	}
	for _, field := range []string{"product_id"} {
		if _, exists := firstProduct[field]; exists {
			t.Fatalf("did not expect internal product field %q in response", field)
		}
	}

	pagination, ok := payload["pagination"].(map[string]any)
	if !ok {
		t.Fatal("expected pagination object in response")
	}
	if pagination["page"] != float64(1) {
		t.Fatalf("expected page 1, got %v", pagination["page"])
	}
	if pagination["pageSize"] != float64(2) {
		t.Fatalf("expected pageSize 2, got %v", pagination["pageSize"])
	}
	if pagination["hasMore"] != true {
		t.Fatalf("expected hasMore true, got %v", pagination["hasMore"])
	}
}

func TestHandlerFiltersByRepeatedCategoryAndCondition(t *testing.T) {
	handler := newTestHandler()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/products?category=audio&category=desktop&condition=excellent&condition=good&page=1&pageSize=20", nil)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var payload ListResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response payload: %v", err)
	}

	if len(payload.Products) != 8 {
		t.Fatalf("expected 8 filtered products, got %d", len(payload.Products))
	}
	if payload.Pagination.Total != 8 {
		t.Fatalf("expected total 8 filtered products, got %d", payload.Pagination.Total)
	}
}

func newTestHandler() http.Handler {
	service := NewService(newTestRepository())
	handler := NewHandler(service)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return mux
}
