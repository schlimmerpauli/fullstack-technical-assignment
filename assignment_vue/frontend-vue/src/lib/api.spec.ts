import { describe, expect, it } from 'vitest'

import { buildProductsPath } from './api'

describe('api helpers', () => {
	it('adds the popularity sort query parameter when requested', () => {
		expect(
			buildProductsPath({
				sort: 'popularity',
				page: 2,
				pageSize: 6,
			}),
		).toBe('/products?sort=popularity&page=2&pageSize=6')
	})
})