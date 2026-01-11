<template>
  <div class="form-input" :class="{ 'has-error': error, 'has-icon': icon }">
    <label v-if="label" :for="inputId" class="form-label">
      {{ label }}
      <span v-if="required" class="required">*</span>
    </label>
    
    <div class="input-wrapper">
      <component v-if="icon" :is="icon" class="input-icon" />
      
      <input
        v-if="type !== 'textarea' && type !== 'select'"
        :id="inputId"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :required="required"
        :disabled="disabled"
        @input="$emit('update:modelValue', $event.target.value)"
        @blur="$emit('blur')"
        class="form-control"
      />
      
      <textarea
        v-else-if="type === 'textarea'"
        :id="inputId"
        :value="modelValue"
        :placeholder="placeholder"
        :required="required"
        :disabled="disabled"
        :rows="rows"
        @input="$emit('update:modelValue', $event.target.value)"
        @blur="$emit('blur')"
        class="form-control"
      ></textarea>
      
      <select
        v-else-if="type === 'select'"
        :id="inputId"
        :value="modelValue"
        :required="required"
        :disabled="disabled"
        @change="$emit('update:modelValue', $event.target.value)"
        @blur="$emit('blur')"
        class="form-control"
      >
        <option value="" disabled>{{ placeholder }}</option>
        <slot name="options"></slot>
      </select>
    </div>
    
    <span v-if="error" class="error-message">{{ error }}</span>
    <span v-else-if="hint" class="hint-message">{{ hint }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  type: {
    type: String,
    default: 'text',
    validator: (value) => ['text', 'email', 'password', 'tel', 'number', 'date', 'textarea', 'select'].includes(value)
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  required: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  },
  hint: {
    type: String,
    default: ''
  },
  icon: {
    type: Object,
    default: null
  },
  rows: {
    type: Number,
    default: 4
  }
})

defineEmits(['update:modelValue', 'blur'])

const inputId = computed(() => {
  return `input-${Math.random().toString(36).substr(2, 9)}`
})
</script>

<style scoped>
.form-input {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-lg);
}

.form-label {
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
  font-size: var(--font-size-sm);
  letter-spacing: 0.02em;
}

.required {
  color: var(--color-gold);
  margin-left: 2px;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: var(--spacing-md);
  color: var(--color-gray);
  width: 20px;
  height: 20px;
  pointer-events: none;
}

.has-icon .form-control {
  padding-left: calc(var(--spacing-md) + 20px + var(--spacing-md));
}

.form-control {
  width: 100%;
  padding: var(--spacing-md) var(--spacing-lg);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-sans);
  font-size: var(--font-size-base);
  color: var(--color-text);
  background: white;
  transition: all var(--transition-base);
}

.form-control:focus {
  outline: none;
  border-color: var(--color-gold);
  box-shadow: 0 0 0 3px rgba(212, 175, 55, 0.1);
}

.form-control:disabled {
  background: var(--color-gray-lighter);
  cursor: not-allowed;
  opacity: 0.6;
}

.form-control::placeholder {
  color: var(--color-gray);
}

textarea.form-control {
  resize: vertical;
  min-height: 100px;
}

select.form-control {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%23666' d='M6 9L1 4h10z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right var(--spacing-md) center;
  padding-right: calc(var(--spacing-lg) + 20px);
}

.has-error .form-control {
  border-color: var(--color-error);
}

.has-error .form-control:focus {
  box-shadow: 0 0 0 3px rgba(193, 41, 46, 0.1);
}

.error-message {
  color: var(--color-error);
  font-size: var(--font-size-sm);
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.error-message::before {
  content: '⚠';
}

.hint-message {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
}

/* Animations */
.form-control {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
