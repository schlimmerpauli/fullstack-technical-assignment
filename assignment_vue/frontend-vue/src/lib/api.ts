import type { Pagination, Product, ProductListQuery, ProductListResponse } from './products'

const rawApiBaseUrl = import.meta.env.VITE_API_BASE_URL?.trim() ?? ''

type ApiProduct = {
	id: string
	name: string
	category: string
	brand: string
	condition: string
	price: number
	discount_percent: number
	bestseller: boolean
	colors: string[]
	image_url: string
	stock: number
}

type ApiProductListResponse = {
	products: ApiProduct[]
	pagination: Pagination
}

function mapProduct(apiProduct: ApiProduct): Product {
	return {
		id: apiProduct.id,
		name: apiProduct.name,
		category: apiProduct.category,
		brand: apiProduct.brand,
		condition: apiProduct.condition,
		price: apiProduct.price,
		discountPercent: apiProduct.discount_percent,
		bestseller: apiProduct.bestseller,
		colors: apiProduct.colors,
		imageUrl: apiProduct.image_url,
		stock: apiProduct.stock,
	}
}

function appendMultiValueParams(params: URLSearchParams, key: string, values: string[] | undefined) {
	if (!values) {
		return
	}

	for (const rawValue of values) {
		const value = rawValue.trim()
		if (value === '') {
			continue
		}

		params.append(key, value)
	}
}

function normalizeBaseUrl(baseUrl: string): string {
	return baseUrl.replace(/\/$/, '')
}

export function buildApiUrl(path: string): string {
	const normalizedPath = path.startsWith('/') ? path : `/${path}`

	if (rawApiBaseUrl !== '') {
		return `${normalizeBaseUrl(rawApiBaseUrl)}${normalizedPath}`
	}

	return normalizedPath
}

export function buildProductsPath(query: ProductListQuery = {}): string {
	const params = new URLSearchParams()
	const trimmedSearch = query.search?.trim()
	const trimmedColor = query.color?.trim()

	appendMultiValueParams(params, 'category', query.categories)
	appendMultiValueParams(params, 'brand', query.brands)
	appendMultiValueParams(params, 'condition', query.conditions)

	if (trimmedSearch) {
		params.set('search', trimmedSearch)
	}

	if (trimmedColor) {
		params.set('color', trimmedColor)
	}

	if (typeof query.bestseller === 'boolean') {
		params.set('bestseller', String(query.bestseller))
	}

	if (query.sort) {
		params.set('sort', query.sort)
	}

	if (typeof query.minPrice === 'number') {
		params.set('minPrice', String(query.minPrice))
	}

	if (typeof query.maxPrice === 'number') {
		params.set('maxPrice', String(query.maxPrice))
	}

	if (typeof query.page === 'number') {
		params.set('page', String(query.page))
	}

	if (typeof query.pageSize === 'number') {
		params.set('pageSize', String(query.pageSize))
	}

	const queryString = params.toString()

	return queryString === '' ? '/products' : `/products?${queryString}`
}

export async function listProducts(query: ProductListQuery = {}, init?: RequestInit): Promise<ProductListResponse> {
	const response = await fetchJSON<ApiProductListResponse>(buildProductsPath(query), init)

	return {
		products: response.products.map(mapProduct),
		pagination: response.pagination,
	}
}

export async function fetchJSON<T>(path: string, init?: RequestInit): Promise<T> {
	const response = await fetch(buildApiUrl(path), {
		headers: {
			Accept: 'application/json',
			...(init?.headers ?? {}),
		},
		...init,
	})

	if (!response.ok) {
		throw new Error(`Failed to fetch ${path}: ${response.status} ${response.statusText}`)
	}

	return response.json() as Promise<T>
}