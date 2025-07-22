
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
import api from './api'

export const authService = {
  async login(credentials) {
    const response = await api.post('/auth/login', credentials)
    return response
  },

  async register(userData) {
    const response = await api.post('/auth/register', userData)
    return response
  },

  async logout() {
    const response = await api.post('/auth/logout')
    return response
  },

  async refreshToken() {
    const response = await api.post('/auth/refresh')
    return response
  },

  async getProfile() {
    const response = await api.get('/user/profile')
    return response
  },

  async updateProfile(profileData) {
    const response = await api.put('/user/profile', profileData)
    return response
  },

  async changePassword(passwordData) {
    const response = await api.post('/user/change-password', passwordData)
    return response
  },

  async forgotPassword(email) {
    const response = await api.post('/auth/forgot-password', { email })
    return response
  },

  async resetPassword(resetData) {
    const response = await api.post('/auth/reset-password', resetData)
    return response
  },

  async verify2FA(data) {
    const response = await api.post('/auth/verify-2fa', data)
    return response
  },

  async enable2FA() {
    const response = await api.post('/user/enable-2fa')
    return response
  },

  async disable2FA(data) {
    const response = await api.post('/user/disable-2fa', data)
    return response
  },

  async generateBackupCodes() {
    const response = await api.post('/user/generate-backup-codes')
    return response
  },

  async validateBackupCode(code) {
    const response = await api.post('/user/validate-backup-code', { code })
    return response
  },
}
