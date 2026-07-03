<template>
	<label
		class="flex h-[70px] w-full max-w-[850px] items-center gap-4 rounded-[12px] bg-[#F3F3F5] px-5 text-[#969DB0] sm:px-6"
	>
		<span class="sr-only">Search products</span>

		<svg
			class="h-8 w-8 shrink-0"
			viewBox="0 0 48 48"
			fill="none"
			aria-hidden="true"
		>
			<circle cx="22" cy="22" r="13" stroke="currentColor" stroke-width="3.5" />
			<path d="M31 31L39 39" stroke="currentColor" stroke-width="3.5" stroke-linecap="round" />
		</svg>

		<input
			type="search"
			name="search"
			:value="localValue"
			placeholder="Search products..."
			aria-label="Search products"
			enterkeyhint="search"
			spellcheck="false"
			class="w-full border-0 bg-transparent p-0 text-[16px] leading-[16px] text-[#6F7487] placeholder:text-[#7E8496] focus:outline-none focus:ring-0"
			@input="onInput"
		/>
	</label>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue'

type SearchInputProps = {
	modelValue?: string
	debounceMs?: number
}

const props = withDefaults(defineProps<SearchInputProps>(), {
	modelValue: '',
	debounceMs: 300,
})

const emit = defineEmits<{
	search: [value: string]
	'update:modelValue': [value: string]
}>()

const localValue = ref(props.modelValue)

let searchTimeoutId: ReturnType<typeof window.setTimeout> | null = null

const clearPendingSearch = () => {
	if (searchTimeoutId === null) {
		return
	}

	window.clearTimeout(searchTimeoutId)
	searchTimeoutId = null
}

watch(
	() => props.modelValue,
	(value) => {
		if (value === localValue.value) {
			return
		}

		localValue.value = value
		clearPendingSearch()
	},
)

const queueSearch = () => {
	clearPendingSearch()

	searchTimeoutId = window.setTimeout(() => {
		emit('search', localValue.value)
		searchTimeoutId = null
	}, props.debounceMs)
}

const onInput = (event: Event) => {
	if (!(event.target instanceof HTMLInputElement)) {
		return
	}

	localValue.value = event.target.value
	emit('update:modelValue', localValue.value)
	queueSearch()
}

onBeforeUnmount(clearPendingSearch)
</script>
