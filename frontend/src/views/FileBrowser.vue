<template>
  <div class="file-browser">
    <!-- Header -->
    <header class="header">
      <div class="header-content">
        <h1 class="title">
          <Folder class="icon" />
          File Browser
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

      <div v-for="file in files" :key="file.path" class="file-item" @click="handleFileClick(file)">
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
        <div class="file-actions" @click.stop>
          <button @click="downloadFile(file)" v-if="!file.isDir" class="action-btn" title="Download">
            <Download class="icon" />
          </button>
          <button @click="startRename(file)" class="action-btn" title="Rename">
            <Edit2 class="icon" />
          </button>
          <button @click="startCopy(file)" class="action-btn" title="Copy">
            <Copy class="icon" />
          </button>
          <button @click="deleteFile(file)" class="action-btn danger" title="Delete">
            <Trash2 class="icon" />
          </button>
        </div>
      </div>
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
import { ref, computed, onMounted } from 'vue'
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
    window.open(fileService.getServeUrl(file.path), '_blank')
  }
}

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
  window.open(fileService.getDownloadUrl(file.path), '_blank')
}

const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const formatDate = (date) => {
  return new Date(date).toLocaleString()
}

onMounted(() => {
  loadFiles()
})
</script>

<style scoped>
.file-browser {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
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
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: #e5e7eb;
  color: #374151;
}

.btn-secondary:hover {
  background: #d1d5db;
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
  background: #f3f4f6;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.2s;
}

.breadcrumb-item:hover {
  background: #e5e7eb;
}

.separator {
  width: 16px;
  height: 16px;
  color: #9ca3af;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #6b7280;
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
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 20px;
  color: #9ca3af;
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
  border-bottom: 1px solid #f3f4f6;
  cursor: pointer;
  transition: background 0.2s;
}

.file-item:hover {
  background: #f9fafb;
}

.file-item:last-child {
  border-bottom: none;
}

.file-icon {
  flex-shrink: 0;
}

.file-icon .folder {
  color: #3b82f6;
  width: 24px;
  height: 24px;
}

.file-icon .file {
  color: #6b7280;
  width: 24px;
  height: 24px;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-weight: 500;
  color: #111827;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-meta {
  font-size: 12px;
  color: #6b7280;
  margin-top: 2px;
}

.file-meta span:not(:last-child)::after {
  content: '•';
  margin: 0 8px;
}

.file-actions {
  display: flex;
  gap: 4px;
}

.action-btn {
  padding: 6px;
  border: none;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  color: #6b7280;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #f3f4f6;
  color: #111827;
}

.action-btn.danger:hover {
  background: #fee2e2;
  color: #dc2626;
}

.action-btn .icon {
  width: 16px;
  height: 16px;
}

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
  background: white;
  border-radius: 8px;
  padding: 24px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal-content h2 {
  margin-bottom: 20px;
  font-size: 20px;
  font-weight: 600;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #3b82f6;
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
  background: #e5e7eb;
  border-radius: 4px;
  overflow: hidden;
  margin: 16px 0;
}

.progress-fill {
  height: 100%;
  background: #3b82f6;
  transition: width 0.3s;
}
</style>