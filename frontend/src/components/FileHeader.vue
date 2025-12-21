<template>
  <header class="header">
    <div class="header-content">
      <h1 class="title">
        <Folder class="icon" />
        OFA Filebrowser<sup>Beta</sup>
      </h1>
      <div class="header-actions">
        <div class="view-toggle">
          <button @click="$emit('change-view', 'list')" :class="{ active: currentView === 'list' }" title="List View">
            <List class="icon" />
          </button>
          <button @click="$emit('change-view', 'grid')" :class="{ active: currentView === 'grid' }" title="Grid View">
            <Grid class="icon" />
          </button>
          <button @click="$emit('change-view', 'tiles')" :class="{ active: currentView === 'tiles' }" title="Tiles View">
            <LayoutGrid class="icon" />
          </button>
        </div>
        <div class="actions">
          <button @click="$emit('upload')" class="btn btn-primary">
            <Upload class="icon" />Upload
          </button>
          <button @click="$emit('create')" class="btn btn-secondary">
            <Plus class="icon" />New
          </button>
        </div>
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
import { Folder, Upload, Plus, Home, ChevronRight, List, Grid, LayoutGrid } from 'lucide-vue-next'
import { computed, ref, watch, onMounted } from 'vue'

const props = defineProps(['currentPath', 'currentView'])
const emit = defineEmits(['upload', 'create', 'navigate', 'change-view'])

const currentPath = computed({
  get: () => props.currentPath,
  set: (val) => emit('navigate', val)
})

watch(currentPath, (val) => {
  localStorage.setItem('currentPath', val)
})

const pathParts = computed(() => {
  // Normalize backslashes to forward slashes first
  const normalized = currentPath.value.replace(/\\/g, '/')
  return normalized.split('/').filter(Boolean)
})

const getPathUpTo = (index) => {
  const parts = pathParts.value.slice(0, index + 1)
  return '/' + parts.join('/')
}

const navigateTo = (path) => {
  currentPath.value = path
}
const goBackOneLevel = () => {
  if (!pathParts.value.length) return
  const parentParts = pathParts.value.slice(0, -1)
  const parentPath = '/' + parentParts.join('/')
  navigateTo(parentPath === '' ? '/' : parentPath)
}

onMounted(() => {
  const savedPath = localStorage.getItem('currentPath')
  if (savedPath) {
    currentPath.value = savedPath
  } else {
    currentPath.value = '/'
  }
})
</script>

<style scoped>
.header { background: var(--bg-secondary); border-radius: 12px; padding: 20px; margin-bottom: 20px; border: 1px solid var(--border-color); }
.header-content { display: flex; justify-content: space-between; align-items: center; margin-bottom: 15px; }
.title { display: flex; align-items: center; gap: 10px; font-size: 24px; font-weight: 600; color: var(--text-primary); }
.header-actions { display: flex; align-items: center; gap: 16px; }
.view-toggle { display: flex; gap: 4px; padding: 4px; background: var(--bg-primary); border-radius: 8px; border: 1px solid var(--border-color); }
.view-toggle button { padding: 8px; background: transparent; border: none; border-radius: 6px; cursor: pointer; color: var(--text-secondary); transition: all 0.2s; display: flex; align-items: center; justify-content: center; }
.view-toggle button:hover { background: var(--bg-hover); color: var(--text-primary); }
.view-toggle button.active { background: var(--accent-color); color: white; }
.actions { display: flex; gap: 10px; }
.btn { display: flex; align-items: center; gap: 8px; padding: 10px 16px; border: none; border-radius: 8px; cursor: pointer; font-size: 14px; font-weight: 500; transition: all 0.2s; }
.btn-primary { background: var(--accent-color); color: white; }
.btn-primary:hover { opacity: 0.9; transform: translateY(-1px); }
.btn-secondary { background: var(--bg-hover); color: var(--text-primary); border: 1px solid var(--border-color); }
.btn-secondary:hover { background: var(--border-color); }
.icon { width: 18px; height: 18px; }
.breadcrumb { display: flex; align-items: center; gap: 8px; }
.breadcrumb-item { padding: 6px 12px; background: var(--bg-primary); border: 1px solid var(--border-color); border-radius: 6px; color: var(--text-primary); display: flex; align-items: center; gap: 4px; cursor: pointer; }
.breadcrumb-item .icon { color: inherit; }
.breadcrumb-item:hover { background: var(--bg-hover); }
.separator { width: 16px; height: 16px; color: var(--text-secondary); }
.back-btn { font-weight: bold; }
</style>