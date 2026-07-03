<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from 'vue'

type RangeSelection = {
	min: number
	max: number
}

type ThumbRole = 'min' | 'max'

type RangeSliderProps = {
	min?: number
	max?: number
	currentMin?: number
	currentMax?: number
	step?: number
}

const props = withDefaults(defineProps<RangeSliderProps>(), {
	min: 0,
	max: 100,
	step: 1,
})

const emit = defineEmits<{
	change: [selection: RangeSelection]
	'update:currentMin': [value: number]
	'update:currentMax': [value: number]
}>()

const selectedMin = ref<number>(0)
const selectedMax = ref<number>(0)
const activeThumb = ref<ThumbRole | null>(null)
const trackElement = ref<HTMLElement | null>(null)
const dragStartSelection = ref<RangeSelection | null>(null)

const minimumBound = computed(() => Math.min(props.min, props.max))
const maximumBound = computed(() => Math.max(props.min, props.max))
const thumbSizePx = 25
const thumbRadiusPx = thumbSizePx / 2

const clamp = (value: number, minimum: number, maximum: number) => Math.min(Math.max(value, minimum), maximum)
const getSpan = () => maximumBound.value - minimumBound.value || 1

const applyDraggedValue = (thumb: ThumbRole, rawValue: number) => {
	const nextValue = clamp(rawValue, minimumBound.value, maximumBound.value)

	if (thumb === 'min') {
		if (nextValue <= selectedMax.value) {
			selectedMin.value = nextValue
			return 'min' as const
		}

		const previousMax = selectedMax.value
		selectedMin.value = previousMax
		selectedMax.value = nextValue
		return 'max' as const
	}

	if (nextValue >= selectedMin.value) {
		selectedMax.value = nextValue
		return 'max' as const
	}

	const previousMin = selectedMin.value
	selectedMax.value = previousMin
	selectedMin.value = nextValue
	return 'min' as const
}

const syncFromProps = () => {
	const rawMin = props.currentMin ?? minimumBound.value
	const rawMax = props.currentMax ?? maximumBound.value
	const nextMin = clamp(Math.min(rawMin, rawMax), minimumBound.value, maximumBound.value)
	const nextMax = clamp(Math.max(rawMin, rawMax), nextMin, maximumBound.value)

	selectedMin.value = nextMin
	selectedMax.value = nextMax
}

watch(() => [props.min, props.max, props.currentMin, props.currentMax], syncFromProps, {
	immediate: true,
})

const emitSelection = () => {
	const selection: RangeSelection = {
		min: selectedMin.value,
		max: selectedMax.value,
	}

	emit('update:currentMin', selection.min)
	emit('update:currentMax', selection.max)
	emit('change', selection)
}

const getInputValue = (event: Event) => {
	if (!(event.target instanceof HTMLInputElement)) {
		return null
	}

	return Number(event.target.value)
}

const hasDraggedSelectionChanged = () =>
	dragStartSelection.value !== null
	&& (
		dragStartSelection.value.min !== selectedMin.value
		|| dragStartSelection.value.max !== selectedMax.value
	)

const onMinInput = (event: Event) => {
	const value = getInputValue(event)

	if (value === null) {
		return
	}

	activeThumb.value = applyDraggedValue('min', value)
	dragStartSelection.value = null
	}

const onMinChange = () => {
	emitSelection()
}

const onMaxInput = (event: Event) => {
	const value = getInputValue(event)

	if (value === null) {
		return
	}

	activeThumb.value = applyDraggedValue('max', value)
	dragStartSelection.value = null
	}

const onMaxChange = () => {
	emitSelection()
}

const resolveValueFromClientX = (clientX: number) => {
	const track = trackElement.value

	if (!track) {
		return null
	}

	const trackBounds = track.getBoundingClientRect()
	const ratio = clamp((clientX - trackBounds.left) / trackBounds.width, 0, 1)
	const rawValue = minimumBound.value + ratio * getSpan()
	const stepSize = props.step > 0 ? props.step : 1
	const steppedValue = minimumBound.value + Math.round((rawValue - minimumBound.value) / stepSize) * stepSize

	return clamp(steppedValue, minimumBound.value, maximumBound.value)
}

const stopDragging = () => {
	const shouldEmitSelection = hasDraggedSelectionChanged()

	activeThumb.value = null
	dragStartSelection.value = null
	window.removeEventListener('pointermove', onPointerMove)
	window.removeEventListener('pointerup', stopDragging)
	window.removeEventListener('pointercancel', stopDragging)

	if (shouldEmitSelection) {
		emitSelection()
	}
}

const onPointerMove = (event: PointerEvent) => {
	if (activeThumb.value === null) {
		return
	}

	const nextValue = resolveValueFromClientX(event.clientX)

	if (nextValue === null) {
		return
	}

	activeThumb.value = applyDraggedValue(activeThumb.value, nextValue)
	event.preventDefault()
}

const startDragging = (thumb: ThumbRole, event: PointerEvent) => {
	stopDragging()
	dragStartSelection.value = {
		min: selectedMin.value,
		max: selectedMax.value,
	}
	activeThumb.value = thumb
	window.addEventListener('pointermove', onPointerMove)
	window.addEventListener('pointerup', stopDragging)
	window.addEventListener('pointercancel', stopDragging)
	onPointerMove(event)
	event.preventDefault()
}

const onTrackPointerDown = (event: PointerEvent) => {
	const nextValue = resolveValueFromClientX(event.clientX)

	if (nextValue === null) {
		return
	}

	const nextThumb: ThumbRole = Math.abs(nextValue - selectedMin.value) <= Math.abs(nextValue - selectedMax.value)
		? 'min'
		: 'max'

	startDragging(nextThumb, event)
}

const priceFormatter = new Intl.NumberFormat('en-IE', {
	style: 'currency',
	currency: 'EUR',
	minimumFractionDigits: 0,
	maximumFractionDigits: 0,
})

const formattedMinimum = computed(() => priceFormatter.format(selectedMin.value))
const formattedMaximum = computed(() => priceFormatter.format(selectedMax.value))
const getThumbOffset = (value: number) => {
	const percentage = ((value - minimumBound.value) / getSpan()) * 100

	return `clamp(0px, calc(${percentage}% - ${thumbRadiusPx}px), calc(100% - ${thumbSizePx}px))`
}

const minimumThumbStyle = computed(() => ({
	left: getThumbOffset(selectedMin.value),
}))

const maximumThumbStyle = computed(() => ({
	left: getThumbOffset(selectedMax.value),
}))

const minimumThumbClass = computed(() => activeThumb.value === 'min' ? 'slider-thumb--active' : 'slider-thumb--inactive')
const maximumThumbClass = computed(() => activeThumb.value === 'max' || activeThumb.value === null ? 'slider-thumb--active' : 'slider-thumb--inactive')
const sliderInstructionsId = 'price-range-instructions'

const rangeProgressStyle = computed(() => {
	const span = getSpan()
	const left = ((selectedMin.value - minimumBound.value) / span) * 100
	const right = ((maximumBound.value - selectedMax.value) / span) * 100

	return {
		left: `clamp(0px, calc(${left}% - ${thumbRadiusPx}px), calc(100% - ${thumbSizePx}px))`,
		right: `clamp(0px, calc(${right}% - ${thumbRadiusPx}px), calc(100% - ${thumbSizePx}px))`,
	}
})

onUnmounted(stopDragging)
</script>

<template>
	<fieldset class="flex flex-col gap-[10px]">
		<legend class="text-[20px] font-medium leading-none tracking-[-0.03em] text-slate-950">Price Range</legend>

		<div ref="trackElement" class="price-range-track relative h-[25px] touch-none mt-4" @pointerdown="onTrackPointerDown">
			<div class="absolute inset-0 rounded-full bg-[#ECECF2]"></div>
			<div class="absolute inset-y-0 rounded-full bg-[#05031A]" :style="rangeProgressStyle"></div>

			<button
				type="button"
				aria-hidden="true"
				tabindex="-1"
				class="slider-thumb"
				:class="minimumThumbClass"
				:style="minimumThumbStyle"
				@pointerdown.stop="startDragging('min', $event)"
			></button>

			<button
				type="button"
				aria-hidden="true"
				tabindex="-1"
				class="slider-thumb"
				:class="maximumThumbClass"
				:style="maximumThumbStyle"
				@pointerdown.stop="startDragging('max', $event)"
			></button>

			<div class="sr-only">
				<p :id="sliderInstructionsId">Use arrow keys to adjust the minimum and maximum price range.</p>

				<label>
					Minimum price
					<input
						type="range"
						:min="minimumBound"
						:max="maximumBound"
						:step="step"
						:value="selectedMin"
						aria-label="Minimum price"
						:aria-describedby="sliderInstructionsId"
						:aria-valuetext="formattedMinimum"
						@input="onMinInput"
						@change="onMinChange"
					>
				</label>

				<label>
					Maximum price
					<input
						type="range"
						:min="minimumBound"
						:max="maximumBound"
						:step="step"
						:value="selectedMax"
						aria-label="Maximum price"
						:aria-describedby="sliderInstructionsId"
						:aria-valuetext="formattedMaximum"
						@input="onMaxInput"
						@change="onMaxChange"
					>
				</label>
			</div>
		</div>

		<div class="grid grid-cols-2 items-end gap-4 text-[#556074]">
			<div class="space-y-[6px]">
				<p class="text-[14px] leading-none">Minimum:</p>
				<p class="text-[18px] font-normal leading-none tracking-[-0.04em]">{{ formattedMinimum }}</p>
			</div>

			<div class="space-y-[6px] text-right">
				<p class="text-[14px] leading-none">Maximum:</p>
				<p class="text-[18px] font-normal leading-none tracking-[-0.04em]">{{ formattedMaximum }}</p>
			</div>
		</div>
	</fieldset>
</template>

<style scoped>
.price-range-track:focus-within {
	outline: 2px solid #05031a;
	outline-offset: 6px;
}

.slider-thumb {
	position: absolute;
	top: 0;
	height: 25px;
	width: 25px;
	border: 3px solid #05031a;
	border-radius: 9999px;
	background: #ffffff;
	box-sizing: border-box;
	cursor: pointer;
	touch-action: none;
}

.slider-thumb--active {
	z-index: 2;
}

.slider-thumb--inactive {
	z-index: 1;
}

.sr-only input[type='range'] {
	position: absolute;
	clip: rect(0, 0, 0, 0);
	clip-path: inset(50%);
	height: 25px;
	width: 25px;
	overflow: hidden;
	white-space: nowrap;
}
</style>
