
import axios from 'axios'
import { useAuth } from '@/composables/useAuth'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:5000/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        const { refreshToken } = useAuth()
        const success = await refreshToken()
        
        if (success) {
          const newToken = localStorage.getItem('token')
          originalRequest.headers.Authorization = `Bearer ${newToken}`
          return api(originalRequest)
        }
      } catch (refreshError) {
        const { logout } = useAuth()
        await logout()
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

export default api
import axios from 'axios'
import { useAuthStore } from '@/store/auth'
import { useToast } from 'vue-toastification'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  async (error) => {
    const authStore = useAuthStore()
    const toast = useToast()
    
    if (error.response) {
      const { status, data } = error.response
      
      // Handle authentication errors
      if (status === 401) {
        if (authStore.isAuthenticated) {
          // Try to refresh token
          const refreshed = await authStore.refreshToken()
          if (refreshed) {
            // Retry the original request
            return api.request(error.config)
          } else {
            authStore.clearAuthData()
            window.location.href = '/login'
          }
        }
      }
      
      // Handle rate limiting
      if (status === 429) {
        toast.error('Too many requests. Please slow down.')
      }
      
      // Handle server errors
      if (status >= 500) {
        toast.error('Server error. Please try again later.')
      }
      
      // Return error with message
      const errorMessage = data?.error || data?.message || 'An error occurred'
      return Promise.reject(new Error(errorMessage))
    }
    
    // Network error
    if (error.request) {
      toast.error('Network error. Please check your connection.')
      return Promise.reject(new Error('Network error'))
    }
    
    return Promise.reject(error)
  }
)

export default api

// Helper functions
export const apiHelpers = {
  // Handle file uploads
  uploadFile: (file, endpoint, onProgress = null) => {
    const formData = new FormData()
    formData.append('file', file)
    
    return api.post(endpoint, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress) {
          const percentCompleted = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total
          )
          onProgress(percentCompleted)
        }
      },
    })
  },
  
  // Handle file downloads
  downloadFile: async (endpoint, filename) => {
    try {
      const response = await api.get(endpoint, {
        responseType: 'blob',
      })
      
      const url = window.URL.createObjectURL(new Blob([response]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', filename)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error('Download error:', error)
      throw error
    }
  },
  
  // Format API errors for display
  formatError: (error) => {
    if (error.response?.data?.error) {
      return error.response.data.error
    }
    if (error.message) {
      return error.message
    }
    return 'An unexpected error occurred'
  },
}
