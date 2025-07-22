
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
