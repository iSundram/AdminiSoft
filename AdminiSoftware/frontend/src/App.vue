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