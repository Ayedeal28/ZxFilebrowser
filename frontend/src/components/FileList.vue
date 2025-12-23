<template>
  <div v-if="loading" class="loading">
    <Loader2 class="spinner" />Loading...
  </div>

  <div v-else class="file-container" :class="`view-${viewMode}`">
    <div v-if="files.length === 0" class="empty-state">
      <FileQuestion class="empty-icon" />
      <p>This folder is empty</p>
    </div>

    <!-- List View -->
    <div v-if="viewMode === 'list'" class="file-list">
      <div v-for="file in files" :key="file.path" class="file-item-list" @click="$emit('file-click', file)" @contextmenu.prevent="$emit('context-menu', $event, file)">
        <div class="file-icon">
          <Folder v-if="file.isDir" class="icon folder" />
          <FileIcon v-else class="icon file" />
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

    <!-- Grid View -->
    <div v-if="viewMode === 'grid'" class="file-grid">
      <div v-for="file in files" :key="file.path" class="file-item-grid" @click="$emit('file-click', file)" @contextmenu.prevent="$emit('context-menu', $event, file)">
        <div class="file-preview">
          <img 
            v-if="isImage(file)" 
            :src="getThumbnailUrl(file)" 
            :key="file.path"
            alt="" 
            class="thumbnail" 
            loading="lazy"
            @error="handleImageError"
          >
          <Folder v-else-if="file.isDir" class="preview-icon folder" />
          <FileIcon v-else class="preview-icon file" />
        </div>
        <div class="file-name">{{ file.name }}</div>
        <div class="file-meta" v-if="!file.isDir">{{ formatSize(file.size) }}</div>
      </div>
    </div>

    <!-- Tiles View -->
    <div v-if="viewMode === 'tiles'" class="file-tiles">
      <div v-for="file in files" :key="file.path" class="file-item-tile" @click="$emit('file-click', file)" @contextmenu.prevent="$emit('context-menu', $event, file)">
        <div class="tile-preview">
          <img 
            v-if="isImage(file)" 
            :src="getThumbnailUrl(file)" 
            :key="file.path"
            alt="" 
            class="thumbnail" 
            loading="lazy"
            @error="handleImageError"
          >
          <Folder v-else-if="file.isDir" class="preview-icon folder" />
          <FileIcon v-else class="preview-icon file" />
          <div class="tile-overlay">
            <div class="tile-name">{{ file.name }}</div>
            <div class="tile-meta" v-if="!file.isDir">{{ formatSize(file.size) }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Folder, File as FileIcon, Loader2, FileQuestion } from 'lucide-vue-next'

const props = defineProps(['files', 'loading', 'viewMode', 'sourceId'])
defineEmits(['file-click', 'context-menu'])

const isImage = (file) => {
  if (file.isDir) return false
  const imageExts = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp', '.svg']
  const ext = file.ext?.toLowerCase() || ''
  return imageExts.includes(ext)
}

const getThumbnailUrl = (file) => {
  // Use serve endpoint for now, can optimize with thumbnail endpoint later
  const baseUrl = import.meta.env.VITE_API_URL || '/api'
  return `${baseUrl}/serve?source=${props.sourceId}&path=${encodeURIComponent(file.path)}`
}

const handleImageError = (e) => {
  e.target.style.display = 'none'
}

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

.file-container { background: var(--bg-secondary); border-radius: 12px; overflow: hidden; border: 1px solid var(--border-color); padding: 16px; }
.empty-state { display: flex; flex-direction: column; align-items: center; padding: 60px 20px; color: var(--text-secondary); }
.empty-icon { width: 64px; height: 64px; margin-bottom: 10px; }

/* List View */
.file-list {
  display: flex;
  flex-direction: column;
}

.file-item-list {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;
  border-bottom: 1px solid var(--border-color);
}

.file-item-list:last-child {
  border-bottom: none;
}

.file-item-list:hover {
  background: var(--bg-hover);
}

.file-icon {
  flex-shrink: 0;
}

.file-icon .folder {
  color: var(--accent-color);
  width: 24px;
  height: 24px;
}

.file-icon .file {
  color: var(--text-secondary);
  width: 24px;
  height: 24px;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-weight: 500;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-meta {
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 2px;
}

.file-meta span:not(:last-child)::after {
  content: '•';
  margin: 0 8px;
}


/* Grid View */
.file-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 16px; }
.file-item-grid { display: flex; flex-direction: column; align-items: center; padding: 12px; border-radius: 12px; cursor: pointer; transition: all 0.2s; border: 1px solid transparent; }
.file-item-grid:hover { background: var(--bg-hover); border-color: var(--border-color); transform: translateY(-2px); }
.file-preview { width: 100px; height: 100px; display: flex; align-items: center; justify-content: center; margin-bottom: 8px; border-radius: 8px; overflow: hidden; background: var(--bg-primary); }
.thumbnail { width: 100%; height: 100%; object-fit: cover; }
.preview-icon { width: 48px; height: 48px; }
.preview-icon.folder { color: var(--accent-color); }
.preview-icon.file { color: var(--text-secondary); }
.file-item-grid .file-name { font-size: 13px; font-weight: 500; color: var(--text-primary); text-align: center; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; width: 100%; }
.file-item-grid .file-meta { font-size: 11px; color: var(--text-secondary); margin-top: 2px; }

/* Tiles View */
.file-tiles { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 16px; }
.file-item-tile { border-radius: 12px; overflow: hidden; cursor: pointer; transition: all 0.2s; border: 1px solid var(--border-color); }
.file-item-tile:hover { transform: translateY(-4px); box-shadow: 0 8px 24px rgba(0,0,0,0.12); }
.tile-preview { position: relative; width: 100%; padding-bottom: 75%; background: var(--bg-primary); overflow: hidden; }
.tile-preview .thumbnail { position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: cover; }
.tile-preview .preview-icon { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); width: 64px; height: 64px; }
.tile-preview .preview-icon.folder { color: var(--accent-color); }
.tile-preview .preview-icon.file { color: var(--text-secondary); }
.tile-overlay { position: absolute; bottom: 0; left: 0; right: 0; background: linear-gradient(to top, rgba(0,0,0,0.8), transparent); padding: 16px 12px 12px; color: white; }
.tile-name { font-size: 14px; font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.tile-meta { font-size: 11px; opacity: 0.9; margin-top: 4px; }
</style>