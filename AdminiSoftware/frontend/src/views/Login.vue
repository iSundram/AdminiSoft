
<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <img class="mx-auto h-12 w-auto" src="/assets/logos/adminisoftware-logo.svg" alt="AdminiSoftware" />
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Sign in to your account
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          AdminiSoftware Control Panel
        </p>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="email" class="sr-only">Email address</label>
            <input
              id="email"
              v-model="form.email"
              name="email"
              type="email"
              autocomplete="email"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="Email address"
            />
          </div>
          <div>
            <label for="password" class="sr-only">Password</label>
            <input
              id="password"
              v-model="form.password"
              name="password"
              type="password"
              autocomplete="current-password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="Password"
            />
          </div>
        </div>

        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              id="remember-me"
              v-model="form.rememberMe"
              name="remember-me"
              type="checkbox"
              class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
            />
            <label for="remember-me" class="ml-2 block text-sm text-gray-900">
              Remember me
            </label>
          </div>
        </div>

        <div v-if="twoFactorRequired" class="space-y-4">
          <div>
            <label for="two-factor-code" class="block text-sm font-medium text-gray-700">
              Two-Factor Authentication Code
            </label>
            <input
              id="two-factor-code"
              v-model="form.twoFactorCode"
              type="text"
              maxlength="6"
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="000000"
            />
          </div>
        </div>

        <div>
          <button
            type="submit"
            :disabled="loading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
          >
            <span v-if="loading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg class="animate-spin h-5 w-5 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            Sign in
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useToast } from 'vue-toastification'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const { login } = useAuth()
    const toast = useToast()
    
    const loading = ref(false)
    const twoFactorRequired = ref(false)
    
    const form = reactive({
      email: '',
      password: '',
      twoFactorCode: '',
      rememberMe: false
    })
    
    const handleLogin = async () => {
      loading.value = true
      
      try {
        const result = await login({
          email: form.email,
          password: form.password,
          twoFactorCode: form.twoFactorCode,
          rememberMe: form.rememberMe
        })
        
        if (result.requiresTwoFactor) {
          twoFactorRequired.value = true
          toast.info('Please enter your 2FA code')
        } else {
          toast.success('Login successful!')
          
          // Redirect based on user role
          const user = result.user
          if (user.role === 'admin') {
            router.push('/admin/dashboard')
          } else if (user.role === 'reseller') {
            router.push('/reseller/dashboard')
          } else {
            router.push('/dashboard')
          }
        }
      } catch (error) {
        toast.error(error.message || 'Login failed')
        twoFactorRequired.value = false
        form.twoFactorCode = ''
      } finally {
        loading.value = false
      }
    }
    
    return {
      form,
      loading,
      twoFactorRequired,
      handleLogin
    }
  }
}
</script>
<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <!-- Logo and Title -->
      <div>
        <img
          class="mx-auto h-12 w-auto"
          src="/assets/logos/adminisoftware-logo.svg"
          alt="AdminiSoftware"
        />
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Sign in to AdminiSoftware
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Control Panel Management System
        </p>
      </div>

      <!-- Login Form -->
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="email" class="sr-only">Email address</label>
            <input
              id="email"
              v-model="form.email"
              name="email"
              type="email"
              autocomplete="email"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="Email address"
            />
          </div>
          <div>
            <label for="password" class="sr-only">Password</label>
            <input
              id="password"
              v-model="form.password"
              name="password"
              type="password"
              autocomplete="current-password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="Password"
            />
          </div>
        </div>

        <!-- Two Factor Authentication -->
        <div v-if="showTwoFactor" class="space-y-4">
          <div>
            <label for="two-factor-code" class="block text-sm font-medium text-gray-700">
              Two-Factor Authentication Code
            </label>
            <input
              id="two-factor-code"
              v-model="form.twoFactorCode"
              type="text"
              maxlength="6"
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              placeholder="123456"
            />
          </div>
        </div>

        <!-- Remember Me -->
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              id="remember-me"
              v-model="form.rememberMe"
              name="remember-me"
              type="checkbox"
              class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
            />
            <label for="remember-me" class="ml-2 block text-sm text-gray-900">
              Remember me
            </label>
          </div>

          <div class="text-sm">
            <a href="#" class="font-medium text-indigo-600 hover:text-indigo-500">
              Forgot your password?
            </a>
          </div>
        </div>

        <!-- Error Messages -->
        <div v-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
          <div class="flex">
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">
                {{ error }}
              </h3>
            </div>
          </div>
        </div>

        <!-- Submit Button -->
        <div>
          <button
            type="submit"
            :disabled="loading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
          >
            <span v-if="loading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <LoadingSpinner class="h-5 w-5 text-indigo-500" />
            </span>
            {{ loading ? 'Signing in...' : 'Sign in' }}
          </button>
        </div>

        <!-- Theme Selector -->
        <div class="mt-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Interface Theme
          </label>
          <select
            v-model="selectedTheme"
            @change="changeTheme"
            class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          >
            <option value="default">Default</option>
            <option value="cpanel">cPanel Style</option>
            <option value="whm">WHM Style</option>
            <option value="dark">Dark Mode</option>
          </select>
        </div>
      </form>

      <!-- Footer -->
      <div class="mt-8 text-center text-xs text-gray-500">
        <p>&copy; 2024 AdminiSoftware. All rights reserved.</p>
        <p class="mt-1">Version 1.0.0</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useTheme } from '@/composables/useTheme'
import { useNotifications } from '@/composables/useNotifications'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'

export default {
  name: 'Login',
  components: {
    LoadingSpinner
  },
  setup() {
    const router = useRouter()
    const { login, isAuthenticated } = useAuth()
    const { theme, setTheme } = useTheme()
    const notifications = useNotifications()

    const form = ref({
      email: '',
      password: '',
      twoFactorCode: '',
      rememberMe: false
    })

    const loading = ref(false)
    const error = ref('')
    const showTwoFactor = ref(false)
    const selectedTheme = ref(theme.value)

    const handleLogin = async () => {
      if (!form.value.email || !form.value.password) {
        error.value = 'Please enter both email and password'
        return
      }

      loading.value = true
      error.value = ''

      try {
        const result = await login({
          email: form.value.email,
          password: form.value.password,
          twoFactorCode: form.value.twoFactorCode,
          rememberMe: form.value.rememberMe
        })

        if (result.requiresTwoFactor) {
          showTwoFactor.value = true
          error.value = 'Please enter your two-factor authentication code'
        } else {
          notifications.success('Logged in successfully')
          
          // Redirect based on user role
          const user = result.user
          if (user.role === 'admin') {
            router.push('/admin/dashboard')
          } else if (user.role === 'reseller') {
            router.push('/reseller/dashboard')
          } else {
            router.push('/dashboard')
          }
        }
      } catch (err) {
        error.value = err.message || 'Login failed. Please try again.'
        showTwoFactor.value = false
      } finally {
        loading.value = false
      }
    }

    const changeTheme = () => {
      setTheme(selectedTheme.value)
      notifications.success(`Theme changed to ${selectedTheme.value}`)
    }

    onMounted(() => {
      // Redirect if already authenticated
      if (isAuthenticated.value) {
        router.push('/')
      }
    })

    return {
      form,
      loading,
      error,
      showTwoFactor,
      selectedTheme,
      handleLogin,
      changeTheme
    }
  }
}
</script>
</template>
