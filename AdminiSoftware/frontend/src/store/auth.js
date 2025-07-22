
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import authService from '@/services/auth'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)
  const isLoading = ref(false)
  const error = ref(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userRole = computed(() => user.value?.role || null)
  const isAdmin = computed(() => userRole.value === 'admin')
  const isReseller = computed(() => userRole.value === 'reseller')
  const isUser = computed(() => userRole.value === 'user')

  // Actions
  const login = async (credentials) => {
    try {
      isLoading.value = true
      error.value = null

      const response = await authService.login(credentials)
      
      token.value = response.token
      user.value = response.user
      
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))

      return { success: true }
    } catch (err) {
      error.value = err.message || 'Login failed'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const logout = async () => {
    try {
      await authService.logout()
    } catch (err) {
      console.error('Logout error:', err)
    } finally {
      token.value = null
      user.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  }

  const refreshToken = async () => {
    try {
      const response = await authService.refreshToken()
      token.value = response.token
      localStorage.setItem('token', response.token)
      return true
    } catch (err) {
      await logout()
      return false
    }
  }

  const checkAuth = async () => {
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')
    
    if (storedToken && storedUser) {
      token.value = storedToken
      user.value = JSON.parse(storedUser)
      
      try {
        const response = await authService.me()
        user.value = response.user
        localStorage.setItem('user', JSON.stringify(response.user))
      } catch (err) {
        const refreshed = await refreshToken()
        if (!refreshed) {
          await logout()
        }
      }
    }
  }

  const updateProfile = async (profileData) => {
    try {
      isLoading.value = true
      const response = await authService.updateProfile(profileData)
      user.value = { ...user.value, ...response.user }
      localStorage.setItem('user', JSON.stringify(user.value))
      return { success: true }
    } catch (err) {
      error.value = err.message || 'Profile update failed'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const changePassword = async (passwordData) => {
    try {
      isLoading.value = true
      await authService.changePassword(passwordData)
      return { success: true }
    } catch (err) {
      error.value = err.message || 'Password change failed'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const enableTwoFactor = async () => {
    try {
      const response = await authService.enableTwoFactor()
      return {
        success: true,
        secret: response.secret,
        qrCode: response.qr_code
      }
    } catch (err) {
      return {
        success: false,
        error: err.message || 'Failed to enable 2FA'
      }
    }
  }

  const verifyTwoFactor = async (secret, code) => {
    try {
      await authService.verifyTwoFactor(secret, code)
      user.value.two_factor_enabled = true
      localStorage.setItem('user', JSON.stringify(user.value))
      return { success: true }
    } catch (err) {
      return {
        success: false,
        error: err.message || 'Invalid verification code'
      }
    }
  }

  const disableTwoFactor = async (code) => {
    try {
      await authService.disableTwoFactor(code)
      user.value.two_factor_enabled = false
      localStorage.setItem('user', JSON.stringify(user.value))
      return { success: true }
    } catch (err) {
      return {
        success: false,
        error: err.message || 'Failed to disable 2FA'
      }
    }
  }

  const clearError = () => {
    error.value = null
  }

  return {
    // State
    user: computed(() => user.value),
    token: computed(() => token.value),
    isLoading: computed(() => isLoading.value),
    error: computed(() => error.value),
    
    // Getters
    isAuthenticated,
    userRole,
    isAdmin,
    isReseller,
    isUser,
    
    // Actions
    login,
    logout,
    refreshToken,
    checkAuth,
    updateProfile,
    changePassword,
    enableTwoFactor,
    verifyTwoFactor,
    disableTwoFactor,
    clearError
  }
})
