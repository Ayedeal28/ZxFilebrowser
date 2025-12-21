import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_URL || '/api'

const api = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
})

export const fileService = {
  // List files in a directory
  async list(sourceID, path = '/') {
    const response = await api.get('/list', { params: { source: sourceID, path } })
    return response.data
  },

  // Get file info
  async getInfo(path) {
    const response = await api.get('/info', { params: { path } })
    return response.data
  },

  // Preview file
  async preview(path) {
    const response = await api.get('/preview', { params: { path } })
    return response.data
  },

  // Create file or folder
  async create(sourceID, path, isDir = false) {
    const response = await api.post('/create', { source: sourceID, path, isDir })
    return response.data
  },
  // Delete file or folder
  async delete(sourceID, path) {
    const response = await api.delete('/delete', { data: { source: sourceID, path } })
    return response.data
  },

  // Copy file or folder
  async copy(sourceId, sourcePath, destId, destination) {
    const response = await api.post('/copy', { 
      sourceId, 
      sourcePath, 
      destId, 
      destination 
    })
    return response.data
  },

  // Move file or folder
  async move(sourceId, sourcePath, destId, destination) {
    const response = await api.post('/move', { 
      sourceId, 
      sourcePath, 
      destId, 
      destination 
    })
    return response.data
  },

  // Rename file or folder
  async rename(sourceID, path, newName) {
    const response = await api.post('/rename', { source: sourceID, path, newName })
    return response.data
  },

  // Upload file
  async upload(sourceID, path, file, onProgress) {
    const formData = new FormData()
    formData.append('file', file)

    const response = await api.post('/upload', formData, {
      params: { source: sourceID, path },
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(percentCompleted)
        }
      },
    })
    return response.data
  },

  async getStorageInfo(sourceID) {
    const response = await api.get('/storage', { params: { source: sourceID } })
    return response.data
  },

   async getSettings() {
    const response = await api.get('/settings')
    return response.data
  },

  async saveSettings(settings) {
    const response = await api.post('/settings/save', settings)
    return response.data
  },

  // Add to fileService object
  async getSources() {
    const response = await api.get('/sources')
    return response.data
  },

  // Get download URL
  getDownloadUrl(sourceID, path) {
    return `${API_BASE}/download?source=${sourceID}&path=${encodeURIComponent(path)}`
  },

  // Get serve URL (for viewing)
  getServeUrl(sourceID, path) {
    return `${API_BASE}/serve?source=${sourceID}&path=${encodeURIComponent(path)}`
  },
  
}

export default fileService
