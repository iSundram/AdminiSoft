
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/services/api'

const user = ref(null)
const token = ref(localStorage.getItem('token'))
const isAuthenticated = computed(() => !!token.value && !!user.value)

export function useAuth() {
  const router = useRouter()

  const login = async (email, password, twoFA = '') => {
    try {
      const response = await api.post('/auth/login', {
        email,
        password,
        two_fa: twoFA
      })

      const { token: authToken, refresh_token, user: userData } = response.data
      
      token.value = authToken
      user.value = userData
      
      localStorage.setItem('token', authToken)
      localStorage.setItem('refresh_token', refresh_token)
      localStorage.setItem('user', JSON.stringify(userData))
      
      return { success: true }
    } catch (error) {
      const message = error.response?.data?.error || 'Login failed'
      return { success: false, error: message }
    }
  }

  const logout = async () => {
    try {
      await api.post('/auth/logout')
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      token.value = null
      user.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user')
    }
  }

  const refreshToken = async () => {
    try {
      const refreshToken = localStorage.getItem('refresh_token')
      if (!refreshToken) throw new Error('No refresh token')

      const response = await api.post('/auth/refresh', {
        refresh_token: refreshToken
      })

      const { token: newToken } = response.data
      token.value = newToken
      localStorage.setItem('token', newToken)
      
      return true
    } catch (error) {
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
      
      // Validate token with backend
      try {
        const response = await api.get('/auth/me')
        user.value = response.data.user
        localStorage.setItem('user', JSON.stringify(response.data.user))
      } catch (error) {
        // Token might be expired, try to refresh
        const refreshed = await refreshToken()
        if (!refreshed) {
          await logout()
        }
      }
    }
  }

  const enableTwoFactor = async () => {
    try {
      const response = await api.post('/auth/2fa/enable')
      return {
        success: true,
        secret: response.data.secret,
        qrCode: response.data.qr_code
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to enable 2FA'
      }
    }
  }

  const verifyTwoFactor = async (secret, code) => {
    try {
      await api.post('/auth/2fa/verify', { secret, code })
      return { success: true }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Invalid verification code'
      }
    }
  }

  return {
    user: computed(() => user.value),
    token: computed(() => token.value),
    isAuthenticated,
    login,
    logout,
    refreshToken,
    checkAuth,
    enableTwoFactor,
    verifyTwoFactor
  }
}
import { computed } from 'vue'
import { useAuthStore } from '@/store/auth'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'

export function useAuth() {
  const authStore = useAuthStore()
  const router = useRouter()
  const toast = useToast()

  const user = computed(() => authStore.user)
  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const isLoading = computed(() => authStore.isLoading)
  const userRole = computed(() => authStore.userRole)
  const isAdmin = computed(() => authStore.isAdmin)
  const isReseller = computed(() => authStore.isReseller)
  const isUser = computed(() => authStore.isUser)

  const login = async (credentials) => {
    try {
      const result = await authStore.login(credentials)
      
      if (result.requires2FA) {
        return result
      }
      
      toast.success('Login successful!')
      
      // Redirect based on role
      const redirectPath = getRedirectPath(authStore.userRole)
      router.push(redirectPath)
      
      return result
    } catch (error) {
      toast.error(error.message || 'Login failed')
      throw error
    }
  }

  const verify2FA = async (code) => {
    try {
      await authStore.verify2FA(code)
      toast.success('Two-factor authentication successful!')
      
      const redirectPath = getRedirectPath(authStore.userRole)
      router.push(redirectPath)
    } catch (error) {
      toast.error(error.message || '2FA verification failed')
      throw error
    }
  }

  const register = async (userData) => {
    try {
      await authStore.register(userData)
      toast.success('Registration successful!')
      router.push('/dashboard')
    } catch (error) {
      toast.error(error.message || 'Registration failed')
      throw error
    }
  }

  const logout = async () => {
    try {
      await authStore.logout()
      toast.success('Logged out successfully')
      router.push('/login')
    } catch (error) {
      console.error('Logout error:', error)
      router.push('/login')
    }
  }

  const updateProfile = async (profileData) => {
    try {
      await authStore.updateProfile(profileData)
      toast.success('Profile updated successfully!')
    } catch (error) {
      toast.error(error.message || 'Failed to update profile')
      throw error
    }
  }

  const changePassword = async (passwordData) => {
    try {
      await authStore.changePassword(passwordData)
      toast.success('Password changed successfully!')
    } catch (error) {
      toast.error(error.message || 'Failed to change password')
      throw error
    }
  }

  const enable2FA = async () => {
    try {
      const result = await authStore.enable2FA()
      toast.success('Two-factor authentication enabled!')
      return result
    } catch (error) {
      toast.error(error.message || 'Failed to enable 2FA')
      throw error
    }
  }

  const disable2FA = async (code) => {
    try {
      await authStore.disable2FA(code)
      toast.success('Two-factor authentication disabled!')
    } catch (error) {
      toast.error(error.message || 'Failed to disable 2FA')
      throw error
    }
  }

  const hasRole = (role) => {
    if (Array.isArray(role)) {
      return role.includes(userRole.value)
    }
    return userRole.value === role
  }

  const hasPermission = (permission) => {
    // Define permission mappings
    const permissions = {
      'admin': ['manage_accounts', 'manage_packages', 'manage_system', 'view_all'],
      'reseller': ['manage_accounts', 'manage_packages', 'view_reseller'],
      'user': ['manage_own_account', 'view_user']
    }
    
    const userPermissions = permissions[userRole.value] || []
    return userPermissions.includes(permission)
  }

  const getRedirectPath = (role) => {
    switch (role) {
      case 'admin':
        return '/admin/dashboard'
      case 'reseller':
        return '/reseller/dashboard'
      case 'user':
      default:
        return '/dashboard'
    }
  }

  return {
    user,
    isAuthenticated,
    isLoading,
    userRole,
    isAdmin,
    isReseller,
    isUser,
    login,
    verify2FA,
    register,
    logout,
    updateProfile,
    changePassword,
    enable2FA,
    disable2FA,
    hasRole,
    hasPermission,
  }
}
