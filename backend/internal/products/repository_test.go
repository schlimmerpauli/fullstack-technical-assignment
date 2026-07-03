package products

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRepositoryListMergesDataFiles(t *testing.T) {
	repository := newTestRepository()

	result, err := repository.List(context.Background(), ListQuery{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if result.Pagination.Total != 16 {
		t.Fatalf("expected total 16 products, got %d", result.Pagination.Total)
	}
	if len(result.Products) != 16 {
		t.Fatalf("expected 16 products on first page, got %d", len(result.Products))
	}

	first := result.Products[0]
	if first.ID != "p1-excellent" {
		t.Fatalf("expected first product id p1-excellent, got %s", first.ID)
	}
	if first.ProductID != "p1" {
		t.Fatalf("expected first product base id p1, got %s", first.ProductID)
	}
	if first.Name != "iPhone 12" {
		t.Fatalf("expected first product name iPhone 12, got %s", first.Name)
	}
	if first.Brand != "Apple" {
		t.Fatalf("expected first product brand Apple, got %s", first.Brand)
	}
	if first.Category != "smartphone" {
		t.Fatalf("expected first product category smartphone, got %s", first.Category)
	}
	if first.Condition != "excellent" {
		t.Fatalf("expected first product condition excellent, got %s", first.Condition)
	}
	if math.Abs(first.Price-331.99) > 0.001 {
		t.Fatalf("expected first product price 331.99, got %.2f", first.Price)
	}
	if first.DiscountPercent != 20 {
		t.Fatalf("expected first product discount 20, got %d", first.DiscountPercent)
	}
	if !first.Bestseller {
		t.Fatal("expected first product to be bestseller")
	}
	if first.Stock != 18 {
		t.Fatalf("expected first product stock 18, got %d", first.Stock)
	}
	if len(first.Colors) != 3 {
		t.Fatalf("expected first product to have 3 colors, got %d", len(first.Colors))
	}
}

func TestRepositoryListFiltersBySearchAndPrice(t *testing.T) {
	repository := newTestRepository()
	maxPrice := 600.0

	result, err := repository.List(context.Background(), ListQuery{
		Search: "Mac Mini",
		PriceRange: &PriceRange{
			Max: &maxPrice,
		},
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if result.Pagination.Total != 1 {
		t.Fatalf("expected one filtered product, got %d", result.Pagination.Total)
	}
	if len(result.Products) != 1 {
		t.Fatalf("expected one product on page, got %d", len(result.Products))
	}
	if result.Products[0].ID != "p7-good" {
		t.Fatalf("expected Mac Mini good variant (p7-good), got %s", result.Products[0].ID)
	}
}

func TestRepositoryListFiltersByColorAndBestseller(t *testing.T) {
	repository := newTestRepository()
	bestseller := true

	result, err := repository.List(context.Background(), ListQuery{
		Color:      "white",
		Bestseller: &bestseller,
		Page:       1,
		PageSize:   10,
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if result.Pagination.Total != 2 {
		t.Fatalf("expected two filtered products, got %d", result.Pagination.Total)
	}
	if len(result.Products) != 2 {
		t.Fatalf("expected two products on page, got %d", len(result.Products))
	}
	expected := []string{"p4-excellent", "p4-good"}
	for index, product := range result.Products {
		if product.ID != expected[index] {
			t.Fatalf("expected product %d to be %s, got %s", index, expected[index], product.ID)
		}
	}
}

func TestRepositoryListFiltersByCategoryBrandAndCondition(t *testing.T) {
	repository := newTestRepository()

	result, err := repository.List(context.Background(), ListQuery{
		Categories: []string{"audio", "desktop"},
		Brands:     []string{"apple"},
		Conditions: []string{"excellent"},
		Page:       1,
		PageSize:   10,
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if result.Pagination.Total != 4 {
		t.Fatalf("expected four filtered products, got %d", result.Pagination.Total)
	}
	if len(result.Products) != 4 {
		t.Fatalf("expected four products on page, got %d", len(result.Products))
	}

	expected := []string{"p4-excellent", "p6-excellent", "p7-excellent", "p8-excellent"}
	for index, product := range result.Products {
		if product.ID != expected[index] {
			t.Fatalf("expected product %d to be %s, got %s", index, expected[index], product.ID)
		}
	}
}

func TestRepositoryListPaginatesResults(t *testing.T) {
	repository := newTestRepository()

	result, err := repository.List(context.Background(), ListQuery{Page: 2, PageSize: 5})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if result.Pagination.Total != 16 {
		t.Fatalf("expected total 16 products, got %d", result.Pagination.Total)
	}
	if !result.Pagination.HasMore {
		t.Fatal("expected page 2 to still have more results")
	}
	if len(result.Products) != 5 {
		t.Fatalf("expected 5 products on page 2, got %d", len(result.Products))
	}

	ids := []string{result.Products[0].ID, result.Products[1].ID, result.Products[2].ID, result.Products[3].ID, result.Products[4].ID}
	expected := []string{"p3-good", "p4-excellent", "p4-good", "p5-excellent", "p5-good"}
	for index, id := range ids {
		if id != expected[index] {
			t.Fatalf("expected product %d to be %s, got %s", index, expected[index], id)
		}
	}
}

func TestRepositoryListSortsByPopularityBeforePagination(t *testing.T) {
	repository := newTestRepository()

	result, err := repository.List(context.Background(), ListQuery{
		Sort:     SortOptionPopularity,
		Page:     1,
		PageSize: 6,
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	expected := []string{"p4-excellent", "p4-good", "p1-excellent", "p1-good", "p2-excellent", "p2-good"}
	if len(result.Products) != len(expected) {
		t.Fatalf("expected %d products on page, got %d", len(expected), len(result.Products))
	}

	for index, product := range result.Products {
		if product.ID != expected[index] {
			t.Fatalf("expected product %d to be %s, got %s", index, expected[index], product.ID)
		}
	}
}

func TestRepositoryListCachesMergedProductsForTTL(t *testing.T) {
	tempDir := t.TempDir()
	metadataPath := filepath.Join(tempDir, "metadata.json")
	detailsPath := filepath.Join(tempDir, "details.json")

	writeJSONFile(t, metadataPath, []ProductMetadata{
		{
			ID:        "p1",
			Name:      "Original Name",
			Brand:     "Apple",
			Category:  "smartphone",
			BasePrice: 100,
		},
	})
	writeJSONFile(t, detailsPath, []ProductVariant{
		{
			ID:              "p1-excellent",
			ProductID:       "p1",
			Condition:       "excellent",
			DiscountPercent: 10,
			Colors:          []string{"blue"},
			Stock:           2,
		},
	})

	now := time.Date(2026, time.July, 3, 12, 0, 0, 0, time.UTC)
	repository := NewRepositoryWithPaths(metadataPath, detailsPath)
	repository.now = func() time.Time { return now }

	firstResult, err := repository.List(context.Background(), ListQuery{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("first List returned error: %v", err)
	}
	if firstResult.Products[0].Name != "Original Name" {
		t.Fatalf("expected first cached name to be Original Name, got %s", firstResult.Products[0].Name)
	}

	writeJSONFile(t, metadataPath, []ProductMetadata{
		{
			ID:        "p1",
			Name:      "Updated Name",
			Brand:     "Apple",
			Category:  "smartphone",
			BasePrice: 100,
		},
	})

	secondResult, err := repository.List(context.Background(), ListQuery{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("second List returned error: %v", err)
	}
	if secondResult.Products[0].Name != "Original Name" {
		t.Fatalf("expected cached name before TTL expiry to remain Original Name, got %s", secondResult.Products[0].Name)
	}

	now = now.Add(repositoryCacheTTL + time.Second)

	thirdResult, err := repository.List(context.Background(), ListQuery{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("third List returned error: %v", err)
	}
	if thirdResult.Products[0].Name != "Updated Name" {
		t.Fatalf("expected refreshed name after TTL expiry to be Updated Name, got %s", thirdResult.Products[0].Name)
	}
}

func newTestRepository() *Repository {
	return NewRepositoryWithPaths(
		filepath.Join("..", "..", "data", "metadata.json"),
		filepath.Join("..", "..", "data", "details.json"),
	)
}

func writeJSONFile(t *testing.T, path string, value any) {
	t.Helper()

	content, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal test json: %v", err)
	}

	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write test json file: %v", err)
	}
}
