<template>
	<article
		class="relative flex aspect-[340/510] w-full max-w-[340px] min-w-0 cursor-pointer flex-col overflow-visible rounded-[12px] p-4 ring-1 transition-[transform,box-shadow,background-color,ring-color] duration-200 ease-out sm:p-5"
		:class="[
			isCardSelected
				? 'bg-[#FFFEFA] -translate-y-1.5 shadow-[0_28px_56px_rgba(15,23,42,0.18)] ring-[#E6DED3]'
				: 'bg-white shadow-[0_20px_40px_rgba(15,23,42,0.14)] ring-slate-950/5',
			'hover:-translate-y-1.5 hover:shadow-[0_28px_56px_rgba(15,23,42,0.18)] hover:ring-[#E6DED3] focus-within:-translate-y-1.5 focus-within:shadow-[0_28px_56px_rgba(15,23,42,0.18)] focus-within:ring-[#E6DED3]',
		]"
		@click="handleCardClick"
	>
		<span
			v-if="product.bestseller"
			class="absolute left-0 top-0 z-10 bg-[#0D7B67] px-2.5 py-1 text-[10px] font-semibold uppercase tracking-[0.08em] text-white shadow-[0_8px_16px_rgba(13,123,103,0.24)] sm:px-3 sm:text-[11px]"
		>
			Bestseller
		</span>

		<span
			v-if="hasDiscount"
			class="absolute right-0 top-0 z-10 inline-flex h-[54px] w-[54px] translate-x-1/2 -translate-y-1/2 rotate-[25deg] items-center justify-center rounded-full bg-[#5B46DB] text-[12px] font-semibold text-white shadow-[0_14px_24px_rgba(79,59,194,0.28)] sm:h-[62px] sm:w-[62px] sm:text-sm"
		>
			-{{ product.discountPercent }}%
		</span>

		<div class="relative isolate flex h-[54.5%] items-center justify-center rounded-[18px] bg-[#F6F6F2] px-4 pb-4 pt-5 sm:px-6 sm:pb-5 sm:pt-6">
			<img
				:src="product.imageUrl"
				:alt="product.name"
				:fetchpriority="props.imageFetchPriority"
				:loading="props.imageLoading"
				class="max-h-[78%] w-auto max-w-full object-contain mix-blend-multiply drop-shadow-[0_18px_26px_rgba(15,23,42,0.12)]"
			>
		</div>

		<div class="mt-5 flex flex-1 flex-col">
			<h3 class="min-h-[48px] text-[22px] font-semibold leading-[1.05] tracking-[-0.04em] text-slate-950 sm:min-h-[56px] sm:text-[24px] 2xl:min-h-[62px] 2xl:text-[26px]">
				{{ product.name }}
			</h3>

			<ul
				v-if="availableColors.length > 0"
				class="mt-4 flex items-center gap-2.5"
				aria-label="Available colors"
			>
				<li
					v-for="color in availableColors"
					:key="`${product.id}-${color.label}`"
				>
					<button
						type="button"
						class="grid h-[18px] w-[18px] place-items-center rounded-full transition-[background-color,border-color,box-shadow] duration-200 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#05031A]/20 focus-visible:ring-offset-2"
						:class="selectedColorKey === color.key ? 'bg-[#05031A] p-[1px]' : ''"
						:aria-label="`Select ${color.label} color`"
						:aria-pressed="selectedColorKey === color.key"
						:title="color.label"
						@click.stop="selectSwatch(color.key)"
					>
						<span
							aria-hidden="true"
							class="block h-full w-full rounded-full"
							:class="selectedColorKey === color.key ? 'border border-white' : 'border border-slate-300'"
							:style="color.swatchStyle"
						></span>
						<span class="sr-only">{{ color.label }}</span>
					</button>
				</li>
			</ul>

			<div class="mt-auto flex items-end gap-3">
				<p class="text-[28px] font-semibold tracking-[-0.05em] text-[#0B8B73] sm:text-[30px] 2xl:text-[32px]">
					{{ formattedPrice }}
				</p>

				<p
					v-if="formattedOriginalPrice"
					class="pb-1 text-[16px] text-slate-400 line-through decoration-2 sm:text-[18px]"
				>
					{{ formattedOriginalPrice }}
				</p>
			</div>
		</div>
	</article>
</template>

<script setup lang="ts">
	import { computed, ref } from 'vue'

	import type { Product } from '../../lib/products'

	const colorSwatches: Record<string, string> = {
		black: '#111827',
		blue: '#5B6EE1',
		gold: '#D4A017',
		gray: '#7B7B7B',
		green: '#1FA34A',
		grey: '#7B7B7B',
		pink: '#F9A8D4',
		purple: '#8B5CF6',
		red: '#D92D48',
		silver: '#D1D5DB',
		white: '#FFFFFF',
	}

	const currencyFormatter = new Intl.NumberFormat('de-DE', {
		style: 'currency',
		currency: 'EUR',
	})

	const props = withDefaults(defineProps<{
		product: Product
		imageLoading?: 'eager' | 'lazy'
		imageFetchPriority?: 'auto' | 'high' | 'low'
	}>(), {
		imageLoading: 'lazy',
		imageFetchPriority: 'auto',
	})

	const emit = defineEmits<{
		'open-detail': [product: Product]
	}>()

	const hasDiscount = computed(() => props.product.discountPercent > 0)
	const normalizeColorName = (color: string) => color.trim().toLowerCase()
	const selectedColorKey = ref<string | null>(props.product.colors[0] ? normalizeColorName(props.product.colors[0]) : null)
	const isCardSelected = ref(false)

	const handleCardClick = () => {
		isCardSelected.value = !isCardSelected.value
		emit('open-detail', props.product)
	}

	const selectSwatch = (colorKey: string) => {
		selectedColorKey.value = colorKey
		isCardSelected.value = true
	}

	const createSwatchStyle = (baseColor: string) => ({
		backgroundColor: baseColor,
		backgroundImage:
			'linear-gradient(135deg, rgba(0, 0, 0, 0.18) 0%, rgba(0, 0, 0, 0.04) 38%, rgba(255, 255, 255, 0.12) 100%)',
	})

	const availableColors = computed(() =>
		props.product.colors.map((color) => ({
			key: normalizeColorName(color),
			label: color,
			swatchStyle: createSwatchStyle(
				colorSwatches[normalizeColorName(color)] ?? '#CBD5E1',
			),
		})),
	)

	const formattedPrice = computed(() => currencyFormatter.format(props.product.price))

	const formattedOriginalPrice = computed(() => {
		if (!hasDiscount.value) {
			return null
		}

		const originalPrice = Number(
			(props.product.price / (1 - props.product.discountPercent / 100)).toFixed(2),
		)

		return currencyFormatter.format(originalPrice)
	})
</script>
