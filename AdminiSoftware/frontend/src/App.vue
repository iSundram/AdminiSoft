<template>
  <div id="app" :class="themeClass">
    <!-- Login page (no layout) -->
    <router-view v-if="$route.name === 'Login'" />

    <!-- Main app layout -->
    <div v-else class="flex h-screen bg-gray-50">
      <!-- Sidebar -->
      <Sidebar v-if="isAuthenticated" />

      <!-- Main content -->
      <div class="flex-1 flex flex-col overflow-hidden">
        <!-- Navbar -->
        <Navbar v-if="isAuthenticated" />

        <!-- Page content -->
        <main class="flex-1 overflow-x-hidden overflow-y-auto bg-gray-50">
          <div class="container mx-auto px-6 py-8">
            <router-view />
            <Notifications />
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { useTheme } from '@/composables/useTheme'
import Navbar from '@/components/common/Navbar.vue'
import Sidebar from '@/components/common/Sidebar.vue'
import Notifications from '@/components/common/Notifications.vue'

export default {
  name: 'App',
  components: {
    Navbar,
    Sidebar,
    Notifications
  },
  setup() {
    const { isAuthenticated, checkAuth } = useAuth()
    const { theme } = useTheme()

    const themeClass = computed(() => {
      return {
        'whm-theme': theme.value === 'whm',
        'cpanel-theme': theme.value === 'cpanel',
        'dark-theme': theme.value === 'dark'
      }
    })

    onMounted(async () => {
      await checkAuth()
    })

    return {
      isAuthenticated: computed(() => isAuthenticated.value),
      themeClass
    }
  }
}
</script>

<style>
#app {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  min-height: 100vh;
  transition: all 0.3s ease;
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
<template>
  <div id="app" :class="currentTheme">
    <router-view />
    <Teleport to="body">
      <div id="modal-root"></div>
    </Teleport>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'
import { useTheme } from '@/composables/useTheme'

export default {
  name: 'App',
  setup() {
    const authStore = useAuthStore()
    const { currentTheme, initTheme } = useTheme()

    onMounted(() => {
      // Initialize authentication state
      authStore.initializeAuth()
      
      // Initialize theme
      initTheme()
      
      // Set up global error handling
      window.addEventListener('unhandledrejection', (event) => {
        console.error('Unhandled promise rejection:', event.reason)
        // You can add toast notification here
      })
    })

    return {
      currentTheme: computed(() => `theme-${currentTheme.value}`)
    }
  }
}
</script>

<style>
/* Global styles */
html {
  scroll-behavior: smooth;
}

body {
  margin: 0;
  padding: 0;
  font-family: 'Inter', system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  min-height: 100vh;
  background-color: var(--bg-color, #f8fafc);
  transition: background-color 0.3s ease;
}

/* Theme transitions */
.theme-default {
  --bg-color: #f8fafc;
  --text-color: #1e293b;
  --primary-color: #3b82f6;
}

.theme-cpanel {
  --bg-color: #f0f9ff;
  --text-color: #1e293b;
  --primary-color: #059669;
}

.theme-whm {
  --bg-color: #f8fafc;
  --text-color: #334155;
  --primary-color: #64748b;
}

.theme-dark {
  --bg-color: #0f172a;
  --text-color: #f1f5f9;
  --primary-color: #60a5fa;
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f5f9;
}

::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* Loading overlay custom styles */
.vld-overlay {
  background: rgba(255, 255, 255, 0.9) !important;
}

.vld-spinner {
  border-color: var(--primary-color, #3b82f6) !important;
}

/* Toast notification custom styles */
.Vue-Toastification__toast {
  border-radius: 8px !important;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1) !important;
}

.Vue-Toastification__toast--success {
  background-color: #10b981 !important;
}

.Vue-Toastification__toast--error {
  background-color: #ef4444 !important;
}

.Vue-Toastification__toast--warning {
  background-color: #f59e0b !important;
}

.Vue-Toastification__toast--info {
  background-color: #3b82f6 !important;
}
</style>
