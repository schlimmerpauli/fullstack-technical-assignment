import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import RangeSlider from './range-slider.vue'

const getSlider = (index: number, sliders: ReturnType<typeof mount> extends infer T ? T extends { findAll: (...args: never[]) => infer R } ? R : never : never) => {
  const slider = sliders[index]

  if (!slider) {
    throw new Error(`Expected slider at index ${index}`)
  }

  return slider
}

describe('range-slider', () => {
  it('updates the local minimum on input and only emits after change', async () => {
    const wrapper = mount(RangeSlider, {
      props: {
        min: 0,
        max: 100,
        currentMin: 20,
        currentMax: 80,
      },
    })

    const minimumSlider = getSlider(0, wrapper.findAll('input[type="range"]'))
    const minimumInput = minimumSlider.element as HTMLInputElement

    minimumInput.value = '30'
    await minimumSlider.trigger('input')

    expect(wrapper.emitted('change')).toBeFalsy()
    expect(minimumInput.value).toBe('30')

    await minimumSlider.trigger('change')

    expect(wrapper.emitted('change')?.[0]).toEqual([{ min: 30, max: 80 }])
  })

  it('emits the selected min and max when the minimum changes', async () => {
    const wrapper = mount(RangeSlider, {
      props: {
        min: 0,
        max: 100,
        currentMin: 20,
        currentMax: 80,
      },
    })

    const minimumSlider = getSlider(0, wrapper.findAll('input[type="range"]'))

    await minimumSlider.setValue(30)

    const changeEvents = wrapper.emitted('change')

    expect(changeEvents).toBeTruthy()
    expect(changeEvents?.[changeEvents.length - 1]).toEqual([{ min: 30, max: 80 }])
  })

  it('swaps the range bounds when the maximum moves below the selected minimum', async () => {
    const wrapper = mount(RangeSlider, {
      props: {
        min: 0,
        max: 100,
        currentMin: 40,
        currentMax: 90,
      },
    })

    const maximumSlider = getSlider(1, wrapper.findAll('input[type="range"]'))

    await maximumSlider.setValue(20)

    const changeEvents = wrapper.emitted('change')
    const minimumSlider = getSlider(0, wrapper.findAll('input[type="range"]'))
    const minimumInput = minimumSlider.element as HTMLInputElement
    const maximumInput = maximumSlider.element as HTMLInputElement

    expect(changeEvents?.[changeEvents.length - 1]).toEqual([{ min: 20, max: 40 }])
    expect(minimumInput.value).toBe('20')
    expect(maximumInput.value).toBe('40')
  })

  it('swaps the range bounds when the minimum moves above the selected maximum', async () => {
    const wrapper = mount(RangeSlider, {
      props: {
        min: 0,
        max: 100,
        currentMin: 20,
        currentMax: 60,
      },
    })

    const minimumSlider = getSlider(0, wrapper.findAll('input[type="range"]'))

    await minimumSlider.setValue(90)

    const changeEvents = wrapper.emitted('change')
    const maximumSlider = getSlider(1, wrapper.findAll('input[type="range"]'))
    const minimumInput = minimumSlider.element as HTMLInputElement
    const maximumInput = maximumSlider.element as HTMLInputElement

    expect(changeEvents?.[changeEvents.length - 1]).toEqual([{ min: 60, max: 90 }])
    expect(minimumInput.value).toBe('60')
    expect(maximumInput.value).toBe('90')
  })

  it('syncs the slider positions when the parent updates the current values', async () => {
    const wrapper = mount(RangeSlider, {
      props: {
        min: 0,
        max: 100,
        currentMin: 10,
        currentMax: 70,
      },
    })

    await wrapper.setProps({
      currentMin: 25,
      currentMax: 55,
    })

    const minimumSlider = getSlider(0, wrapper.findAll('input[type="range"]'))
    const maximumSlider = getSlider(1, wrapper.findAll('input[type="range"]'))
    const minimumInput = minimumSlider.element as HTMLInputElement
    const maximumInput = maximumSlider.element as HTMLInputElement

    expect(minimumInput.value).toBe('25')
    expect(maximumInput.value).toBe('55')
  })

  it('keeps the decorative thumb buttons out of the tab order', () => {
    const wrapper = mount(RangeSlider, {
      props: {
        min: 0,
        max: 100,
        currentMin: 10,
        currentMax: 70,
      },
    })

    const thumbButtons = wrapper.findAll('button[aria-hidden="true"]')

    expect(thumbButtons).toHaveLength(2)
    expect(thumbButtons.every((button) => button.attributes('tabindex') === '-1')).toBe(true)
  })
})