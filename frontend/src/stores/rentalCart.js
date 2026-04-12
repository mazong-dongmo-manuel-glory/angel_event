import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'

const STORAGE_KEY = 'angel-event-rental-cart'

function normalizeItem(item) {
  return {
    id: item.id,
    title: item.title,
    description: item.description || '',
    price: Number(item.price) || 0,
    image_url: item.image_url || '',
    category_id: item.category_id ?? item.category?.id ?? null,
    category_enum: item.category_enum || item.category?.slug || '',
    category: item.category
      ? {
          id: item.category.id,
          name: item.category.name,
          slug: item.category.slug,
        }
      : null,
  }
}

function loadStoredItems() {
  if (typeof window === 'undefined') {
    return []
  }

  try {
    const rawCart = window.localStorage.getItem(STORAGE_KEY)
    const parsedCart = rawCart ? JSON.parse(rawCart) : []

    if (!Array.isArray(parsedCart)) {
      return []
    }

    return parsedCart
      .filter((item) => item && typeof item.id === 'number')
      .map(normalizeItem)
  } catch (error) {
    console.error('Failed to read rental cart from storage:', error)
    return []
  }
}

export const useRentalCartStore = defineStore('rentalCart', () => {
  const cartItems = ref(loadStoredItems())

  watch(
    cartItems,
    (items) => {
      if (typeof window === 'undefined') {
        return
      }

      window.localStorage.setItem(STORAGE_KEY, JSON.stringify(items))
    },
    { deep: true }
  )

  const cartCount = computed(() => cartItems.value.length)
  const cartTotal = computed(() => {
    return cartItems.value.reduce((sum, item) => sum + (Number(item.price) || 0), 0)
  })
  const cartItemIds = computed(() => cartItems.value.map((item) => item.id))

  function hasItem(itemId) {
    return cartItems.value.some((item) => item.id === itemId)
  }

  function addItem(item) {
    const normalizedItem = normalizeItem(item)
    const existingIndex = cartItems.value.findIndex((cartItem) => cartItem.id === normalizedItem.id)

    if (existingIndex >= 0) {
      const updatedItems = [...cartItems.value]
      updatedItems[existingIndex] = normalizedItem
      cartItems.value = updatedItems
      return
    }

    cartItems.value = [...cartItems.value, normalizedItem]
  }

  function removeItem(itemId) {
    cartItems.value = cartItems.value.filter((item) => item.id !== itemId)
  }

  function toggleItem(item) {
    if (hasItem(item.id)) {
      removeItem(item.id)
      return false
    }

    addItem(item)
    return true
  }

  function clearCart() {
    cartItems.value = []
  }

  function syncCatalog(items) {
    const itemsById = new Map(items.map((item) => [item.id, normalizeItem(item)]))
    cartItems.value = cartItems.value.map((item) => itemsById.get(item.id) || item)
  }

  return {
    cartItems,
    cartCount,
    cartTotal,
    cartItemIds,
    hasItem,
    addItem,
    removeItem,
    toggleItem,
    clearCart,
    syncCatalog,
  }
})
