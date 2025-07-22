
import api from './api'

export default {
  async login(credentials) {
    const response = await api.post('/auth/login', credentials)
    return response.data
  },

  async register(userData) {
    const response = await api.post('/auth/register', userData)
    return response.data
  },

  async logout() {
    const response = await api.post('/auth/logout')
    return response.data
  },

  async refreshToken() {
    const response = await api.post('/auth/refresh')
    return response.data
  },

  async getCurrentUser() {
    const response = await api.get('/auth/me')
    return response.data
  },

  async updateProfile(userData) {
    const response = await api.put('/auth/profile', userData)
    return response.data
  },

  async changePassword(passwordData) {
    const response = await api.put('/auth/password', passwordData)
    return response.data
  },

  async enable2FA() {
    const response = await api.post('/auth/2fa/enable')
    return response.data
  },

  async disable2FA(code) {
    const response = await api.post('/auth/2fa/disable', { code })
    return response.data
  },

  async verify2FA(code) {
    const response = await api.post('/auth/2fa/verify', { code })
    return response.data
  }
}
