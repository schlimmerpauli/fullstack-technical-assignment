<template>
	<div class="space-y-8 sm:space-y-10">
		<MultiSelect
			:model-value="props.selectedCategories"
			title="Category"
			:options="props.categoryOptions"
			@update:model-value="emit('update:selectedCategories', $event)"
		/>

		<MultiSelect
			:model-value="props.selectedBrands"
			title="Brand"
			:options="props.brandOptions"
			@update:model-value="emit('update:selectedBrands', $event)"
		/>

		<MultiSelect
			:model-value="props.selectedConditions"
			title="Condition"
			:options="props.conditionOptions"
			@update:model-value="emit('update:selectedConditions', $event)"
		/>

		<MultiSelect
			:model-value="selectedColorValues"
			title="Color"
			:options="props.colorOptions"
			selection-mode="single"
			@update:model-value="onColorChange"
		/>
		<RangeSlider
			:max="props.priceUpperBound"
			:current-min="props.minPrice"
			:current-max="props.maxPrice"
			@update:current-min="emit('update:minPrice', $event)"
			@update:current-max="emit('update:maxPrice', $event)"
		/>
	</div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import MultiSelect from '../multi-select/multi-select.vue'
import RangeSlider from '../range-slider/range-slider.vue'

type FilterOption = {
	label: string
	value: string
}

type FiltersPanelProps = {
	selectedCategories?: string[]
	selectedBrands?: string[]
	selectedConditions?: string[]
	selectedColor?: string | null
	minPrice?: number
	maxPrice?: number
	priceUpperBound?: number
	categoryOptions?: FilterOption[]
	brandOptions?: FilterOption[]
	conditionOptions?: FilterOption[]
	colorOptions?: FilterOption[]
}

const props = withDefaults(defineProps<FiltersPanelProps>(), {
	selectedCategories: () => [],
	selectedBrands: () => [],
	selectedConditions: () => [],
	selectedColor: null,
	minPrice: 0,
	maxPrice: 100,
	priceUpperBound: 100,
	categoryOptions: () => [],
	brandOptions: () => [],
	conditionOptions: () => [],
	colorOptions: () => [],
})

const emit = defineEmits<{
	'update:selectedCategories': [values: string[]]
	'update:selectedBrands': [values: string[]]
	'update:selectedConditions': [values: string[]]
	'update:selectedColor': [value: string | null]
	'update:minPrice': [value: number]
	'update:maxPrice': [value: number]
}>()

const selectedColorValues = computed(() => props.selectedColor === null ? [] : [props.selectedColor])

const onColorChange = (values: string[]) => {
	emit('update:selectedColor', values[0] ?? null)
}
</script>
