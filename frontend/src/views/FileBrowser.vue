<template>
  <div class="file-browser" :class="{ 'dark-mode': isDarkMode, 'sidebar-pinned': isSidebarPinned }">
    <!-- Sidebar -->
    <aside class="sidebar" :class="{ collapsed: !isSidebarPinned }">
      <!-- Settings Card -->
      <div class="sidebar-card settings-card">
        <button @click="toggleSidebarPin" class="sidebar-btn" :title="isSidebarPinned ? 'Unpin Sidebar' : 'Pin Sidebar'">
          <Pin :class="{ 'pinned': isSidebarPinned }" />
        </button>
        <button @click="toggleDarkMode" class="sidebar-btn" title="Toggle Dark Mode">
          <Moon v-if="!isDarkMode" />
          <Sun v-else />
        </button>
        <button @click="toggleEncryption" class="sidebar-btn" :class="{ 'active': isEncryptionEnabled }" title="Auto Encryption">
          <Lock v-if="isEncryptionEnabled" />
          <Unlock v-else />
        </button>
      </div>

      <!-- Storage Card -->
      <div class="sidebar-card storage-card">
        <h3 class="card-title">Storage</h3>
        <div class="storage-item">
          <div class="storage-header">
            <HardDrive class="storage-icon" />
            <span class="storage-name">Local Storage</span>
          </div>
          <div class="storage-bar">
            <div class="storage-fill" :style="{ width: storagePercent + '%' }"></div>
          </div>
          <div class="storage-info">
            <span>{{ formatSize(storageUsed) }} / {{ formatSize(storageTotal) }}</span>
            <span>{{ storagePercent }}%</span>
          </div>
        </div>
      </div>
    </aside>

    <!-- Toggle Button (when sidebar is collapsed) -->
    <button v-if="!isSidebarPinned" @click="isSidebarPinned = true" class="sidebar-toggle">
      <PanelLeftOpen />
    </button>

    <!-- Main Content -->
    <main class="main-content">
      <!-- Header -->
      <header class="header">
        <div class="header-content">
          <h1 class="title">
            <Folder class="icon" />
            ZxFileBrowser
          </h1>
          <div class="actions">
            <button @click="showUploadModal = true" class="btn btn-primary">
              <Upload class="icon" />
              Upload
            </button>
            <button @click="showCreateModal = true" class="btn btn-secondary">
              <Plus class="icon" />
              New
            </button>
          </div>
        </div>

        <!-- Breadcrumb -->
        <div class="breadcrumb">
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

      <!-- Loading State -->
      <div v-if="loading" class="loading">
        <Loader2 class="spinner" />
        Loading...
      </div>

      <!-- File List -->
      <div v-else class="file-list">
        <div v-if="files.length === 0" class="empty-state">
          <FileQuestion class="empty-icon" />
          <p>This folder is empty</p>
        </div>

        <div
          v-for="file in files"
          :key="file.path"
          class="file-item"
          @click="handleFileClick(file)"
          @contextmenu.prevent="showContextMenu($event, file)"
        >
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
    </main>

    <!-- Context Menu -->
    <div
      v-if="contextMenu.visible"
      ref="contextMenuRef"
      class="context-menu"
      :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
      @click.stop
    >
      <button @click="downloadFile(contextMenu.file)" v-if="!contextMenu.file?.isDir" class="context-item">
        <Download class="icon" />
        Download
      </button>
      <button @click="startRename(contextMenu.file)" class="context-item">
        <Edit2 class="icon" />
        Rename
      </button>
      <button @click="startCopy(contextMenu.file)" class="context-item">
        <Copy class="icon" />
        Copy
      </button>
      <div class="context-divider"></div>
      <button @click="deleteFile(contextMenu.file)" class="context-item danger">
        <Trash2 class="icon" />
        Delete
      </button>
    </div>

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
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  Folder,
  File,
  Upload,
  Plus,
  Home,
  ChevronRight,
  Download,
  Edit2,
  Copy,
  Trash2,
  Loader2,
  FileQuestion,
  Pin,
  Moon,
  Sun,
  Lock,
  Unlock,
  HardDrive,
  PanelLeftOpen,
} from 'lucide-vue-next'
import fileService from '../api/fileService'

const currentPath = ref('/')
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

// Sidebar state
const isSidebarPinned = ref(true)
const isSidebarOpen = ref(true)
const isDarkMode = ref(false)
const isEncryptionEnabled = ref(false)

// Storage info (mock data for now)
// Storage info
const storageUsed = ref(0)
const storageTotal = ref(0)
const storagePercent = computed(() => {
  if (storageTotal.value === 0) return 0
  return Math.round((storageUsed.value / storageTotal.value) * 100)
})

// Load storage info
const loadStorageInfo = async () => {
  try {
    const response = await fileService.getStorageInfo()
    if (response.success) {
      storageTotal.value = response.data.total
      storageUsed.value = response.data.used
    }
  } catch (error) {
    console.error('Failed to load storage info:', error)
  }
}
// Context menu
const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  file: null,
})

const pathParts = computed(() => {
  return currentPath.value.split('/').filter(Boolean)
})

const loadFiles = async () => {
  loading.value = true
  try {
    const response = await fileService.list(currentPath.value)
    if (response.success) {
      files.value = response.data || []
    }
  } catch (error) {
    console.error('Failed to load files:', error)
    alert('Failed to load files')
  } finally {
    loading.value = false
  }
}

const navigateTo = (path) => {
  currentPath.value = path || '/'
  loadFiles()
}

const getPathUpTo = (index) => {
  const parts = pathParts.value.slice(0, index + 1)
  return '/' + parts.join('/')
}

const handleFileClick = (file) => {
  if (file.isDir) {
    navigateTo(file.path)
  } else {
    window.location.href = fileService.getServeUrl(file.path)  // Opens in same tab
  }
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

// Sidebar Functions
const toggleSidebarPin = () => {
  isSidebarPinned.value = !isSidebarPinned.value
  isSidebarOpen.value = isSidebarPinned.value
}


const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value
}

const toggleEncryption = () => {
  isEncryptionEnabled.value = !isEncryptionEnabled.value
  // TODO: Implement encryption logic
  console.log('Encryption:', isEncryptionEnabled.value ? 'Enabled' : 'Disabled')
}

// File Operations
const handleFileSelect = (event) => {
  selectedFile.value = event.target.files[0]
}

const uploadFile = async () => {
  if (!selectedFile.value) return

  try {
    uploadProgress.value = 0
    await fileService.upload(currentPath.value, selectedFile.value, (progress) => {
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
    const newPath = currentPath.value === '/' 
      ? '/' + createName.value 
      : currentPath.value + '/' + createName.value
    
    await fileService.create(newPath, createType.value === 'folder')
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
    await fileService.rename(renameFile.value.path, renameName.value)
    showRenameModal.value = false
    renameFile.value = null
    renameName.value = ''
    loadFiles()
  } catch (error) {
    console.error('Rename failed:', error)
    alert('Failed to rename')
  }
}

const startCopy = async (file) => {
  hideContextMenu()
  const newName = prompt('Enter new name:', file.name + '_copy')
  if (!newName) return

  try {
    const destPath = currentPath.value === '/' 
      ? '/' + newName 
      : currentPath.value + '/' + newName
    
    await fileService.copy(file.path, destPath)
    loadFiles()
  } catch (error) {
    console.error('Copy failed:', error)
    alert('Failed to copy')
  }
}

const deleteFile = async (file) => {
  hideContextMenu()
  if (!confirm(`Are you sure you want to delete "${file.name}"?`)) return

  try {
    await fileService.delete(file.path)
    loadFiles()
  } catch (error) {
    console.error('Delete failed:', error)
    alert('Failed to delete')
  }
}

const downloadFile = (file) => {
  hideContextMenu()
  window.location.href = fileService.getDownloadUrl(file.path)
}

const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const formatDate = (date) => {
  return new Date(date).toLocaleString()
}

onMounted(() => {
  loadFiles()
  loadStorageInfo()
  document.addEventListener('click', hideContextMenu)
})

onUnmounted(() => {
  document.removeEventListener('click', hideContextMenu)
})
</script>

<style scoped>
.file-browser {
  display: flex;
  min-height: 100vh;
  background: var(--bg-primary);
  transition: background 0.3s;
}

/* Dark Mode Variables */
.file-browser {
  --bg-primary: #f9fafb;
  --bg-secondary: #ffffff;
  --bg-hover: #f3f4f6;
  --text-primary: #111827;
  --text-secondary: #6b7280;
  --border-color: #e5e7eb;
  --accent-color: #3b82f6;
}

.file-browser.dark-mode {
  --bg-primary: #111827;
  --bg-secondary: #1f2937;
  --bg-hover: #374151;
  --text-primary: #f9fafb;
  --text-secondary: #9ca3af;
  --border-color: #374151;
  --accent-color: #60a5fa;
}

/* Sidebar */
.sidebar {
  width: 280px;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  transform: translateX(0);
  transition: transform 0.3s;
  display: flex;
  flex-direction: column;
  gap: 20px; /* ← add this */
  padding: 20px; /* optional, keeps inner spacing from edges */
}


.sidebar.collapsed {
  transform: translateX(-100%);
}

.sidebar-card {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 16px;
}

.settings-card {
  display: flex;
  gap: 8px;
  justify-content: space-around;
}

.sidebar-btn {
  flex: 1;
  padding: 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  color: var(--text-secondary);
}

.sidebar-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
  transform: translateY(-2px);
}

.sidebar-btn.active {
  background: var(--accent-color);
  color: white;
  border-color: var(--accent-color);
}

.sidebar-btn .pinned {
  transform: rotate(45deg);
}

.card-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.storage-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.storage-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.storage-icon {
  width: 18px;
  height: 18px;
  color: var(--text-secondary);
}

.storage-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
}

.storage-bar {
  height: 8px;
  background: var(--border-color);
  border-radius: 4px;
  overflow: hidden;
}

.storage-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6, #8b5cf6);
  transition: width 0.3s;
}

.storage-info {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--text-secondary);
}

.sidebar-toggle {
  position: fixed;
  left: 20px;
  top: 20px;
  z-index: 99;
  padding: 10px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  color: var(--text-primary);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* Main Content */
.main-content {
  flex: 1;
  padding: 20px;
  margin-left: 60px;
  transition: margin-left 0.3s;
}

.sidebar-pinned .main-content {
  margin-left: 280px;
}

.header {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  border: 1px solid var(--border-color);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
}

.actions {
  display: flex;
  gap: 10px;
}

.btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-primary {
  background: var(--accent-color);
  color: white;
}

.btn-primary:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.btn-secondary {
  background: var(--bg-hover);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--border-color);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.icon {
  width: 18px;
  height: 18px;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-primary);
  transition: all 0.2s;
}

.breadcrumb-item:hover {
  background: var(--bg-hover);
}

.separator {
  width: 16px;
  height: 16px;
  color: var(--text-secondary);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.spinner {
  width: 40px;
  height: 40px;
  margin-bottom: 10px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.file-list {
  background: var(--bg-secondary);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--border-color);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin-bottom: 10px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  transition: background 0.2s;
}

.file-item:hover {
  background: var(--bg-hover);
}

.file-item:last-child {
  border-bottom: none;
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

/* Context Menu */
.context-menu {
  position: fixed;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 6px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  z-index: 1000;
  min-width: 180px;
}

.context-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-primary);
  transition: background 0.2s;
  text-align: left;
}

.context-item:hover {
  background: var(--bg-hover);
}

.context-item.danger {
  color: #ef4444;
}

.context-item.danger:hover {
  background: #fee2e2;
}

.context-divider {
  height: 1px;
  background: var(--border-color);
  margin: 6px 0;
}

/* Modal */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 24px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-content h2 {
  margin-bottom: 20px;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: var(--accent-color);
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: var(--border-color);
  border-radius: 4px;
  overflow: hidden;
  margin: 16px 0;
}

.progress-fill {
  height: 100%;
  background: var(--accent-color);
  transition: width 0.3s;
}
</style>