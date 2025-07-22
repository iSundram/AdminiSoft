
<template>
  <nav class="bg-white shadow-lg">
    <div class="max-w-7xl mx-auto px-4">
      <div class="flex justify-between h-16">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <img class="h-8 w-auto" :src="logoSrc" :alt="brandName" />
          </div>
          <div class="hidden md:block ml-6">
            <div class="flex space-x-4">
              <router-link
                v-for="item in navigation"
                :key="item.name"
                :to="item.href"
                class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
                :class="{ 'text-blue-600': $route.path === item.href }"
              >
                {{ item.name }}
              </router-link>
            </div>
          </div>
        </div>

        <div class="flex items-center space-x-4">
          <!-- Notifications -->
          <button class="relative p-2 text-gray-400 hover:text-gray-500">
            <BellIcon class="h-6 w-6" />
            <span v-if="unreadCount > 0" class="absolute top-0 right-0 h-4 w-4 bg-red-500 text-white rounded-full text-xs flex items-center justify-center">
              {{ unreadCount }}
            </span>
          </button>

          <!-- User menu -->
          <div class="relative">
            <button
              @click="showUserMenu = !showUserMenu"
              class="flex items-center text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <img
                class="h-8 w-8 rounded-full"
                :src="user.avatar || '/assets/default-avatar.png'"
                :alt="user.name"
              />
              <span class="ml-2 text-gray-700">{{ user.name }}</span>
              <ChevronDownIcon class="ml-1 h-4 w-4" />
            </button>

            <div
              v-if="showUserMenu"
              class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-50"
            >
              <router-link
                to="/profile"
                class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
              >
                Profile Settings
              </router-link>
              <router-link
                to="/change-password"
                class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
              >
                Change Password
              </router-link>
              <button
                @click="logout"
                class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
              >
                Sign out
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { BellIcon, ChevronDownIcon } from '@heroicons/vue/24/outline'

export default {
  name: 'Navbar',
  components: {
    BellIcon,
    ChevronDownIcon
  },
  setup() {
    const router = useRouter()
    const { user, logout: authLogout } = useAuth()
    const showUserMenu = ref(false)
    const unreadCount = ref(0)

    const navigation = computed(() => {
      const role = user.value?.role
      const baseNavigation = []

      if (role === 'admin') {
        baseNavigation.push(
          { name: 'AdminiCore', href: '/admin/dashboard' },
          { name: 'AdminiReseller', href: '/reseller/dashboard' },
          { name: 'AdminiPanel', href: '/dashboard' }
        )
      } else if (role === 'reseller') {
        baseNavigation.push(
          { name: 'AdminiReseller', href: '/reseller/dashboard' },
          { name: 'AdminiPanel', href: '/dashboard' }
        )
      } else {
        baseNavigation.push(
          { name: 'AdminiPanel', href: '/dashboard' }
        )
      }

      return baseNavigation
    })

    const logoSrc = computed(() => {
      return '/assets/logos/adminisoftware-logo.svg'
    })

    const brandName = computed(() => {
      const role = user.value?.role
      if (role === 'admin') return 'AdminiCore'
      if (role === 'reseller') return 'AdminiReseller'
      return 'AdminiPanel'
    })

    const logout = async () => {
      try {
        await authLogout()
        router.push('/login')
      } catch (error) {
        console.error('Logout failed:', error)
      }
    }

    const handleClickOutside = (event) => {
      if (!event.target.closest('.relative')) {
        showUserMenu.value = false
      }
    }

    onMounted(() => {
      document.addEventListener('click', handleClickOutside)
    })

    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      user,
      navigation,
      logoSrc,
      brandName,
      showUserMenu,
      unreadCount,
      logout
    }
  }
}
</script>
<template>
  <nav class="bg-white shadow-lg border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex items-center">
          <img class="h-8 w-auto" :src="logoSrc" :alt="siteName" />
          <span class="ml-2 text-xl font-semibold text-gray-900">{{ siteName }}</span>
        </div>
        
        <div class="flex items-center space-x-4">
          <div class="relative">
            <button @click="showUserMenu = !showUserMenu" class="flex items-center text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
              <span class="sr-only">Open user menu</span>
              <div class="h-8 w-8 rounded-full bg-indigo-500 flex items-center justify-center">
                <span class="text-white font-medium">{{ userInitials }}</span>
              </div>
            </button>
            
            <div v-if="showUserMenu" class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white ring-1 ring-black ring-opacity-5 focus:outline-none z-50">
              <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Profile</a>
              <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Settings</a>
              <a @click="logout" href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Sign out</a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import { useAuthStore } from '@/store/auth'

export default {
  name: 'Navbar',
  data() {
    return {
      showUserMenu: false
    }
  },
  computed: {
    siteName() {
      const route = this.$route.path
      if (route.includes('/admin')) return 'AdminiCore'
      if (route.includes('/reseller')) return 'AdminiReseller'
      return 'AdminiPanel'
    },
    logoSrc() {
      const route = this.$route.path
      if (route.includes('/admin')) return '/assets/logos/whm-style-logo.svg'
      if (route.includes('/reseller')) return '/assets/logos/adminisoftware-logo.svg'
      return '/assets/logos/cpanel-style-logo.svg'
    },
    userInitials() {
      const auth = useAuthStore()
      const name = auth.user?.username || 'U'
      return name.substring(0, 2).toUpperCase()
    }
  },
  methods: {
    logout() {
      const auth = useAuthStore()
      auth.logout()
      this.$router.push('/login')
    }
  }
}
</script>
<template>
  <nav class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <img class="h-8 w-auto" src="/assets/logos/adminisoftware-logo.svg" alt="AdminiSoftware" />
          </div>
          <div class="ml-6">
            <h1 class="text-xl font-semibold text-gray-900">{{ pageTitle }}</h1>
          </div>
        </div>

        <div class="flex items-center space-x-4">
          <!-- Notifications -->
          <button
            @click="showNotifications = !showNotifications"
            class="relative p-2 text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 rounded-md"
          >
            <BellIcon class="h-6 w-6" />
            <span v-if="notificationCount > 0" class="absolute -top-1 -right-1 h-4 w-4 bg-red-500 text-white text-xs rounded-full flex items-center justify-center">
              {{ notificationCount }}
            </span>
          </button>

          <!-- User menu -->
          <div class="relative">
            <button
              @click="showUserMenu = !showUserMenu"
              class="flex items-center space-x-3 p-2 text-sm rounded-md text-gray-700 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <img class="h-8 w-8 rounded-full" :src="user.avatar || '/assets/images/default-avatar.png'" :alt="user.username" />
              <span class="hidden md:block">{{ user.username }}</span>
              <ChevronDownIcon class="h-4 w-4" />
            </button>

            <!-- User dropdown -->
            <div
              v-if="showUserMenu"
              class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5 z-50"
            >
              <div class="py-1">
                <router-link
                  to="/profile"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  Profile Settings
                </router-link>
                <router-link
                  to="/preferences"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  Preferences
                </router-link>
                <hr class="my-1">
                <button
                  @click="logout"
                  class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  Sign Out
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Notifications panel -->
    <div
      v-if="showNotifications"
      class="absolute right-4 top-16 w-96 bg-white rounded-lg shadow-lg border border-gray-200 z-50"
    >
      <div class="p-4 border-b border-gray-200">
        <h3 class="text-lg font-medium text-gray-900">Notifications</h3>
      </div>
      <div class="max-h-96 overflow-y-auto">
        <div v-if="notifications.length === 0" class="p-4 text-center text-gray-500">
          No notifications
        </div>
        <div v-else>
          <div
            v-for="notification in notifications"
            :key="notification.id"
            class="p-4 border-b border-gray-100 hover:bg-gray-50"
          >
            <div class="flex items-start">
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900">{{ notification.title }}</p>
                <p class="text-sm text-gray-600 mt-1">{{ notification.message }}</p>
                <p class="text-xs text-gray-400 mt-2">{{ formatDate(notification.created_at) }}</p>
              </div>
              <button
                @click="markAsRead(notification.id)"
                class="ml-2 text-gray-400 hover:text-gray-600"
              >
                <XMarkIcon class="h-4 w-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useNotifications } from '@/composables/useNotifications'
import { BellIcon, ChevronDownIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import dayjs from 'dayjs'

export default {
  name: 'Navbar',
  components: {
    BellIcon,
    ChevronDownIcon,
    XMarkIcon
  },
  props: {
    pageTitle: {
      type: String,
      default: 'Dashboard'
    }
  },
  setup() {
    const router = useRouter()
    const { user, logout: authLogout } = useAuth()
    const { notifications, markAsRead } = useNotifications()
    
    const showUserMenu = ref(false)
    const showNotifications = ref(false)

    const notificationCount = computed(() => {
      return notifications.value.filter(n => !n.read).length
    })

    const logout = async () => {
      await authLogout()
      router.push('/login')
    }

    const formatDate = (date) => {
      return dayjs(date).format('MMM D, h:mm A')
    }

    const handleClickOutside = (event) => {
      if (!event.target.closest('.relative')) {
        showUserMenu.value = false
        showNotifications.value = false
      }
    }

    onMounted(() => {
      document.addEventListener('click', handleClickOutside)
    })

    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      user,
      notifications,
      notificationCount,
      showUserMenu,
      showNotifications,
      logout,
      markAsRead,
      formatDate
    }
  }
}
</script>
