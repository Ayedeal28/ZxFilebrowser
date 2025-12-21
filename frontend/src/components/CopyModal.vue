<template>
  <div class="modal" @click.self="$emit('close')">
    <div class="modal-content">
      <h2>{{ mode === 'move' ? 'Move' : 'Copy' }} "{{ fileName }}" to...</h2>
      
      <!-- Source Selector -->
      <div class="form-group">
        <label>Destination Source:</label>
        <select v-model="selectedSource">
          <option v-for="src in sources" :key="src.id" :value="src.id">
            {{ src.name }}
          </option>
        </select>
      </div>

      <!-- Path Browser -->
      <div class="path-browser">
        <div class="path-header">
          <button @click="navigateUp" :disabled="currentPath === '/'" class="nav-btn">
            <ArrowLeft class="icon" />
          </button>
          <span class="current-path">{{ currentPath }}</span>
        </div>
        
        <div class="folder-list">
          <div 
            v-for="folder in folders" 
            :key="folder.path" 
            @click="navigateInto(folder)"
            @dblclick="selectAndCopy(folder.path)"
            class="folder-item"
            :class="{ selected: selectedPath === folder.path }"
          >
            <Folder class="icon" />
            <span>{{ folder.name }}</span>
          </div>
          <div v-if="folders.length === 0" class="empty">No folders</div>
        </div>
      </div>

      <div class="modal-actions">
        <button @click="$emit('close')" class="btn btn-secondary">Cancel</button>
        <button @click="confirmCopy" class="btn btn-primary">
          {{ mode === 'move' ? 'Move Here' : 'Copy Here' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Folder, ArrowLeft } from 'lucide-vue-next'
import fileService from '../api/fileService'

const props = defineProps(['fileName', 'sourcePath', 'sourceId', 'sources', 'mode'])
const emit = defineEmits(['close', 'confirm'])

const selectedSource = ref(props.sourceId)
const currentPath = ref('/')
const selectedPath = ref('/')
const folders = ref([])
const loading = ref(false)

const loadFolders = async () => {
  loading.value = true
  try {
    const response = await fileService.list(selectedSource.value, currentPath.value)
    if (response.success) {
      folders.value = response.data.filter(f => f.isDir)
    }
  } catch (error) {
    console.error('Failed to load folders:', error)
  } finally {
    loading.value = false
  }
}

const navigateInto = (folder) => {
  currentPath.value = folder.path
  selectedPath.value = folder.path
  loadFolders()
}

const navigateUp = () => {
  const parts = currentPath.value.split('/').filter(Boolean)
  parts.pop()
  currentPath.value = '/' + parts.join('/')
  selectedPath.value = currentPath.value
  loadFolders()
}

const selectAndCopy = (path) => {
  selectedPath.value = path
  confirmCopy()
}

const confirmCopy = () => {
  emit('confirm', {
    sourceId: selectedSource.value,
    path: currentPath.value
  })
}

watch(selectedSource, () => {
  currentPath.value = '/'
  selectedPath.value = '/'
  loadFolders()
})

loadFolders()
</script>

<style scoped>
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-content { background: var(--bg-secondary); border-radius: 12px; padding: 24px; width: 90%; max-width: 600px; max-height: 80vh; display: flex; flex-direction: column; }
.modal-content h2 { margin-bottom: 20px; font-size: 20px; font-weight: 600; color: var(--text-primary); }
.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 6px; font-size: 14px; font-weight: 500; color: var(--text-primary); }
.form-group select { width: 100%; padding: 10px; border: 1px solid var(--border-color); border-radius: 8px; background: var(--bg-primary); color: var(--text-primary); }
.path-browser { flex: 1; border: 1px solid var(--border-color); border-radius: 8px; overflow: hidden; display: flex; flex-direction: column; }
.path-header { display: flex; align-items: center; gap: 12px; padding: 12px; background: var(--bg-primary); border-bottom: 1px solid var(--border-color); }
.nav-btn { padding: 6px; background: transparent; border: 1px solid var(--border-color); border-radius: 6px; cursor: pointer; display: flex; align-items: center; }
.nav-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.current-path { font-size: 14px; color: var(--text-primary); font-weight: 500; }
.folder-list { flex: 1; overflow-y: auto; padding: 8px; }
.folder-item { display: flex; align-items: center; gap: 10px; padding: 10px 12px; border-radius: 6px; cursor: pointer; transition: background 0.2s; }
.folder-item:hover { background: var(--bg-hover); }
.folder-item.selected { background: var(--accent-color); color: white; }
.folder-item .icon { width: 20px; height: 20px; color: var(--accent-color); }
.folder-item.selected .icon { color: white; }
.empty { padding: 40px; text-align: center; color: var(--text-secondary); }
.modal-actions { display: flex; gap: 10px; justify-content: flex-end; margin-top: 20px; }
.btn { padding: 10px 16px; border: none; border-radius: 8px; cursor: pointer; font-weight: 500; }
.btn-primary { background: var(--accent-color); color: white; }
.btn-secondary { background: var(--bg-hover); color: var(--text-primary); }
.icon { width: 18px; height: 18px; }
</style>