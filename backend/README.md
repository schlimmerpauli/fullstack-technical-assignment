# Backend

## Prerequisites
- Go 1.22 or higher

## Running the Server

The backend starts through the application entrypoint in `cmd/app`.

A `.env` file is optional. If no `.env` file is present, the backend falls back to the default host and port from the config package.

Optional environment overrides:

- `SERVER_HOST`
- `SERVER_PORT`

```bash
cd backend
go run ./cmd/app
```

The server will start on `http://localhost:8080`

## Testing

```bash
# Health check
curl http://localhost:8080/health

# Products endpoint
curl http://localhost:8080/products

# Search
curl "http://localhost:8080/products?search=iphone"

# Filter by color and bestseller
curl "http://localhost:8080/products?color=blue&bestseller=true"

# Price range and pagination
curl "http://localhost:8080/products?minPrice=100&maxPrice=500&page=1&pageSize=5"

# Load more via pagination
curl "http://localhost:8080/products?page=2&pageSize=5"
```

## Data Files

- `data/metadata.json` - Base product metadata (id, name, brand, category, base_price, image_url)
- `data/details.json` - Product condition variants (id, product_id, condition, discount_percent, bestseller, colors, stock)

## Architecture Notes

- The backend is organized around a `products` feature package under `internal/products`.
- The repository reads and merges the two JSON data sources into a single aggregated product list.
- The HTTP contract is intentionally listing-screen driven: keyword search, color, bestseller, and load-more pagination map directly to the product discovery UI.
- Filtering and load-more behavior are exposed through `GET /products` using query parameters.
- A simple in-memory cache with a 30 second TTL is applied to the full aggregated product list before filtering and pagination.

## Production Improvements

- Improve the package structure further by separating API DTOs/contracts from domain models and HTTP handlers.
- Replace file-backed storage with a proper database such as PostgreSQL for product and variant data.
- Replace the in-process cache with Redis for a production-ready distributed caching layer.
- Add a `Dockerfile.prod` for production builds.
- Add a `docker-compose` setup that can spin up backend and frontend dependencies quickly for local development.
- Add a CI/CD pipeline for validation, build, and deployment automation.
- Add a dev-container oriented workflow so backend and frontend development do not depend on local machine setup.

## Final Thoughts

The current implementation is intentionally lightweight for the assignment, but it is structured to make the next production-oriented steps straightforward.

