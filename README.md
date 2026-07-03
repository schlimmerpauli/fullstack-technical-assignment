# Refurbed Fullstack Assignment Submission

This repository contains a Go backend and a Vue 3 + Tailwind frontend for the Refurbed product discovery assignment.

The submission is intentionally lightweight and timeboxed around the assignment scope. The main priorities were a working end-to-end product listing flow, clear separation of responsibilities, and straightforward local setup and validation.

## Repository Structure

```text
backend/                         Go API and data aggregation
assignment_vue/frontend-vue/     Vue 3 + Vite + Tailwind frontend
```

## Prerequisites

- Go 1.22 or higher
- Node.js 18 or higher

## How to Run the Backend

From the repository root:

```bash
cd backend
go run ./cmd/app
```

The backend starts on `http://localhost:8080` by default.

Optional environment overrides:

- `SERVER_HOST`
- `SERVER_PORT`
- `CORS_ALLOWED_ORIGINS`

## How to Run the Frontend

From the repository root:

```bash
cd assignment_vue/frontend-vue
npm install
npm run dev
```

The frontend starts on `http://localhost:5173`.

If the frontend needs to call a backend hosted on another origin, set:

```bash
VITE_API_BASE_URL=http://localhost:8080
```

An example env file is available at `assignment_vue/frontend-vue/.env.example`.

## Validation

Backend:

```bash
cd backend
go test ./...
```

Frontend:

```bash
cd assignment_vue/frontend-vue
npm test
npm run typecheck
npm run build
```

## API Overview

The backend exposes a single endpoint:

```text
GET /products
```

Supported query parameters:

- `search`
- repeated `category`
- repeated `brand`
- repeated `condition`
- `color`
- `bestseller`
- `minPrice`
- `maxPrice`
- `sort=popularity`
- `page`
- `pageSize`

Example requests:

```bash
curl "http://localhost:8080/products?search=iphone"
curl "http://localhost:8080/products?color=blue&bestseller=true"
curl "http://localhost:8080/products?sort=popularity"
curl "http://localhost:8080/products?minPrice=100&maxPrice=500&page=1&pageSize=5"
```

## Architecture Notes

- The backend is organized around a `products` feature package under `backend/internal/products`.
- Product data is aggregated from the JSON metadata and details sources into a listing-oriented API model.
- Search, filtering, popularity sorting, and pagination are applied server-side through the `/products` endpoint.
- A simple in-memory cache with a 30 second TTL is applied to the full aggregated product list before filters and pagination.
- The frontend is a single-page Vue application with page-level listing state in `src/App.vue` and reusable UI components under `src/components`.
- The frontend delegates search, filters, sorting, and pagination to the backend and derives filter options from the returned catalog data for the current assignment scope.

## Production Improvements

- Split backend transport DTOs from domain models and handlers more explicitly.
- Replace file-backed storage with a database-backed product source.
- Replace the in-process cache with a distributed cache such as Redis.
- Add a dedicated facets or filter-metadata endpoint so the frontend does not need to derive filter options from listing data.
- Persist frontend search and filter state in the URL for shareable and refresh-safe product views.
- Improve loading and empty-state UX with skeleton cards and more polished error handling.
- Add containerized local setup and CI for automated validation.

## Final Thoughts

This implementation is intentionally scoped to the assignment rather than a production deployment. The goal was to show clear structure, practical tradeoffs, and a complete fullstack slice that is easy to run, review, and extend.
