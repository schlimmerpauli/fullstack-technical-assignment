import { mount } from '@vue/test-utils'
import { afterEach, describe, expect, it, vi } from 'vitest'

import SearchInput from './search-input.vue'

describe('search-input', () => {
	afterEach(() => {
		vi.useRealTimers()
	})

	it('emits a debounced search event after typing settles', async () => {
		vi.useFakeTimers()

		const wrapper = mount(SearchInput)
		const input = wrapper.get('input[type="search"]')

		await input.setValue('iphone')

		expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['iphone'])
		expect(wrapper.emitted('search')).toBeFalsy()

		vi.advanceTimersByTime(299)
		expect(wrapper.emitted('search')).toBeFalsy()

		vi.advanceTimersByTime(1)
		expect(wrapper.emitted('search')?.[0]).toEqual(['iphone'])
	})

	it('restarts the debounce window on each keystroke and only emits the latest value', async () => {
		vi.useFakeTimers()

		const wrapper = mount(SearchInput)
		const input = wrapper.get('input[type="search"]')

		await input.setValue('iph')
		vi.advanceTimersByTime(200)
		await input.setValue('iphone')

		vi.advanceTimersByTime(299)
		expect(wrapper.emitted('search')).toBeFalsy()

		vi.advanceTimersByTime(1)
		const searchEvents = wrapper.emitted('search')

		expect(searchEvents).toHaveLength(1)
		expect(searchEvents?.[0]).toEqual(['iphone'])
	})

	it('keeps the pending debounce when the parent echoes the same modelValue back', async () => {
		vi.useFakeTimers()

		const wrapper = mount(SearchInput, {
			props: {
				modelValue: '',
			},
		})
		const input = wrapper.get('input[type="search"]')

		await input.setValue('mac')
		await wrapper.setProps({
			modelValue: 'mac',
		})

		vi.advanceTimersByTime(300)

		expect(wrapper.emitted('search')?.[0]).toEqual(['mac'])
	})

	it('syncs the input value when the parent updates modelValue', async () => {
		const wrapper = mount(SearchInput, {
			props: {
				modelValue: 'apple',
			},
		})

		await wrapper.setProps({
			modelValue: 'google',
		})

		const input = wrapper.get('input[type="search"]').element as HTMLInputElement

		expect(input.value).toBe('google')
	})
})