import { flushPromises, mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import App from './App.vue'
import { listProducts } from './lib/api'
import type { ProductListResponse } from './lib/products'

vi.mock('./lib/api', () => ({
  listProducts: vi.fn(),
}))

const mockedListProducts = vi.mocked(listProducts)

const productListResponse: ProductListResponse = {
  products: [
    {
      id: 'speaker-1',
      name: 'Studio Speaker',
      category: 'audio',
      brand: 'refurbed',
      condition: 'excellent',
      price: 399,
      discountPercent: 10,
      bestseller: false,
      colors: ['black'],
      imageUrl: 'https://example.com/speaker.png',
      stock: 4,
    },
  ],
  pagination: {
    page: 1,
    pageSize: 100,
    total: 1,
    hasMore: false,
  },
}

const emptyProductListResponse: ProductListResponse = {
  products: [],
  pagination: {
    page: 1,
    pageSize: 6,
    total: 0,
    hasMore: false,
  },
}

describe('App', () => {
  beforeEach(() => {
    mockedListProducts.mockReset()
    mockedListProducts.mockResolvedValue(productListResponse)
  })

  afterEach(() => {
    document.body.innerHTML = ''
  })

  it('renders the empty results state without the sort control', async () => {
    mockedListProducts.mockResolvedValue(emptyProductListResponse)

    const wrapper = mount(App, {
      global: {
        stubs: {
          FiltersPanel: true,
          ItemCard: true,
          SearchInput: true,
        },
      },
    })

    await flushPromises()
    await nextTick()
    await flushPromises()

    expect(wrapper.text()).toContain('0 products found')
    expect(wrapper.text()).toContain('No products found')
    expect(wrapper.text()).toContain("We couldn't find any products matching your current filters. Try adjusting your search criteria.")
    expect(wrapper.find('#products-sort').exists()).toBe(false)
    expect(wrapper.findAll('button').some((button) => button.text() === 'Clear all filters')).toBe(true)
  })

  it('loads visible products before background filter metadata', async () => {
    mount(App, {
      global: {
        stubs: {
          FiltersPanel: true,
          ItemCard: true,
          SearchInput: true,
        },
      },
    })

    await flushPromises()

    expect(mockedListProducts.mock.calls[0]?.[0]).toEqual({
      page: 1,
      pageSize: 6,
    })
    expect(mockedListProducts.mock.calls[1]?.[0]).toEqual({
      page: 1,
      pageSize: 100,
    })
  })

  it('moves focus back to the filters trigger before hiding the mobile panel', async () => {
    const wrapper = mount(App, {
      attachTo: document.body,
      global: {
        stubs: {
          FiltersPanel: true,
          ItemCard: true,
          SearchInput: true,
        },
      },
    })

    await flushPromises()

    const filtersTrigger = wrapper.get('button[aria-controls="mobile-filters-panel"]')

    await filtersTrigger.trigger('click')
    await nextTick()

    const closeButton = wrapper.get('#mobile-filters-panel button[aria-label="Close filters"]')

    expect(document.activeElement).toBe(closeButton.element)

    await closeButton.trigger('click')
    await nextTick()

    expect(document.activeElement).toBe(filtersTrigger.element)
    expect(wrapper.get('#mobile-filters-panel').attributes('aria-hidden')).toBe('true')

    wrapper.unmount()
  })
})