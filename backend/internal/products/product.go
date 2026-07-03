package products

type ListQuery struct {
	Search     string      `json:"search,omitempty"`
	Color      string      `json:"color,omitempty"`
	Bestseller *bool       `json:"bestseller,omitempty"`
	PriceRange *PriceRange `json:"priceRange,omitempty"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
}

type PriceRange struct {
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
}

type Pagination struct {
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
	Total    int  `json:"total"`
	HasMore  bool `json:"hasMore"`
}

type ProductMetadata struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	BasePrice float64 `json:"base_price"`
	ImageURL  string  `json:"image_url"`
	Category  string  `json:"category,omitempty"`
	Brand     string  `json:"brand,omitempty"`
}

type ProductVariant struct {
	ID              string   `json:"id"`
	ProductID       string   `json:"product_id"`
	Condition       string   `json:"condition"`
	DiscountPercent int      `json:"discount_percent"`
	Bestseller      bool     `json:"bestseller"`
	Colors          []string `json:"colors"`
	Stock           int      `json:"stock"`
}

type Product struct {
	ID              string   `json:"id"`
	ProductID       string   `json:"-"`
	Name            string   `json:"name"`
	Category        string   `json:"-"`
	Brand           string   `json:"-"`
	Condition       string   `json:"-"`
	Price           float64  `json:"price"`
	DiscountPercent int      `json:"discount_percent"`
	Bestseller      bool     `json:"bestseller"`
	Colors          []string `json:"colors"`
	ImageURL        string   `json:"image_url"`
	Stock           int      `json:"stock"`
}

type ListResult struct {
	Products   []Product
	Pagination Pagination
}

// usually this should land in a contracts or data transfer objects (dto) package
type ListResponse struct {
	Products   []Product  `json:"products"`
	Pagination Pagination `json:"pagination"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
