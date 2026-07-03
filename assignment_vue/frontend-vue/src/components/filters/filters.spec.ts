import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import FiltersPanel from './filters.vue'
import MultiSelect from '../multi-select/multi-select.vue'
import RangeSlider from '../range-slider/range-slider.vue'

const filterProps = {
	selectedCategories: ['audio'],
	selectedBrands: ['apple'],
	selectedConditions: ['good'],
	selectedColor: 'blue',
	minPrice: 100,
	maxPrice: 500,
	priceUpperBound: 1450,
	categoryOptions: [
		{ label: 'Audio', value: 'audio' },
		{ label: 'Desktop', value: 'desktop' },
	],
	brandOptions: [
		{ label: 'Apple', value: 'apple' },
	],
	conditionOptions: [
		{ label: 'Excellent', value: 'excellent' },
		{ label: 'Good', value: 'good' },
	],
	colorOptions: [
		{ label: 'Blue', value: 'blue' },
		{ label: 'Red', value: 'red' },
	],
}

describe('filters-panel', () => {
	it('passes the current filter state down to the child controls', () => {
		const wrapper = mount(FiltersPanel, {
			props: filterProps,
		})

		const multiSelects = wrapper.findAllComponents(MultiSelect)
		const rangeSlider = wrapper.getComponent(RangeSlider)

		expect(multiSelects).toHaveLength(4)
		expect(multiSelects[0]?.props('title')).toBe('Category')
		expect(multiSelects[0]?.props('modelValue')).toEqual(['audio'])
		expect(multiSelects[1]?.props('title')).toBe('Brand')
		expect(multiSelects[1]?.props('modelValue')).toEqual(['apple'])
		expect(multiSelects[2]?.props('title')).toBe('Condition')
		expect(multiSelects[2]?.props('modelValue')).toEqual(['good'])
		expect(multiSelects[3]?.props('title')).toBe('Color')
		expect(multiSelects[3]?.props('modelValue')).toEqual(['blue'])
		expect(multiSelects[3]?.props('selectionMode')).toBe('single')

		expect(rangeSlider.props('currentMin')).toBe(100)
		expect(rangeSlider.props('currentMax')).toBe(500)
		expect(rangeSlider.props('max')).toBe(1450)
	})

	it('forwards category, brand, and condition updates', async () => {
		const wrapper = mount(FiltersPanel, {
			props: filterProps,
		})

		const multiSelects = wrapper.findAllComponents(MultiSelect)

		multiSelects[0]?.vm.$emit('update:modelValue', ['desktop'])
		multiSelects[1]?.vm.$emit('update:modelValue', ['apple', 'refurbed'])
		multiSelects[2]?.vm.$emit('update:modelValue', ['excellent'])

		await wrapper.vm.$nextTick()

		expect(wrapper.emitted('update:selectedCategories')?.[0]).toEqual([['desktop']])
		expect(wrapper.emitted('update:selectedBrands')?.[0]).toEqual([['apple', 'refurbed']])
		expect(wrapper.emitted('update:selectedConditions')?.[0]).toEqual([['excellent']])
	})

	it('maps color multi-select updates to a single selectedColor value', async () => {
		const wrapper = mount(FiltersPanel, {
			props: filterProps,
		})

		const colorSelect = wrapper.findAllComponents(MultiSelect)[3]

		colorSelect?.vm.$emit('update:modelValue', ['red'])
		colorSelect?.vm.$emit('update:modelValue', [])

		await wrapper.vm.$nextTick()

		expect(wrapper.emitted('update:selectedColor')?.[0]).toEqual(['red'])
		expect(wrapper.emitted('update:selectedColor')?.[1]).toEqual([null])
	})

	it('forwards price range updates from the slider', async () => {
		const wrapper = mount(FiltersPanel, {
			props: filterProps,
		})

		const rangeSlider = wrapper.getComponent(RangeSlider)

		rangeSlider.vm.$emit('update:currentMin', 150)
		rangeSlider.vm.$emit('update:currentMax', 450)

		await wrapper.vm.$nextTick()

		expect(wrapper.emitted('update:minPrice')?.[0]).toEqual([150])
		expect(wrapper.emitted('update:maxPrice')?.[0]).toEqual([450])
	})
})