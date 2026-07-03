import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import ItemCard from './item-card.vue'

const product = {
	id: 'p1-excellent',
	name: 'iPhone 12',
	category: 'smartphone',
	brand: 'Apple',
	condition: 'excellent',
	price: 331.99,
	discountPercent: 20,
	bestseller: true,
	colors: ['blue', 'red', 'green'],
	imageUrl: 'https://example.com/iphone-12.webp',
	stock: 18,
}

describe('item-card', () => {
	it('emits open-detail when the card is clicked', async () => {
		const wrapper = mount(ItemCard, {
			props: {
				product,
			},
		})

		await wrapper.get('article').trigger('click')

		expect(wrapper.emitted('open-detail')?.[0]).toEqual([product])
	})

	it('does not emit open-detail when a color swatch is clicked', async () => {
		const wrapper = mount(ItemCard, {
			props: {
				product,
			},
		})

		await wrapper.get('button[aria-label="Select red color"]').trigger('click')

		expect(wrapper.emitted('open-detail')).toBeFalsy()
	})

	it('allows the first image to opt into eager loading', () => {
		const wrapper = mount(ItemCard, {
			props: {
				product,
				imageLoading: 'eager',
				imageFetchPriority: 'high',
			},
		})

		const image = wrapper.get('img')

		expect(image.attributes('loading')).toBe('eager')
		expect(image.attributes('fetchpriority')).toBe('high')
	})
})