<script setup>
import './assets/styles/variables.css'
import './assets/styles/animations.css'
</script>

<template>
  <RouterView />
</template>

<style>
/* Smooth scrolling */
html {
  scroll-behavior: smooth;
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 10px;
}

::-webkit-scrollbar-track {
  background: var(--color-gray-lighter);
}

::-webkit-scrollbar-thumb {
  background: var(--color-gold);
  border-radius: var(--radius-full);
}

::-webkit-scrollbar-thumb:hover {
  background: var(--color-gold-dark);
}
</style>
