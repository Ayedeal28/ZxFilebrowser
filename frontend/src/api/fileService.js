import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_URL || '/api'

const api = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
})

export const fileService = {
  // List files in a directory
  async list(path = '/') {
    const response = await api.get('/list', { params: { path } })
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
  async create(path, isDir = false) {
    const response = await api.post('/create', { path, isDir })
    return response.data
  },

  // Delete file or folder
  async delete(path) {
    const response = await api.delete('/delete', { data: { path } })
    return response.data
  },

  // Copy file or folder
  async copy(source, destination) {
    const response = await api.post('/copy', { source, destination })
    return response.data
  },

  // Move file or folder
  async move(source, destination) {
    const response = await api.post('/move', { source, destination })
    return response.data
  },

  // Rename file or folder
  async rename(path, newName) {
    const response = await api.post('/rename', { path, newName })
    return response.data
  },

  // Upload file
  async upload(path, file, onProgress) {
    const formData = new FormData()
    formData.append('file', file)

    const response = await api.post('/upload', formData, {
      params: { path },
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

  // Get download URL
  getDownloadUrl(path) {
    return `${API_BASE}/download?path=${encodeURIComponent(path)}`
  },

  // Get serve URL (for viewing)
  getServeUrl(path) {
    return `${API_BASE}/serve?path=${encodeURIComponent(path)}`
  },
}

export default fileService
