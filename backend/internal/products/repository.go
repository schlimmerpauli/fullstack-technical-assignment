package products

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const repositoryCacheTTL = 30 * time.Second

type Repository struct {
	metadataPath string
	detailsPath  string
	cacheTTL     time.Duration
	now          func() time.Time

	mu             sync.RWMutex
	cachedProducts []Product
	cachedAt       time.Time
	cacheFilled    bool
}

func NewRepository() *Repository {
	return NewRepositoryWithPaths(
		filepath.Join("data", "metadata.json"),
		filepath.Join("data", "details.json"),
	)
}

func NewRepositoryWithPaths(metadataPath string, detailsPath string) *Repository {
	return &Repository{
		metadataPath: metadataPath,
		detailsPath:  detailsPath,
		cacheTTL:     repositoryCacheTTL,
		now:          time.Now,
	}
}

func (r *Repository) List(ctx context.Context, query ListQuery) (ListResult, error) {
	products, err := r.loadProducts(ctx)
	if err != nil {
		return ListResult{}, err
	}

	filteredProducts := filterProducts(products, query)
	paginatedProducts := paginateProducts(filteredProducts, query.Page, query.PageSize)

	return ListResult{
		Products: paginatedProducts,
		Pagination: Pagination{
			Page:     query.Page,
			PageSize: query.PageSize,
			Total:    len(filteredProducts),
			HasMore:  query.Page*query.PageSize < len(filteredProducts),
		},
	}, nil
}

func (r *Repository) loadProducts(ctx context.Context) ([]Product, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	now := r.now()

	r.mu.RLock()
	if r.cacheFilled && now.Sub(r.cachedAt) < r.cacheTTL {
		products := cloneProducts(r.cachedProducts)
		r.mu.RUnlock()
		return products, nil
	}
	r.mu.RUnlock()

	metadata, err := r.loadMetadata()
	if err != nil {
		return nil, fmt.Errorf("read metadata: %w", err)
	}

	variants, err := r.loadDetails()
	if err != nil {
		return nil, fmt.Errorf("read details: %w", err)
	}

	products := mergeProducts(metadata, variants)
	cachedProducts := cloneProducts(products)

	r.mu.Lock()
	r.cachedProducts = cachedProducts
	r.cachedAt = now
	r.cacheFilled = true
	r.mu.Unlock()

	return products, nil
}

func (r *Repository) loadMetadata() ([]ProductMetadata, error) {
	var metadata []ProductMetadata
	if err := loadJSONFile(r.metadataPath, &metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func (r *Repository) loadDetails() ([]ProductVariant, error) {
	var variants []ProductVariant
	if err := loadJSONFile(r.detailsPath, &variants); err != nil {
		return nil, err
	}

	return variants, nil
}

func loadJSONFile(path string, destination any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, destination); err != nil {
		return err
	}

	return nil
}

func mergeProducts(metadata []ProductMetadata, variants []ProductVariant) []Product {
	variantsByProductID := make(map[string][]ProductVariant, len(variants))
	for _, variant := range variants {
		variantsByProductID[variant.ProductID] = append(variantsByProductID[variant.ProductID], variant)
	}

	products := make([]Product, 0, len(variants))
	for _, item := range metadata {
		productVariants, exists := variantsByProductID[item.ID]
		if !exists {
			continue
		}

		for _, variant := range productVariants {
			products = append(products, Product{
				ID:              variant.ID,
				ProductID:       item.ID,
				Name:            item.Name,
				Category:        item.Category,
				Brand:           item.Brand,
				Condition:       variant.Condition,
				Price:           discountedPrice(item.BasePrice, variant.DiscountPercent),
				DiscountPercent: variant.DiscountPercent,
				Bestseller:      variant.Bestseller,
				Colors:          append([]string(nil), variant.Colors...),
				ImageURL:        item.ImageURL,
				Stock:           variant.Stock,
			})
		}
	}

	return products
}

func discountedPrice(basePrice float64, discountPercent int) float64 {
	discounted := basePrice * (1 - float64(discountPercent)/100)
	return math.Round(discounted*100) / 100
}

func filterProducts(products []Product, query ListQuery) []Product {
	filtered := make([]Product, 0, len(products))
	for _, product := range products {
		if !matchesSearch(product, query.Search) {
			continue
		}
		if !matchesColor(product, query.Color) {
			continue
		}
		if !matchesBestseller(product, query.Bestseller) {
			continue
		}
		if !matchesPrice(product.Price, query.PriceRange) {
			continue
		}

		filtered = append(filtered, product)
	}

	return filtered
}

func matchesSearch(product Product, rawSearch string) bool {
	terms := strings.Fields(strings.ToLower(rawSearch))
	if len(terms) == 0 {
		return true
	}

	name := strings.ToLower(product.Name)
	for _, term := range terms {
		if !strings.Contains(name, term) {
			return false
		}
	}

	return true
}

func matchesColor(product Product, rawColor string) bool {
	color := strings.ToLower(strings.TrimSpace(rawColor))
	if color == "" {
		return true
	}

	for _, candidate := range product.Colors {
		if strings.ToLower(strings.TrimSpace(candidate)) == color {
			return true
		}
	}

	return false
}

func matchesBestseller(product Product, bestseller *bool) bool {
	if bestseller == nil {
		return true
	}

	return product.Bestseller == *bestseller
}

func matchesPrice(price float64, priceRange *PriceRange) bool {
	if priceRange == nil {
		return true
	}

	if priceRange.Min != nil && price < *priceRange.Min {
		return false
	}

	if priceRange.Max != nil && price > *priceRange.Max {
		return false
	}

	return true
}

func paginateProducts(products []Product, page int, pageSize int) []Product {
	start := (page - 1) * pageSize
	if start >= len(products) {
		return []Product{}
	}

	end := start + pageSize
	if end > len(products) {
		end = len(products)
	}

	return append([]Product(nil), products[start:end]...)
}

func cloneProducts(products []Product) []Product {
	cloned := make([]Product, len(products))
	for index, product := range products {
		cloned[index] = product
		cloned[index].Colors = append([]string(nil), product.Colors...)
	}

	return cloned
}
