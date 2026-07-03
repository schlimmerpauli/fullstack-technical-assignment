package products

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/products", h.handleList)
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	query, err := parseListQuery(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.service.List(r.Context(), query)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to load products"})
		return
	}

	writeJSON(w, http.StatusOK, ListResponse{
		Products:   result.Products,
		Pagination: result.Pagination,
	})
}

func parseListQuery(r *http.Request) (ListQuery, error) {
	values := r.URL.Query()
	query := ListQuery{
		Search:     strings.TrimSpace(values.Get("search")),
		Categories: parseMultiValue(values["category"]),
		Brands:     parseMultiValue(values["brand"]),
		Conditions: parseMultiValue(values["condition"]),
		Color:      strings.TrimSpace(values.Get("color")),
	}

	sortOption, err := parseSortOption(values.Get("sort"))
	if err != nil {
		return ListQuery{}, err
	}
	query.Sort = sortOption

	bestseller, err := parseOptionalBool(values.Get("bestseller"), "bestseller")
	if err != nil {
		return ListQuery{}, err
	}
	query.Bestseller = bestseller

	page, err := parsePositiveInt(values.Get("page"), "page", defaultPage)
	if err != nil {
		return ListQuery{}, err
	}
	query.Page = page

	pageSize, err := parsePositiveInt(values.Get("pageSize"), "pageSize", defaultPageSize)
	if err != nil {
		return ListQuery{}, err
	}
	if pageSize > maxPageSize {
		return ListQuery{}, fmt.Errorf("invalid pageSize: must be less than or equal to %d", maxPageSize)
	}
	query.PageSize = pageSize

	minPrice, err := parseOptionalFloat(values.Get("minPrice"), "minPrice")
	if err != nil {
		return ListQuery{}, err
	}

	maxPrice, err := parseOptionalFloat(values.Get("maxPrice"), "maxPrice")
	if err != nil {
		return ListQuery{}, err
	}

	if minPrice != nil || maxPrice != nil {
		if minPrice != nil && maxPrice != nil && *minPrice > *maxPrice {
			return ListQuery{}, fmt.Errorf("invalid priceRange: minPrice cannot be greater than maxPrice")
		}

		query.PriceRange = &PriceRange{
			Min: minPrice,
			Max: maxPrice,
		}
	}

	return query, nil
}

func parseMultiValue(rawValues []string) []string {
	values := make([]string, 0, len(rawValues))
	for _, rawValue := range rawValues {
		value := strings.TrimSpace(rawValue)
		if value == "" {
			continue
		}

		values = append(values, value)
	}

	if len(values) == 0 {
		return nil
	}

	return values
}

func parseSortOption(raw string) (SortOption, error) {
	value := strings.TrimSpace(strings.ToLower(raw))
	if value == "" {
		return SortOptionDefault, nil
	}

	if SortOption(value) == SortOptionPopularity {
		return SortOptionPopularity, nil
	}

	return SortOptionDefault, fmt.Errorf("invalid sort: must be popularity")
}

func parseOptionalBool(raw string, fieldName string) (*bool, error) {
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid %s: must be true or false", fieldName)
	}

	return &value, nil
}

func parsePositiveInt(raw string, fieldName string, defaultValue int) (int, error) {
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be an integer", fieldName)
	}
	if value < 1 {
		return 0, fmt.Errorf("invalid %s: must be greater than 0", fieldName)
	}

	return value, nil
}

func parseOptionalFloat(raw string, fieldName string) (*float64, error) {
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid %s: must be a number", fieldName)
	}

	return &value, nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
