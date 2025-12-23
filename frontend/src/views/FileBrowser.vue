<template>
  <div class="file-browser" :class="{ 'dark-mode': isDarkMode }">
    <Sidebar 
      :is-open="isSidebarOpen"
      :is-pinned="isSidebarPinned"
      :is-dark="isDarkMode"
      :is-encryption="isEncryptionEnabled"
      :sources="sources"
      :active-source="activeSource"
      @close="isSidebarOpen = false"
      @toggle-pin="toggleSidebarPin"
      @toggle-dark="toggleDarkMode"
      @toggle-encryption="toggleEncryption"
      @switch-source="switchSource"
    />

    <button v-if="!isSidebarOpen" @click="isSidebarOpen = true" class="sidebar-toggle">
      <Menu />
    </button>

    <main class="main-content" :class="{'sidebar-collapsed': !isSidebarOpen}">
      <FileHeader 
        :current-path="currentPath"
        :current-view="currentView"
        @upload="showUploadModal = true"
        @create="showCreateModal = true"
        @navigate="navigateTo"
        @change-view="changeView"
      />

      <FileList 
        :files="files"
        :loading="loading"
        :view-mode="currentView"
        :source-id="activeSource"
        @file-click="handleFileClick"
        @context-menu="showContextMenu"
      />
    </main>

    <!-- Context Menu -->
    <ContextMenu 
      :visible="contextMenu.visible"
      :x="contextMenu.x"
      :y="contextMenu.y"
      :file="contextMenu.file"
      @download="downloadFile"
      @rename="startRename"
      @copy="startCopy"
      @move="startMove"
      @delete="deleteFile"
    />

    <!-- Upload Modal -->
    <div v-if="showUploadModal" class="modal" @click.self="showUploadModal = false">
      <div class="modal-content">
        <h2>Upload File</h2>
        <input type="file" @change="handleFileSelect" ref="fileInput" multiple />
        <div v-if="uploadProgress > 0" class="progress-bar">
          <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
        </div>
        <div class="modal-actions">
          <button @click="showUploadModal = false" class="btn btn-secondary">Cancel</button>
          <button @click="uploadFile" :disabled="!selectedFile" class="btn btn-primary">Upload</button>
        </div>
      </div>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="modal" @click.self="showCreateModal = false">
      <div class="modal-content">
        <h2>Create New</h2>
        <div class="form-group">
          <label>Type:</label>
          <select v-model="createType">
            <option value="folder">Folder</option>
            <option value="file">File</option>
          </select>
        </div>
        <div class="form-group">
          <label>Name:</label>
          <input v-model="createName" type="text" placeholder="Enter name" @keyup.enter="createItem" />
        </div>
        <div class="modal-actions">
          <button @click="showCreateModal = false" class="btn btn-secondary">Cancel</button>
          <button @click="createItem" class="btn btn-primary">Create</button>
        </div>
      </div>
    </div>

    <!-- Rename Modal -->
    <div v-if="showRenameModal" class="modal" @click.self="showRenameModal = false">
      <div class="modal-content">
        <h2>Rename</h2>
        <div class="form-group">
          <label>New name:</label>
          <input v-model="renameName" type="text" @keyup.enter="confirmRename" />
        </div>
        <div class="modal-actions">
          <button @click="showRenameModal = false" class="btn btn-secondary">Cancel</button>
          <button @click="confirmRename" class="btn btn-primary">Rename</button>
        </div>
      </div>
    </div>

    <CopyModal 
      v-if="showCopyModal"
      :file-name="copyingFile?.name"
      :source-path="copyingFile?.path"
      :source-id="activeSource"
      :sources="sources"
      :mode="copyMode"
      @close="showCopyModal = false"
      @confirm="confirmCopy"
    />

    <FileViewer 
      :visible="showFileViewer"
      :file-url="viewingFile ? fileService.getServeUrl(activeSource, viewingFile.path) : ''"
      :file-name="viewingFile?.name"
      :download-url="viewingFile ? fileService.getDownloadUrl(activeSource, viewingFile.path) : ''"
      @close="closeViewer"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Menu} from 'lucide-vue-next'
import Sidebar from '@/components/Sidebar.vue'
import FileHeader from '@/components/FileHeader.vue'
import FileList from '@/components/FileList.vue'
import fileService from '../api/fileService'
import ContextMenu from '@/components/ContextMenu.vue'
import CopyModal from '@/components/CopyModal.vue'
import FileViewer from '@/components/FileViewer.vue'

const currentPath = ref('/')
const currentView = ref('list')
const files = ref([])
const loading = ref(false)
const showUploadModal = ref(false)
const showCreateModal = ref(false)
const showRenameModal = ref(false)
const selectedFile = ref(null)
const uploadProgress = ref(0)
const createType = ref('folder')
const createName = ref('')
const renameFile = ref(null)
const renameName = ref('')
const showCopyModal = ref(false)
const copyingFile = ref(null)
const copyMode = ref('copy')
const showFileViewer = ref(false)
const viewingFile = ref(null)


// Sidebar state
const isSidebarPinned = ref(true)
const isSidebarOpen = ref(true)
const isDarkMode = ref(false)
const isEncryptionEnabled = ref(false)

// Storage info
const storageUsed = ref(0)
const storageTotal = ref(0)
const storageName = ref('Local Storage')
const storagePath = ref('')
const sources = ref([])
// const activeSource = ref('')
const activeSource = ref(localStorage.getItem('activeSource') || '')
const storagePercent = computed(() => {
  if (storageTotal.value === 0) return 0
  return Math.round((storageUsed.value / storageTotal.value) * 100)
})

const changeView = (view) => {
  currentView.value = view
  localStorage.setItem('fileView', view)
}

//Copy Menu



// Context menu
const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  file: null,
})

// Load settings
const loadSettings = async () => {
  try {
    const response = await fileService.getSettings()
    if (response.success) {
      isDarkMode.value = response.data.darkMode
      isSidebarPinned.value = response.data.sidebarPin
      isEncryptionEnabled.value = response.data.encryption
      isSidebarOpen.value = response.data.sidebarPin
    }
  } catch (error) {
    console.error('Failed to load settings:', error)
  }
}

// Save settings
const saveSettings = async () => {
  try {
    await fileService.saveSettings({
      darkMode: isDarkMode.value,
      sidebarPin: isSidebarPinned.value,
      encryption: isEncryptionEnabled.value,
    })
  } catch (error) {
    console.error('Failed to save settings:', error)
  }
}

const toggleSidebarPin = () => {
  isSidebarPinned.value = !isSidebarPinned.value
  saveSettings()
}

const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value
  saveSettings()
}

const toggleEncryption = () => {
  isEncryptionEnabled.value = !isEncryptionEnabled.value
  saveSettings()
}

const loadSources = async () => {
  try {
    const response = await fileService.getSources()
    if (response.success) {
      sources.value = response.data.map(src => ({
        ...src,
        used: 0,
        total: 0,
        percent: 0
      }))
      
      // Restore saved source or use first
      const savedSource = localStorage.getItem('activeSource')
      if (savedSource && sources.value.find(s => s.id === savedSource)) {
        activeSource.value = savedSource
      } else if (sources.value.length > 0) {
        activeSource.value = sources.value[0].id
        localStorage.setItem('activeSource', activeSource.value)
      }
      
      for (const source of sources.value) {
        loadStorageInfo(source.id)
      }
    }
  } catch (error) {
    console.error('Failed to load sources:', error)
  }
}

const loadStorageInfo = async (sourceID) => {
  try {
    const response = await fileService.getStorageInfo(sourceID)
    if (response.success) {
      const sourceIndex = sources.value.findIndex(s => s.id === sourceID)
      if (sourceIndex !== -1) {
        sources.value[sourceIndex].total = response.data.total
        sources.value[sourceIndex].used = response.data.used
        sources.value[sourceIndex].percent = Math.round((response.data.used / response.data.total) * 100)
      }
    }
  } catch (error) {
    console.error('Failed to load storage info:', error)
  }
}

const switchSource = (sourceID) => {
  activeSource.value = sourceID
  localStorage.setItem('activeSource', sourceID) // Save to localStorage
  currentPath.value = '/'
  loadFiles()
  loadStorageInfo(sourceID)
}

const loadFiles = async () => {
  loading.value = true
  try {
    const response = await fileService.list(activeSource.value, currentPath.value)
    if (response.success) {
      files.value = response.data || []
    }
  } catch (error) {
    console.error('Failed to load files:', error)
  } finally {
    loading.value = false
  }
}

const navigateTo = (path) => {
  currentPath.value = path || '/'
  loadFiles()
}

const handleFileClick = (file) => {
  if (file.isDir) {
    navigateTo(file.path)
  } else {
    viewingFile.value = file
    showFileViewer.value = true
  }
}

const closeViewer = () => {
  showFileViewer.value = false
  viewingFile.value = null
}

// Context Menu Functions
const showContextMenu = (event, file) => {
  contextMenu.value = {
    visible: true,
    x: event.clientX,
    y: event.clientY,
    file: file,
  }
}

const hideContextMenu = () => {
  contextMenu.value.visible = false
}

// File Operations
const handleFileSelect = (event) => {
  selectedFile.value = event.target.files[0]
}

const uploadFile = async () => {
  if (!selectedFile.value) return
  try {
    await fileService.upload(activeSource.value, currentPath.value, selectedFile.value, (progress) => {
      uploadProgress.value = progress
    })
    showUploadModal.value = false
    selectedFile.value = null
    uploadProgress.value = 0
    loadFiles()
  } catch (error) {
    console.error('Upload failed:', error)
    alert('Upload failed')
  }
}

const createItem = async () => {
  if (!createName.value) return

  try {
    const newPath =
      currentPath.value === '/'
        ? '/' + createName.value
        : currentPath.value + '/' + createName.value

    await fileService.create(
      activeSource.value,
      newPath,
      createType.value === 'folder'
    )

    showCreateModal.value = false
    createName.value = ''
    loadFiles()
  } catch (error) {
    console.error('Create failed:', error)
    alert('Failed to create item')
  }
}


const startRename = (file) => {
  renameFile.value = file
  renameName.value = file.name
  showRenameModal.value = true
  hideContextMenu()
}

const confirmRename = async () => {
  if (!renameName.value || !renameFile.value) return

  try {
    await fileService.rename(activeSource.value, renameFile.value.path, renameName.value)
    showRenameModal.value = false
    renameFile.value = null
    renameName.value = ''
    loadFiles()
  } catch (error) {
    console.error('Rename failed:', error)
    alert('Failed to rename')
  }
}

const startCopy = (file) => {
  hideContextMenu()
  copyingFile.value = file
  copyMode.value = 'copy'
  showCopyModal.value = true
}

const startMove = (file) => {
  hideContextMenu()
  copyingFile.value = file
  copyMode.value = 'move'
  showCopyModal.value = true
}

const confirmCopy = async (destination) => {
  if (!copyingFile.value) return
  
  try {
    if (copyMode.value === 'move') {
      await fileService.move(
        activeSource.value,
        copyingFile.value.path,
        destination.sourceId,
        destination.path + '/' + copyingFile.value.name
      )
      alert('Moved successfully!')
    } else {
      await fileService.copy(
        activeSource.value,
        copyingFile.value.path,
        destination.sourceId,
        destination.path + '/' + copyingFile.value.name
      )
      alert('Copied successfully!')
    }
    
    showCopyModal.value = false
    copyingFile.value = null
    loadFiles()
  } catch (error) {
    console.error(`${copyMode.value} failed:`, error)
    alert(`Failed to ${copyMode.value}`)
  }
}


const deleteFile = async (file) => {
  hideContextMenu()
  if (!confirm(`Are you sure you want to delete "${file.name}"?`)) return

  try {
    await fileService.delete(activeSource.value, file.path)
    loadFiles()
  } catch (error) {
    console.error('Delete failed:', error)
    alert('Failed to delete')
  }
}

const downloadFile = (file) => {
  hideContextMenu()
  window.location.href = fileService.getDownloadUrl(activeSource.value, file.path)
}

onMounted(() => {
    // Restore saved view
  const savedView = localStorage.getItem('fileView')
  if (savedView) {
    currentView.value = savedView
  }
  loadSettings()
  loadSources()
  loadFiles()
  loadStorageInfo()
  document.addEventListener('click', hideContextMenu)
})

onUnmounted(() => {
  document.removeEventListener('click', hideContextMenu)
})
</script>

<style scoped>
.file-browser { display: flex; min-height: 100vh; background: var(--bg-primary); transition: background 0.3s; }
.file-browser { --bg-primary: #f9fafb; --bg-secondary: #ffffff; --bg-hover: #f3f4f6; --text-primary: #111827; --text-secondary: #6b7280; --border-color: #e5e7eb; --accent-color: #3b82f6; }
.file-browser.dark-mode { --bg-primary: #111827; --bg-secondary: #1f2937; --bg-hover: #374151; --text-primary: #f9fafb; --text-secondary: #9ca3af; --border-color: #374151; --accent-color: #60a5fa; }
.sidebar-toggle { position: fixed; left: 10px; top: 20px; z-index: 99; padding: 10px; background: var(--bg-secondary); border: 1px solid var(--border-color); border-radius: 8px; cursor: pointer; color: var(--text-primary); box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1); }
.main-content { flex: 1; padding: 20px; margin-left: 280px; transition: margin-left 0.3s ease; min-height: 100vh;  }
.main-content.sidebar-collapsed { margin-left: 0; padding-left: 70px; }
.icon { width: 18px; height: 18px; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-content { background: var(--bg-secondary); border-radius: 12px; padding: 24px; width: 90%; max-width: 500px; box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3); }
.modal-content h2 { margin-bottom: 20px; font-size: 20px; font-weight: 600; color: var(--text-primary); }
.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 6px; font-size: 14px; font-weight: 500; color: var(--text-primary); }
.form-group input, .form-group select { width: 100%; padding: 10px; border: 1px solid var(--border-color); border-radius: 8px; font-size: 14px; background: var(--bg-primary); color: var(--text-primary); }
.form-group input:focus, .form-group select:focus { outline: none; border-color: var(--accent-color); }
.modal-actions { display: flex; gap: 10px; justify-content: flex-end; margin-top: 20px; }
.btn { display: flex; align-items: center; gap: 8px; padding: 10px 16px; border: none; border-radius: 8px; cursor: pointer; font-size: 14px; font-weight: 500; transition: all 0.2s; }
.btn-primary { background: var(--accent-color); color: white; }
.btn-primary:hover { opacity: 0.9; }
.btn-secondary { background: var(--bg-hover); color: var(--text-primary); border: 1px solid var(--border-color); }
.btn-secondary:hover { background: var(--border-color); }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.progress-bar { width: 100%; height: 8px; background: var(--border-color); border-radius: 4px; overflow: hidden; margin: 16px 0; }
.progress-fill { height: 100%; background: var(--accent-color); transition: width 0.3s; }

</style>