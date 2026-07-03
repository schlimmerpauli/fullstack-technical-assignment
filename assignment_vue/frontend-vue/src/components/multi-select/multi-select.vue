<template>
	<fieldset class="flex flex-col gap-[14px]">
		<legend :id="legendId" class="text-[18px] font-medium leading-none tracking-[-0.03em] text-slate-950">{{ title }}</legend>

		<div
			class="flex flex-col gap-[6px] mt-4"
			:role="selectionMode === 'single' ? 'radiogroup' : 'group'"
			:aria-labelledby="legendId"
		>
			<button
				v-for="option in options"
				:key="option.value"
				type="button"
				class="flex min-h-[32px] w-full items-center gap-3 text-left"
				:role="selectionMode === 'single' ? 'radio' : 'checkbox'"
				:aria-checked="isSelected(option.value)"
				@click="toggleValue(option.value)"
			>
				<span
					class="flex h-[24px] w-[24px] shrink-0 items-center justify-center rounded-[8px] border text-white transition-colors duration-150"
					:class="isSelected(option.value) ? 'border-[#05031A] bg-[#05031A]' : 'border-[#E2E5EE] bg-[#FFFFFF]'"
				>
					<svg
						viewBox="0 0 20 20"
						fill="none"
						aria-hidden="true"
						class="h-4 w-4 transition-opacity duration-150"
						:class="isSelected(option.value) ? 'opacity-100' : 'opacity-0'"
					>
						<path
							d="M4.5 10.5L8 14L15.5 6.5"
							stroke="currentColor"
							stroke-width="2.25"
							stroke-linecap="round"
							stroke-linejoin="round"
						/>
					</svg>
				</span>

				<span class="text-[16px] font-normal leading-none text-slate-950">{{ option.label }}</span>
			</button>
		</div>
	</fieldset>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type MultiSelectOption = {
	label: string
	value: string
}

type SelectionMode = 'single' | 'multiple'

type MultiSelectProps = {
	title?: string
	options?: MultiSelectOption[]
	modelValue?: string[]
	selectionMode?: SelectionMode
}

const props = withDefaults(defineProps<MultiSelectProps>(), {
	title: 'Brand',
	options: () => [],
	modelValue: () => [],
	selectionMode: 'multiple',
})

const emit = defineEmits<{
	change: [values: string[]]
	'update:modelValue': [values: string[]]
}>()

const legendId = computed(() => `multi-select-${props.title.toLowerCase().replace(/[^a-z0-9]+/g, '-')}`)

const isSelected = (value: string) => props.modelValue.includes(value)

const toggleValue = (value: string) => {
	const nextValues = props.selectionMode === 'single'
		? (isSelected(value) ? [] : [value])
		: (isSelected(value)
			? props.modelValue.filter((entry) => entry !== value)
			: [...props.modelValue, value])

	emit('update:modelValue', nextValues)
	emit('change', nextValues)
}
</script>
