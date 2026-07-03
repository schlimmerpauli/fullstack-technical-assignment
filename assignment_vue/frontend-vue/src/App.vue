<template>
  <div class="min-h-screen bg-[#F9FAFB] text-slate-900">
    <div class="mx-auto flex min-h-screen w-full max-w-none flex-col px-4 py-4 sm:px-6 sm:py-6 lg:px-[44px] lg:py-8">
      <header>
        <div class="flex items-center gap-3">
          <SearchInput v-model="searchTerm" class="min-w-0 flex-1" @search="commitSearchTerm" />

          <button
            ref="mobileFiltersTriggerRef"
            type="button"
            class="inline-flex h-[70px] w-[70px] shrink-0 items-center justify-center rounded-[12px] border border-[#05031A] bg-white text-[#05031A] transition hover:bg-[#05031A] hover:text-white focus:outline-none focus:ring-2 focus:ring-[#05031A] focus:ring-offset-2 xl:hidden"
            :class="isMobileFiltersOpen ? 'bg-[#05031A] text-white' : ''"
            :aria-label="isMobileFiltersOpen ? 'Close filters' : 'Open filters'"
            aria-controls="mobile-filters-panel"
            :aria-expanded="isMobileFiltersOpen"
            aria-haspopup="dialog"
            @click="toggleMobileFilters"
          >
            <svg class="h-7 w-7" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path
                d="M4 7H20"
                stroke="currentColor"
                stroke-width="1.9"
                stroke-linecap="round"
              />
              <path
                d="M7 12H17"
                stroke="currentColor"
                stroke-width="1.9"
                stroke-linecap="round"
              />
              <path
                d="M10 17H14"
                stroke="currentColor"
                stroke-width="1.9"
                stroke-linecap="round"
              />
              <circle cx="16" cy="7" r="2.25" fill="currentColor" />
              <circle cx="9" cy="12" r="2.25" fill="currentColor" />
              <circle cx="12" cy="17" r="2.25" fill="currentColor" />
            </svg>
          </button>
        </div>
      </header>

      <main class="mt-6 grid flex-1 gap-6 lg:mt-8 xl:grid-cols-[290px_minmax(0,1fr)] xl:items-start xl:gap-10">
        <aside class="hidden rounded-[12px] bg-[#FFFFFF] p-5 shadow-[0_20px_40px_rgba(15,23,42,0.06)] xl:block">
          <FiltersPanel
            v-model:selected-categories="selectedCategories"
            v-model:selected-brands="selectedBrands"
            v-model:selected-conditions="selectedConditions"
            v-model:selected-color="selectedColor"
            v-model:min-price="selectedMinPrice"
            v-model:max-price="selectedMaxPrice"
            :category-options="categoryOptions"
            :brand-options="brandOptions"
            :condition-options="conditionOptions"
            :color-options="colorOptions"
            :price-upper-bound="priceUpperBound"
          />
        </aside>

        <section class="min-w-0">
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <p
              class="text-slate-500"
              :class="hasProducts ? 'text-xs font-semibold uppercase tracking-[0.22em]' : 'text-sm font-medium tracking-normal text-slate-600'"
            >
              {{ totalProducts }} products found
            </p>

            <div v-if="hasProducts" class="flex flex-wrap items-center gap-3">
              <label class="relative z-20 flex items-center gap-2 text-sm text-slate-600" for="products-sort">
                <span class="font-medium text-slate-700">Sort by</span>

                <div class="relative z-20">
                  <select
                    id="products-sort"
                    v-model="selectedSort"
                    class="h-10 appearance-none rounded-[12px] border border-slate-300 bg-white pl-3 pr-10 text-left text-sm font-medium text-slate-700 shadow-sm transition focus:outline-none focus:ring-2 focus:ring-slate-300 focus:ring-offset-2"
                  >
                    <option value="">Default order</option>
                    <option value="popularity">Popularity</option>
                  </select>

                  <svg
                    class="pointer-events-none absolute right-3 top-1/2 z-10 h-4 w-4 -translate-y-1/2 text-slate-500"
                    viewBox="0 0 20 20"
                    fill="none"
                    aria-hidden="true"
                  >
                    <path
                      d="M5 7.5L10 12.5L15 7.5"
                      stroke="currentColor"
                      stroke-width="1.8"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                </div>
              </label>

              <p v-if="isLoadingMore" class="text-sm text-slate-500">Loading more...</p>
            </div>
          </div>

          <div
            v-if="errorMessage"
            class="mt-5 rounded-[12px] border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
          >
            <div class="flex items-center justify-between gap-4">
              <p>{{ errorMessage }}</p>
              <button
                type="button"
                class="inline-flex items-center rounded-[10px] border border-red-300 bg-white px-3 py-1.5 font-semibold text-red-700 transition hover:border-red-400 hover:text-red-800 focus:outline-none focus:ring-2 focus:ring-red-200 focus:ring-offset-2"
                @click="reloadProducts"
              >
                Retry
              </button>
            </div>
          </div>

          <div
            v-else-if="isInitialLoadPending"
            class="mt-5 rounded-[12px] border border-slate-200 bg-white px-4 py-8 text-sm text-slate-500 shadow-sm"
          >
            Loading products...
          </div>

          <div
            v-else-if="!hasProducts"
            class="mt-5 flex min-h-[420px] flex-col items-center justify-center px-6 py-12 text-center"
          >
            <div class="flex h-20 w-20 items-center justify-center rounded-full bg-[#F1F2F5] text-[#A5ADBC]">
              <svg class="h-10 w-10" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <path
                  d="M12 3.75L18.75 7.5L12 11.25L5.25 7.5L12 3.75Z"
                  stroke="currentColor"
                  stroke-width="1.7"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M18.75 7.5V15.75L12 19.5V11.25"
                  stroke="currentColor"
                  stroke-width="1.7"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M5.25 7.5V15.75L12 19.5"
                  stroke="currentColor"
                  stroke-width="1.7"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M9 5.4L15 8.7"
                  stroke="currentColor"
                  stroke-width="1.7"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </div>

            <h2 class="mt-7 text-[26px] font-medium tracking-[-0.04em] text-slate-900 sm:text-[30px]">
              No products found
            </h2>

            <p class="mt-3 max-w-[460px] text-[15px] leading-7 text-slate-500 sm:text-[16px]">
              We couldn't find any products matching your current filters. Try adjusting your search criteria.
            </p>

            <button
              type="button"
              class="mt-8 inline-flex h-11 items-center justify-center rounded-[12px] border border-slate-200 bg-white px-5 text-sm font-semibold text-slate-800 shadow-sm transition hover:border-slate-300 hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-slate-200 focus:ring-offset-2"
              @click="clearSearchCriteria"
            >
              Clear all filters
            </button>
          </div>

          <div v-else class="mt-10 grid grid-cols-1 justify-items-center gap-6 md:grid-cols-2 md:justify-items-stretch lg:grid-cols-3 lg:gap-8 xl:gap-10">
            <ItemCard
              v-for="(product, index) in products"
              :key="product.id"
              :product="product"
              :image-loading="index === 0 ? 'eager' : 'lazy'"
              :image-fetch-priority="index === 0 ? 'high' : 'auto'"
              class="mx-auto w-full md:mx-0"
              @open-detail="handleProductOpen"
            />
          </div>

          <div
            v-if="loadMoreErrorMessage"
            class="mt-4 rounded-[12px] border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800"
          >
            {{ loadMoreErrorMessage }}
          </div>

          <div v-if="canLoadMore" class="mt-6 flex justify-center pb-2">
            <button
              type="button"
              class="inline-flex items-center rounded-[12px] border border-slate-300 bg-white px-5 py-2.5 text-sm font-semibold text-slate-700 shadow-sm transition hover:border-slate-400 hover:text-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-300 focus:ring-offset-2"
              :disabled="isInitialLoadPending || isLoadingMore"
              @click="loadMore"
            >
              {{ isLoadingMore ? 'Loading...' : 'Load more' }}
            </button>
          </div>
        </section>
      </main>

      <div
        class="fixed inset-0 z-40 bg-black/30 transition-opacity duration-200 xl:hidden"
        :class="isMobileFiltersOpen ? 'opacity-100' : 'pointer-events-none opacity-0'"
        @click="closeMobileFilters"
      />

      <div
        id="mobile-filters-panel"
        ref="mobileFiltersPanelRef"
        class="fixed inset-y-0 right-0 z-50 flex w-full max-w-[390px] transform flex-col border-l border-black/10 bg-white shadow-[0_28px_64px_rgba(15,23,42,0.22)] transition-transform duration-300 ease-out xl:hidden"
        :class="isMobileFiltersOpen ? 'translate-x-0' : 'translate-x-full'"
        :aria-hidden="!isMobileFiltersOpen"
        :inert="!isMobileFiltersOpen || undefined"
        role="dialog"
        aria-modal="true"
        aria-labelledby="mobile-filters-title"
        @keydown.esc.window="closeMobileFilters"
      >
        <div class="flex items-start justify-between gap-4 border-b border-slate-200 px-5 py-4">
          <div>
            <p class="text-xs font-semibold uppercase tracking-[0.22em] text-slate-500">Filters</p>
            <h2 id="mobile-filters-title" class="mt-1 text-[20px] font-medium tracking-[-0.03em] text-slate-950">Refine your results</h2>
          </div>

          <button
            ref="mobileFiltersCloseButtonRef"
            type="button"
            class="inline-flex h-11 w-11 items-center justify-center rounded-[12px] border border-slate-200 text-slate-700 transition hover:border-slate-300 hover:text-slate-950 focus:outline-none focus:ring-2 focus:ring-slate-300 focus:ring-offset-2"
            aria-label="Close filters"
            @click="closeMobileFilters"
          >
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="none" aria-hidden="true">
              <path d="M5 5L15 15" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
              <path d="M15 5L5 15" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
          </button>
        </div>

        <div class="flex-1 overflow-y-auto px-5 py-5">
          <FiltersPanel
            v-model:selected-categories="selectedCategories"
            v-model:selected-brands="selectedBrands"
            v-model:selected-conditions="selectedConditions"
            v-model:selected-color="selectedColor"
            v-model:min-price="selectedMinPrice"
            v-model:max-price="selectedMaxPrice"
            :category-options="categoryOptions"
            :brand-options="brandOptions"
            :condition-options="conditionOptions"
            :color-options="colorOptions"
            :price-upper-bound="priceUpperBound"
          />
        </div>

        <div class="border-t border-slate-200 px-5 py-4">
          <div class="flex items-center gap-3">
            <button
              type="button"
              class="inline-flex h-11 items-center justify-center rounded-[12px] border border-[#05031A] px-4 text-sm font-semibold text-[#05031A] transition hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-[#05031A] focus:ring-offset-2"
              @click="clearFilters"
            >
              Clear all
            </button>

            <button
              type="button"
              class="inline-flex h-11 flex-1 items-center justify-center rounded-[12px] bg-[#05031A] px-4 text-sm font-semibold text-white transition hover:bg-slate-900 focus:outline-none focus:ring-2 focus:ring-[#05031A] focus:ring-offset-2"
              @click="closeMobileFilters"
            >
              View results
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import FiltersPanel from './components/filters/filters.vue'
import ItemCard from './components/item-card/item-card.vue'
import SearchInput from './components/search-input/search-input.vue'
import { listProducts } from './lib/api'
import { filterClassicColors } from './lib/colors'
import type { Pagination, Product, ProductListQuery, ProductSort } from './lib/products'

const LOAD_MORE_BATCH_SIZE = 6
const FILTER_METADATA_PAGE_SIZE = 100
const DEFAULT_MIN_PRICE = 0
const INITIAL_MAX_PRICE = 2000

type FilterOption = {
  label: string
  value: string
}

const normalizeText = (value: string) => value.trim().toLowerCase()

const capitalizeLabel = (value: string) => value.charAt(0).toUpperCase() + value.slice(1)

const formatFilterLabel = (value: string) => value
  .split(/[\s_-]+/)
  .filter(Boolean)
  .map((segment) => capitalizeLabel(segment.toLowerCase()))
  .join(' ')

const buildFilterOptions = (values: string[]) => {
  const uniqueValues = new Map<string, string>()

  for (const value of values) {
    const normalizedValue = normalizeText(value)

    if (normalizedValue === '' || uniqueValues.has(normalizedValue)) {
      continue
    }

    uniqueValues.set(normalizedValue, value)
  }

  return Array.from(uniqueValues.entries())
    .sort((firstEntry, secondEntry) => firstEntry[1].localeCompare(secondEntry[1]))
    .map(([filterValue, label]) => ({
      value: filterValue,
      label: formatFilterLabel(label),
    }))
}

const deriveMaxPrice = (products: Product[]) => {
  if (products.length === 0) {
    return INITIAL_MAX_PRICE
  }

  return Math.ceil(Math.max(...products.map((product) => product.price)) / 50) * 50
}

const buildListQuery = (baseQuery: ProductListQuery, page: number): ProductListQuery => ({
  ...baseQuery,
  page,
  pageSize: LOAD_MORE_BATCH_SIZE,
})

const products = ref<Product[]>([])
const catalogProducts = ref<Product[]>([])
const pagination = ref<Pagination | null>(null)
const currentPage = ref(1)
const searchTerm = ref('')
const debouncedSearchTerm = ref('')
const selectedSort = ref<ProductSort | ''>('')
const selectedCategories = ref<string[]>([])
const selectedBrands = ref<string[]>([])
const selectedConditions = ref<string[]>([])
const selectedColor = ref<string | null>(null)
const selectedMinPrice = ref(DEFAULT_MIN_PRICE)
const selectedMaxPrice = ref(INITIAL_MAX_PRICE)
const priceUpperBound = ref(INITIAL_MAX_PRICE)
const isInitialLoadPending = ref(false)
const isLoadingMore = ref(false)
const errorMessage = ref('')
const loadMoreErrorMessage = ref('')
const isMobileFiltersOpen = ref(false)
const mobileFiltersTriggerRef = ref<HTMLButtonElement | null>(null)
const mobileFiltersPanelRef = ref<HTMLElement | null>(null)
const mobileFiltersCloseButtonRef = ref<HTMLButtonElement | null>(null)

let productRequestController: AbortController | null = null
let catalogRequestController: AbortController | null = null
let hasRequestedFilterMetadata = false

const categoryOptions = computed<FilterOption[]>(() => buildFilterOptions(catalogProducts.value.map((product) => product.category)))
const brandOptions = computed<FilterOption[]>(() => buildFilterOptions(catalogProducts.value.map((product) => product.brand)))
const conditionOptions = computed<FilterOption[]>(() => buildFilterOptions(catalogProducts.value.map((product) => product.condition)))
const colorOptions = computed<FilterOption[]>(() => buildFilterOptions(
  catalogProducts.value.flatMap((product) => filterClassicColors(product.colors)),
))

const hasProducts = computed(() => products.value.length > 0)
const totalProducts = computed(() => pagination.value?.total ?? 0)
const canLoadMore = computed(() => pagination.value?.hasMore ?? false)

const baseQuery = computed<ProductListQuery>(() => {
  const query: ProductListQuery = {}
  const normalizedSearch = debouncedSearchTerm.value.trim()

  if (normalizedSearch !== '') {
    query.search = normalizedSearch
  }

  if (selectedCategories.value.length > 0) {
    query.categories = [...selectedCategories.value]
  }

  if (selectedBrands.value.length > 0) {
    query.brands = [...selectedBrands.value]
  }

  if (selectedConditions.value.length > 0) {
    query.conditions = [...selectedConditions.value]
  }

  if (selectedColor.value !== null) {
    query.color = selectedColor.value
  }

  if (selectedSort.value !== '') {
    query.sort = selectedSort.value
  }

  if (selectedMinPrice.value > DEFAULT_MIN_PRICE) {
    query.minPrice = selectedMinPrice.value
  }

  if (selectedMaxPrice.value < priceUpperBound.value) {
    query.maxPrice = selectedMaxPrice.value
  }

  return query
})

const fetchProducts = async (page: number, append: boolean) => {
  productRequestController?.abort()

  const requestController = new AbortController()
  productRequestController = requestController

  if (append) {
    isLoadingMore.value = true
    loadMoreErrorMessage.value = ''
  } else {
    isInitialLoadPending.value = true
    errorMessage.value = ''
    loadMoreErrorMessage.value = ''
    products.value = []
    pagination.value = null
  }

  try {
    const response = await listProducts(buildListQuery(baseQuery.value, page), {
      signal: requestController.signal,
    })

    if (requestController.signal.aborted) {
      return
    }

    products.value = append ? [...products.value, ...response.products] : response.products
    pagination.value = response.pagination
    currentPage.value = response.pagination.page
  } catch (error) {
    if (requestController.signal.aborted) {
      return
    }

    const message = error instanceof Error ? error.message : 'Failed to load products.'

    if (append) {
      loadMoreErrorMessage.value = message
    } else {
      errorMessage.value = message
      currentPage.value = 1
    }
  } finally {
    if (productRequestController === requestController) {
      productRequestController = null
    }

    if (!append && !requestController.signal.aborted) {
      ensureFilterMetadataLoaded()
    }

    if (append) {
      isLoadingMore.value = false
    } else {
      isInitialLoadPending.value = false
    }
  }
}

const loadFilterMetadata = async () => {
  catalogRequestController?.abort()

  const requestController = new AbortController()
  catalogRequestController = requestController

  try {
    const response = await listProducts(
      {
        page: 1,
        pageSize: FILTER_METADATA_PAGE_SIZE,
      },
      {
        signal: requestController.signal,
      },
    )

    if (requestController.signal.aborted) {
      return
    }

    const previousUpperBound = priceUpperBound.value
    const derivedMaxPrice = deriveMaxPrice(response.products)

    catalogProducts.value = response.products
    priceUpperBound.value = derivedMaxPrice

    if (selectedMaxPrice.value === previousUpperBound || selectedMaxPrice.value > derivedMaxPrice) {
      selectedMaxPrice.value = derivedMaxPrice
    }
  } catch {
    // The listing can still function without derived filter metadata.
  } finally {
    if (catalogRequestController === requestController) {
      catalogRequestController = null
    }
  }
}

const ensureFilterMetadataLoaded = () => {
  if (hasRequestedFilterMetadata) {
    return
  }

  hasRequestedFilterMetadata = true
  void loadFilterMetadata()
}

const loadMore = () => {
  if (isInitialLoadPending.value || isLoadingMore.value || !canLoadMore.value) {
    return
  }

  void fetchProducts(currentPage.value + 1, true)
}

const reloadProducts = () => {
  currentPage.value = 1
  void fetchProducts(1, false)
}

const handleProductOpen = (product: Product) => {
  console.log('Open detail - out of scope', product.id)
}

const commitSearchTerm = (value: string) => {
  debouncedSearchTerm.value = value
}

const clearFilters = () => {
  selectedCategories.value = []
  selectedBrands.value = []
  selectedConditions.value = []
  selectedColor.value = null
  selectedMinPrice.value = DEFAULT_MIN_PRICE
  selectedMaxPrice.value = priceUpperBound.value
}

const clearSearchCriteria = () => {
  searchTerm.value = ''
  debouncedSearchTerm.value = ''
  selectedSort.value = ''
  clearFilters()
}

const closeMobileFilters = () => {
  isMobileFiltersOpen.value = false
}

const toggleMobileFilters = () => {
  isMobileFiltersOpen.value = !isMobileFiltersOpen.value
}

const syncBodyScrollLock = (isLocked: boolean) => {
  if (typeof document === 'undefined') {
    return
  }

  document.body.style.overflow = isLocked ? 'hidden' : ''
}

const moveFocusOutOfMobileFilters = () => {
  if (typeof document === 'undefined') {
    return
  }

  const activeElement = document.activeElement

  if (!(activeElement instanceof HTMLElement) || !mobileFiltersPanelRef.value?.contains(activeElement)) {
    return
  }

  if (typeof window !== 'undefined' && window.innerWidth < 1280) {
    mobileFiltersTriggerRef.value?.focus()
  }

  if (mobileFiltersPanelRef.value?.contains(document.activeElement)) {
    activeElement.blur()
  }
}

const handleViewportResize = () => {
  if (window.innerWidth >= 1280) {
    closeMobileFilters()
  }
}

watch(isMobileFiltersOpen, async (isOpen) => {
  syncBodyScrollLock(isOpen)

  if (isOpen) {
    await nextTick()
    mobileFiltersCloseButtonRef.value?.focus()
    return
  }

  moveFocusOutOfMobileFilters()
})

watch(baseQuery, () => {
  currentPage.value = 1
  void fetchProducts(1, false)
}, { immediate: true })

onMounted(() => {
  window.addEventListener('resize', handleViewportResize)
})

onBeforeUnmount(() => {
  productRequestController?.abort()
  catalogRequestController?.abort()
  syncBodyScrollLock(false)
  window.removeEventListener('resize', handleViewportResize)
})
</script>
