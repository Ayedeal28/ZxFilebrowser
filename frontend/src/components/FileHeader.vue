<template>
  <header class="header">
    <div class="header-content">
      <h1 class="title">
        <Folder class="icon" />
        OFA Filebrowser<sup>Beta</sup>
      </h1>
      <div class="actions">
        <button @click="$emit('upload')" class="btn btn-primary">
          <Upload class="icon" />Upload
        </button>
        <button @click="$emit('create')" class="btn btn-secondary">
          <Plus class="icon" />New
        </button>
      </div>
    </div>

    <div class="breadcrumb">
      <button v-if="pathParts.length" @click="goBackOneLevel" class="breadcrumb-item back-btn">
        ←
      </button>
      <button @click="navigateTo('/')" class="breadcrumb-item">
        <Home class="icon" />
      </button>
      <ChevronRight v-if="pathParts.length" class="separator" />
      <template v-for="(part, index) in pathParts" :key="index">
        <button @click="navigateTo(getPathUpTo(index))" class="breadcrumb-item">
          {{ part }}
        </button>
        <ChevronRight v-if="index < pathParts.length - 1" class="separator" />
      </template>
    </div>
  </header>
</template>

<script setup>
import { Folder, Upload, Plus, Home, ChevronRight } from 'lucide-vue-next'
import { computed, ref, watch, onMounted } from 'vue'

const props = defineProps(['currentPath'])
const emit = defineEmits(['upload', 'create', 'navigate'])

// Use a reactive currentPath that persists in localStorage
const storedPath = localStorage.getItem('currentPath') || props.currentPath || '/'
const currentPath = ref(storedPath)

// Save currentPath to localStorage whenever it changes
watch(currentPath, (val) => {
  localStorage.setItem('currentPath', val)
})

// Compute path parts for breadcrumb
const pathParts = computed(() => {
  // Normalize backslashes to forward slashes first
  const normalized = currentPath.value.replace(/\\/g, '/')
  return normalized.split('/').filter(Boolean)
})

const getPathUpTo = (index) => {
  const parts = pathParts.value.slice(0, index + 1)
  return '/' + parts.join('/')
}

// Navigate and update currentPath + emit event
const navigateTo = (path) => {
  currentPath.value = path
  emit('navigate', path)
}

// Go back one level
const goBackOneLevel = () => {
  if (!pathParts.value.length) return
  const parentParts = pathParts.value.slice(0, -1)
  const parentPath = '/' + parentParts.join('/')
  // Emit navigate event
  navigateTo(parentPath === '' ? '/' : parentPath)
}

// Optional: onMounted, restore from localStorage if empty
onMounted(() => {
  if (!currentPath.value) currentPath.value = '/'
})
</script>

<style scoped>
.header { background: var(--bg-secondary); border-radius: 12px; padding: 20px; margin-bottom: 20px; border: 1px solid var(--border-color); }
.header-content { display: flex; justify-content: space-between; align-items: center; margin-bottom: 15px; }
.title { display: flex; align-items: center; gap: 10px; font-size: 24px; font-weight: 600; color: var(--text-primary); }
.actions { display: flex; gap: 10px; }
.btn { display: flex; align-items: center; gap: 8px; padding: 10px 16px; border: none; border-radius: 8px; cursor: pointer; font-size: 14px; font-weight: 500; transition: all 0.2s; }
.btn-primary { background: var(--accent-color); color: white; }
.btn-primary:hover { opacity: 0.9; transform: translateY(-1px); }
.btn-secondary { background: var(--bg-hover); color: var(--text-primary); border: 1px solid var(--border-color); }
.btn-secondary:hover { background: var(--border-color); }
.icon { width: 18px; height: 18px; }
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px; /* ← Make sure this exists */
}

.breadcrumb-item {
  padding: 6px 12px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-primary); /* text color */
  display: flex;
  align-items: center;
  gap: 4px;
}
.breadcrumb-item .icon {
  color: inherit;
}
body.dark-mode {
  --bg-primary: #2b2b2b;       /* darker bg */
  --text-primary: #ffffff;     /* white text */
  --border-color: #444;        /* border in dark */
}
.breadcrumb-item:hover { background: var(--bg-hover); }
.separator { width: 16px; height: 16px; color: var(--text-secondary); }
.back-btn { font-weight: bold; }
</style>
