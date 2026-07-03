import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import MultiSelect from './multi-select.vue'

const options = [
  { label: 'Apple', value: 'apple' },
  { label: 'Samsung', value: 'samsung' },
  { label: 'Google', value: 'google' },
  { label: 'Xiaomi', value: 'xiaomi' },
]

const getOption = (
  index: number,
  optionButtons: ReturnType<typeof mount> extends infer T ? T extends { findAll: (...args: never[]) => infer R } ? R : never : never,
) => {
  const optionButton = optionButtons[index]

  if (!optionButton) {
    throw new Error(`Expected option button at index ${index}`)
  }

  return optionButton
}

describe('multi-select', () => {
  it('emits the selected values when an unselected option is chosen', async () => {
    const wrapper = mount(MultiSelect, {
      props: {
        title: 'Brand',
        options,
        modelValue: ['apple', 'samsung'],
      },
    })

    const googleOption = getOption(2, wrapper.findAll('button[role="checkbox"]'))

    await googleOption.trigger('click')

    const changeEvents = wrapper.emitted('change')
    const updateEvents = wrapper.emitted('update:modelValue')

    expect(changeEvents).toBeTruthy()
    expect(updateEvents).toBeTruthy()
    expect(changeEvents?.[changeEvents.length - 1]).toEqual([['apple', 'samsung', 'google']])
    expect(updateEvents?.[updateEvents.length - 1]).toEqual([['apple', 'samsung', 'google']])

    await wrapper.setProps({
      modelValue: ['apple', 'samsung', 'google'],
    })

    expect(googleOption.attributes('aria-checked')).toBe('true')
  })

  it('emits a single selected value in single selection mode', async () => {
    const wrapper = mount(MultiSelect, {
      props: {
        title: 'Color',
        options,
        modelValue: ['apple'],
        selectionMode: 'single',
      },
    })

    const samsungOption = getOption(1, wrapper.findAll('button[role="radio"]'))

    await samsungOption.trigger('click')

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([['samsung']])
    expect(wrapper.emitted('change')?.[0]).toEqual([['samsung']])
  })

  it('exposes the option group with a shared legend', () => {
    const wrapper = mount(MultiSelect, {
      props: {
        title: 'Brand',
        options,
      },
    })

    expect(wrapper.get('fieldset').text()).toContain('Brand')
    expect(wrapper.get('[role="group"]').attributes('aria-labelledby')).toBeTruthy()
  })

  it('emits the remaining values when a selected option is cleared', async () => {
    const wrapper = mount(MultiSelect, {
      props: {
        title: 'Brand',
        options,
        modelValue: ['apple', 'samsung'],
      },
    })

    const samsungOption = getOption(1, wrapper.findAll('button[role="checkbox"]'))

    await samsungOption.trigger('click')

    const changeEvents = wrapper.emitted('change')
    const updateEvents = wrapper.emitted('update:modelValue')

    expect(changeEvents?.[changeEvents.length - 1]).toEqual([['apple']])
    expect(updateEvents?.[updateEvents.length - 1]).toEqual([['apple']])

    await wrapper.setProps({
      modelValue: ['apple'],
    })

    expect(samsungOption.attributes('aria-checked')).toBe('false')
  })

  it('syncs the checked state when the parent updates modelValue', async () => {
    const wrapper = mount(MultiSelect, {
      props: {
        title: 'Brand',
        options,
        modelValue: ['apple', 'samsung'],
      },
    })

    await wrapper.setProps({
      modelValue: ['google'],
    })

    const optionButtons = wrapper.findAll('button[role="checkbox"]')
    const appleOption = getOption(0, optionButtons)
    const samsungOption = getOption(1, optionButtons)
    const googleOption = getOption(2, optionButtons)

    expect(appleOption.attributes('aria-checked')).toBe('false')
    expect(samsungOption.attributes('aria-checked')).toBe('false')
    expect(googleOption.attributes('aria-checked')).toBe('true')
  })
})