
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
import { defineStore } from 'pinia'
import { authService } from '@/services/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token'),
    isLoading: false,
    loginAttempts: 0,
    requires2FA: false,
    tempToken: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token && !!state.user,
    userRole: (state) => state.user?.role || null,
    isAdmin: (state) => state.user?.role === 'admin',
    isReseller: (state) => state.user?.role === 'reseller',
    isUser: (state) => state.user?.role === 'user',
    userName: (state) => state.user ? `${state.user.first_name} ${state.user.last_name}` : '',
  },

  actions: {
    async login(credentials) {
      this.isLoading = true
      try {
        const response = await authService.login(credentials)
        
        if (response.requires_2fa) {
          this.requires2FA = true
          this.tempToken = response.temp_token
          return { requires2FA: true }
        }
        
        this.setAuthData(response.token, response.user)
        this.loginAttempts = 0
        this.requires2FA = false
        this.tempToken = null
        
        return { success: true }
      } catch (error) {
        this.loginAttempts++
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async verify2FA(code) {
      this.isLoading = true
      try {
        const response = await authService.verify2FA({
          code,
          temp_token: this.tempToken
        })
        
        this.setAuthData(response.token, response.user)
        this.requires2FA = false
        this.tempToken = null
        
        return { success: true }
      } catch (error) {
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async register(userData) {
      this.isLoading = true
      try {
        const response = await authService.register(userData)
        this.setAuthData(response.token, response.user)
        return { success: true }
      } catch (error) {
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async logout() {
      this.clearAuthData()
      // Optionally call logout endpoint
      try {
        await authService.logout()
      } catch (error) {
        console.error('Logout error:', error)
      }
    },

    async refreshToken() {
      if (!this.token) return false
      
      try {
        const response = await authService.refreshToken()
        this.setAuthData(response.token, response.user)
        return true
      } catch (error) {
        this.clearAuthData()
        return false
      }
    },

    async fetchProfile() {
      if (!this.isAuthenticated) return
      
      try {
        const response = await authService.getProfile()
        this.user = response.user
      } catch (error) {
        console.error('Failed to fetch profile:', error)
        if (error.response?.status === 401) {
          this.clearAuthData()
        }
      }
    },

    async updateProfile(profileData) {
      try {
        const response = await authService.updateProfile(profileData)
        this.user = { ...this.user, ...response.user }
        return { success: true }
      } catch (error) {
        throw error
      }
    },

    async changePassword(passwordData) {
      try {
        await authService.changePassword(passwordData)
        return { success: true }
      } catch (error) {
        throw error
      }
    },

    async enable2FA() {
      try {
        const response = await authService.enable2FA()
        return response
      } catch (error) {
        throw error
      }
    },

    async disable2FA(code) {
      try {
        await authService.disable2FA({ code })
        this.user.two_factor_enabled = false
        return { success: true }
      } catch (error) {
        throw error
      }
    },

    setAuthData(token, user) {
      this.token = token
      this.user = user
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(user))
    },

    clearAuthData() {
      this.token = null
      this.user = null
      this.requires2FA = false
      this.tempToken = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },

    initializeAuth() {
      const token = localStorage.getItem('token')
      const user = localStorage.getItem('user')
      
      if (token && user) {
        try {
          this.token = token
          this.user = JSON.parse(user)
          // Fetch fresh profile data
          this.fetchProfile()
        } catch (error) {
          console.error('Error initializing auth:', error)
          this.clearAuthData()
        }
      }
    },
  },
})
