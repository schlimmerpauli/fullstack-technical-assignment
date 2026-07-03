export interface Product {
	id: string
	name: string
	category: string
	brand: string
	condition: string
	price: number
	discountPercent: number
	bestseller: boolean
	colors: string[]
	imageUrl: string
	stock: number
}

export type ProductSort = 'popularity'

export interface ProductListQuery {
	search?: string
	categories?: string[]
	brands?: string[]
	conditions?: string[]
	color?: string
	bestseller?: boolean
	sort?: ProductSort
	minPrice?: number
	maxPrice?: number
	page?: number
	pageSize?: number
}

export interface Pagination {
	page: number
	pageSize: number
	total: number
	hasMore: boolean
}

export interface ProductListResponse {
	products: Product[]
	pagination: Pagination
}