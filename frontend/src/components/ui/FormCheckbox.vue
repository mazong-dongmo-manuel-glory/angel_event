<template>
  <div class="form-checkbox" :class="{ 'has-error': error }">
    <label class="checkbox-label">
      <input
        type="checkbox"
        :checked="modelValue"
        :disabled="disabled"
        @change="$emit('update:modelValue', $event.target.checked)"
        class="checkbox-input"
      />
      <span class="checkbox-custom"></span>
      <span class="checkbox-text">
        {{ label }}
        <span v-if="required" class="required">*</span>
      </span>
    </label>
    
    <span v-if="error" class="error-message">{{ error }}</span>
    <span v-else-if="hint" class="hint-message">{{ hint }}</span>
  </div>
</template>

<script setup>
defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  label: {
    type: String,
    required: true
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
  }
})

defineEmits(['update:modelValue'])
</script>

<style scoped>
.form-checkbox {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-lg);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  cursor: pointer;
  user-select: none;
}

.checkbox-input {
  position: absolute;
  opacity: 0;
  cursor: pointer;
}

.checkbox-custom {
  position: relative;
  width: 20px;
  height: 20px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: white;
  transition: all var(--transition-base);
  flex-shrink: 0;
}

.checkbox-input:checked + .checkbox-custom {
  background: var(--color-gold);
  border-color: var(--color-gold);
}

.checkbox-input:checked + .checkbox-custom::after {
  content: '';
  position: absolute;
  left: 6px;
  top: 2px;
  width: 5px;
  height: 10px;
  border: solid white;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.checkbox-input:focus + .checkbox-custom {
  box-shadow: 0 0 0 3px rgba(212, 175, 55, 0.1);
}

.checkbox-input:disabled + .checkbox-custom {
  background: var(--color-gray-lighter);
  cursor: not-allowed;
  opacity: 0.6;
}

.checkbox-text {
  color: var(--color-text);
  font-size: var(--font-size-base);
}

.required {
  color: var(--color-gold);
  margin-left: 2px;
}

.has-error .checkbox-custom {
  border-color: var(--color-error);
}

.error-message {
  color: var(--color-error);
  font-size: var(--font-size-sm);
  margin-left: calc(20px + var(--spacing-md));
}

.hint-message {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
  margin-left: calc(20px + var(--spacing-md));
}
</style>
