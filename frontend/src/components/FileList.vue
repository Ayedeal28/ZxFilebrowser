<template>
  <div v-if="loading" class="loading">
    <Loader2 class="spinner" />Loading...
  </div>

  <div v-else class="file-list">
    <div v-if="files.length === 0" class="empty-state">
      <FileQuestion class="empty-icon" />
      <p>This folder is empty</p>
    </div>

    <div v-for="file in files" :key="file.path" class="file-item" @click="$emit('file-click', file)" @contextmenu.prevent="$emit('context-menu', $event, file)">
      <div class="file-icon">
        <Folder v-if="file.isDir" class="icon folder" />
        <File v-else class="icon file" />
      </div>
      <div class="file-info">
        <div class="file-name">{{ file.name }}</div>
        <div class="file-meta">
          <span v-if="!file.isDir">{{ formatSize(file.size) }}</span>
          <span>{{ formatDate(file.modTime) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Folder, File, Loader2, FileQuestion } from 'lucide-vue-next'

defineProps(['files', 'loading'])
defineEmits(['file-click', 'context-menu'])

const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const formatDate = (date) => new Date(date).toLocaleString()
</script>

<style scoped>
.loading { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 60px 20px; color: var(--text-secondary); }
.spinner { width: 40px; height: 40px; margin-bottom: 10px; animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
.file-list { background: var(--bg-secondary); border-radius: 12px; overflow: hidden; border: 1px solid var(--border-color); }
.empty-state { display: flex; flex-direction: column; align-items: center; padding: 60px 20px; color: var(--text-secondary); }
.empty-icon { width: 64px; height: 64px; margin-bottom: 10px; }
.file-item { display: flex; align-items: center; gap: 12px; padding: 12px 16px; border-bottom: 1px solid var(--border-color); cursor: pointer; transition: background 0.2s; }
.file-item:hover { background: var(--bg-hover); }
.file-item:last-child { border-bottom: none; }
.file-icon { flex-shrink: 0; }
.file-icon .folder { color: var(--accent-color); width: 24px; height: 24px; }
.file-icon .file { color: var(--text-secondary); width: 24px; height: 24px; }
.file-info { flex: 1; min-width: 0; }
.file-name { font-weight: 500; color: var(--text-primary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.file-meta { font-size: 12px; color: var(--text-secondary); margin-top: 2px; }
.file-meta span:not(:last-child)::after { content: '•'; margin: 0 8px; }
</style>