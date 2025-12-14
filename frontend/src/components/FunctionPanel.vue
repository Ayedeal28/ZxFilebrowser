<template>
  <div :class="['FunctionPanel', { pinned }]">
    <button @click="toggleDarkMode">🌓 {{ isDark ? 'Dark On' : 'Dark Off' }}</button>
    <button @click="toggleEncryption">{{ encryptionOn ? '🔒 Encryption On' : '🔓 Encryption Off' }}</button>
    <button @click="pinned = !pinned">{{ pinned ? '📌 Pinned' : '📍 Unpinned' }}</button>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const pinned = ref(false)
const encryptionOn = ref(false)
const isDark = ref(false)

const toggleDarkMode = () => {
  isDark.value = !isDark.value
  document.body.classList.toggle('dark', isDark.value)
}

const toggleEncryption = () => {
  encryptionOn.value = !encryptionOn.value
}

// Optional: persist dark mode in localStorage
watch(isDark, val => localStorage.setItem('darkMode', val))
if (localStorage.getItem('darkMode') === 'true') toggleDarkMode()
</script>

<style scoped>
.FunctionPanel {
  position: fixed;
  top: 20px;
  right: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: var(--color-surface);
  padding: 12px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  transition: transform 0.3s ease, background 0.3s ease;
}

.FunctionPanel:not(.pinned) {
  transform: translateX(100%);
}

.FunctionPanel button {
  padding: 8px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  background: var(--color-background);
  transition: background 0.2s;
}

.FunctionPanel button:hover {
  background: var(--color-primary);
  color: white;
}
</style>
