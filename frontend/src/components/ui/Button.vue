<template>
  <button 
    :class="['btn', `btn-${variant}`, `btn-${size}`, { 'btn-loading': loading, 'w-full': block }]"
    :disabled="disabled || loading"
    v-bind="$attrs"
  >
    <div v-if="loading" class="spinner-border"></div>
    <span v-else class="btn-content">
      <slot name="icon-left"></slot>
      <slot></slot>
      <slot name="icon-right"></slot>
    </span>
  </button>
</template>

<script setup>
defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'outline', 'ghost', 'white'].includes(value)
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  },
  loading: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  block: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid transparent;
  border-radius: var(--radius-sm); /* Sharper, more luxury */
  cursor: pointer;
  font-family: var(--font-sans);
  font-weight: 500;
  letter-spacing: 0.1em; /* Luxury spacing */
  text-transform: uppercase;
  transition: all var(--transition-base);
  position: relative;
  overflow: hidden;
  font-size: var(--font-size-sm);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Sizes */
.btn-sm {
  padding: 0.5rem 1rem;
  font-size: 0.75rem;
  height: 36px;
}

.btn-md {
  padding: 0.75rem 2rem;
  height: 48px;
}

.btn-lg {
  padding: 1rem 3rem;
  font-size: 1rem;
  height: 60px;
}

/* Variants */
.btn-primary {
  background-color: var(--color-black);
  color: var(--color-white);
  border-color: var(--color-black);
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--color-black-soft);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.btn-secondary {
  background-color: var(--color-gold);
  color: var(--color-white);
}

.btn-secondary:hover:not(:disabled) {
  background-color: #C5A028;
}

.btn-outline {
  background-color: transparent;
  border-color: var(--color-black);
  color: var(--color-black);
}

.btn-outline:hover:not(:disabled) {
  background-color: var(--color-black);
  color: var(--color-white);
}

.btn-white {
  background-color: var(--color-white);
  color: var(--color-black);
  border-color: var(--color-white);
}

.btn-white:hover:not(:disabled) {
  background-color: var(--color-gray-light);
}

.btn-ghost {
  background-color: transparent;
  color: var(--color-black);
}

.btn-ghost:hover:not(:disabled) {
  background-color: rgba(0, 0, 0, 0.05);
}

.btn-content {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.w-full {
  width: 100%;
}

/* Loading Spinner */
.spinner-border {
  width: 1.2em;
  height: 1.2em;
  border: 2px solid currentColor;
  border-right-color: transparent;
  border-radius: 50%;
  animation: spin 0.75s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
