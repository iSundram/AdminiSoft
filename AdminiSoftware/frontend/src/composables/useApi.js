
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/store/auth'

export function useApi() {
  const loading = ref(false)
  const error = ref(null)
  const data = ref(null)

  const state = reactive({
    loading: false,
    error: null,
    data: null
  })

  async function execute(apiCall, options = {}) {
    const { showLoading = true, throwError = true } = options

    if (showLoading) {
      loading.value = true
      state.loading = true
    }

    error.value = null
    state.error = null

    try {
      const result = await apiCall()
      data.value = result.data
      state.data = result.data
      return result
    } catch (err) {
      const errorMessage = err.response?.data?.message || err.message || 'An error occurred'
      error.value = errorMessage
      state.error = errorMessage

      // Handle authentication errors
      if (err.response?.status === 401) {
        const auth = useAuthStore()
        auth.logout()
      }

      if (throwError) {
        throw err
      }
    } finally {
      if (showLoading) {
        loading.value = false
        state.loading = false
      }
    }
  }

  function reset() {
    loading.value = false
    error.value = null
    data.value = null
    state.loading = false
    state.error = null
    state.data = null
  }

  return {
    loading,
    error,
    data,
    state,
    execute,
    reset
  }
}
