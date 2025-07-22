
<template>
  <div class="w-64 bg-gray-800 min-h-screen">
    <div class="flex flex-col h-full">
      <!-- Logo -->
      <div class="flex items-center justify-center h-16 px-4 bg-gray-900">
        <img :src="logoSrc" :alt="brandName" class="h-8 w-auto" />
        <span class="ml-2 text-white font-semibold">{{ brandName }}</span>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 px-2 py-4 bg-gray-800">
        <div class="space-y-1">
          <template v-for="item in navigationItems" :key="item.name">
            <!-- Single item -->
            <router-link
              v-if="!item.children"
              :to="item.href"
              class="group flex items-center px-2 py-2 text-sm font-medium rounded-md"
              :class="{
                'bg-gray-900 text-white': isActiveRoute(item.href),
                'text-gray-300 hover:bg-gray-700 hover:text-white': !isActiveRoute(item.href)
              }"
            >
              <component
                :is="item.icon"
                class="mr-3 h-5 w-5 flex-shrink-0"
                :class="{
                  'text-gray-300': isActiveRoute(item.href),
                  'text-gray-400 group-hover:text-gray-300': !isActiveRoute(item.href)
                }"
              />
              {{ item.name }}
            </router-link>

            <!-- Expandable item -->
            <div v-else>
              <button
                @click="toggleExpanded(item.name)"
                class="group w-full flex items-center px-2 py-2 text-sm font-medium rounded-md text-gray-300 hover:bg-gray-700 hover:text-white"
              >
                <component
                  :is="item.icon"
                  class="mr-3 h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-gray-300"
                />
                {{ item.name }}
                <ChevronRightIcon
                  class="ml-auto h-4 w-4 transition-transform duration-200"
                  :class="{ 'rotate-90': expandedItems.includes(item.name) }"
                />
              </button>
              
              <div
                v-if="expandedItems.includes(item.name)"
                class="mt-1 space-y-1"
              >
                <router-link
                  v-for="child in item.children"
                  :key="child.name"
                  :to="child.href"
                  class="group flex items-center pl-11 pr-2 py-2 text-sm font-medium rounded-md"
                  :class="{
                    'bg-gray-900 text-white': isActiveRoute(child.href),
                    'text-gray-300 hover:bg-gray-700 hover:text-white': !isActiveRoute(child.href)
                  }"
                >
                  {{ child.name }}
                </router-link>
              </div>
            </div>
          </template>
        </div>
      </nav>

      <!-- Footer -->
      <div class="flex-shrink-0 p-4 bg-gray-900">
        <div class="flex items-center">
          <img
            class="h-8 w-8 rounded-full"
            :src="user.avatar || '/assets/default-avatar.png'"
            :alt="user.name"
          />
          <div class="ml-3">
            <p class="text-sm font-medium text-white">{{ user.name }}</p>
            <p class="text-xs text-gray-400">{{ user.role }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import {
  HomeIcon,
  UsersIcon,
  ServerIcon,
  ShieldCheckIcon,
  CogIcon,
  DocumentTextIcon,
  EnvelopeIcon,
  GlobeAltIcon,
  ChartBarIcon,
  ChevronRightIcon
} from '@heroicons/vue/24/outline'

export default {
  name: 'Sidebar',
  components: {
    ChevronRightIcon
  },
  setup() {
    const route = useRoute()
    const { user } = useAuth()
    const expandedItems = ref([])

    const navigationItems = computed(() => {
      const role = user.value?.role
      let items = []

      if (role === 'admin') {
        items = [
          { name: 'Dashboard', href: '/admin/dashboard', icon: HomeIcon },
          {
            name: 'Account Management',
            icon: UsersIcon,
            children: [
              { name: 'List Accounts', href: '/admin/accounts' },
              { name: 'Create Account', href: '/admin/accounts/create' },
              { name: 'Suspended Accounts', href: '/admin/accounts/suspended' }
            ]
          },
          {
            name: 'Security Center',
            icon: ShieldCheckIcon,
            children: [
              { name: 'Overview', href: '/admin/security' },
              { name: 'Two-Factor Auth', href: '/admin/security/2fa' },
              { name: 'Brute Force Protection', href: '/admin/security/bruteforce' }
            ]
          },
          {
            name: 'DNS Management',
            icon: GlobeAltIcon,
            children: [
              { name: 'DNS Zone Manager', href: '/admin/dns' },
              { name: 'Add DNS Zone', href: '/admin/dns/add' }
            ]
          },
          {
            name: 'Server Configuration',
            icon: ServerIcon,
            children: [
              { name: 'Basic Setup', href: '/admin/server/basic' },
              { name: 'Web Servers', href: '/admin/server/web' },
              { name: 'Service Manager', href: '/admin/server/services' }
            ]
          }
        ]
      } else if (role === 'reseller') {
        items = [
          { name: 'Dashboard', href: '/reseller/dashboard', icon: HomeIcon },
          {
            name: 'Account Management',
            icon: UsersIcon,
            children: [
              { name: 'List Accounts', href: '/reseller/accounts' },
              { name: 'Create Account', href: '/reseller/accounts/create' }
            ]
          },
          {
            name: 'Package Management',
            icon: DocumentTextIcon,
            children: [
              { name: 'List Packages', href: '/reseller/packages' },
              { name: 'Create Package', href: '/reseller/packages/create' }
            ]
          }
        ]
      } else {
        items = [
          { name: 'Dashboard', href: '/dashboard', icon: HomeIcon },
          {
            name: 'File Management',
            icon: DocumentTextIcon,
            children: [
              { name: 'File Manager', href: '/files' },
              { name: 'FTP Accounts', href: '/files/ftp' },
              { name: 'Disk Usage', href: '/files/usage' }
            ]
          },
          {
            name: 'Email Management',
            icon: EnvelopeIcon,
            children: [
              { name: 'Email Accounts', href: '/email' },
              { name: 'Forwarders', href: '/email/forwarders' },
              { name: 'Autoresponders', href: '/email/autoresponders' }
            ]
          },
          {
            name: 'Database Management',
            icon: ServerIcon,
            children: [
              { name: 'MySQL Databases', href: '/databases' },
              { name: 'Database Wizard', href: '/databases/wizard' }
            ]
          },
          {
            name: 'Statistics',
            icon: ChartBarIcon,
            children: [
              { name: 'AWStats', href: '/stats/awstats' },
              { name: 'Bandwidth Usage', href: '/stats/bandwidth' },
              { name: 'Error Logs', href: '/stats/errors' }
            ]
          }
        ]
      }

      return items
    })

    const brandName = computed(() => {
      const role = user.value?.role
      if (role === 'admin') return 'AdminiCore'
      if (role === 'reseller') return 'AdminiReseller'
      return 'AdminiPanel'
    })

    const logoSrc = computed(() => {
      return '/assets/logos/adminisoftware-logo.svg'
    })

    const isActiveRoute = (href) => {
      return route.path === href || route.path.startsWith(href + '/')
    }

    const toggleExpanded = (itemName) => {
      const index = expandedItems.value.indexOf(itemName)
      if (index > -1) {
        expandedItems.value.splice(index, 1)
      } else {
        expandedItems.value.push(itemName)
      }
    }

    return {
      user,
      navigationItems,
      brandName,
      logoSrc,
      expandedItems,
      isActiveRoute,
      toggleExpanded
    }
  }
}
</script>
<template>
  <div class="w-64 bg-gray-800 text-white min-h-screen">
    <div class="p-4">
      <h2 class="text-lg font-semibold mb-4">{{ panelTitle }}</h2>
      
      <nav class="space-y-2">
        <template v-for="item in menuItems" :key="item.name">
          <div v-if="item.children" class="space-y-1">
            <button @click="toggleSubmenu(item.name)" class="w-full flex items-center justify-between px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-700">
              <span>{{ item.name }}</span>
              <ChevronRightIcon :class="[item.expanded ? 'rotate-90' : '', 'h-4 w-4 transition-transform']" />
            </button>
            <div v-if="item.expanded" class="ml-4 space-y-1">
              <router-link v-for="child in item.children" :key="child.name" :to="child.path" class="block px-3 py-2 rounded-md text-sm text-gray-300 hover:bg-gray-700 hover:text-white">
                {{ child.name }}
              </router-link>
            </div>
          </div>
          
          <router-link v-else :to="item.path" class="block px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-700">
            {{ item.name }}
          </router-link>
        </template>
      </nav>
    </div>
  </div>
</template>

<script>
import { ChevronRightIcon } from '@heroicons/vue/24/outline'

export default {
  name: 'Sidebar',
  components: {
    ChevronRightIcon
  },
  data() {
    return {
      expandedMenus: new Set()
    }
  },
  computed: {
    panelTitle() {
      const route = this.$route.path
      if (route.includes('/admin')) return 'AdminiCore'
      if (route.includes('/reseller')) return 'AdminiReseller'
      return 'AdminiPanel'
    },
    menuItems() {
      const route = this.$route.path
      
      if (route.includes('/admin')) {
        return [
          { name: 'Dashboard', path: '/admin/dashboard' },
          {
            name: 'Server Configuration',
            expanded: this.expandedMenus.has('Server Configuration'),
            children: [
              { name: 'Basic Setup', path: '/admin/server/basic' },
              { name: 'Web Servers', path: '/admin/server/web' },
              { name: 'Database Servers', path: '/admin/server/database' }
            ]
          },
          {
            name: 'Account Management',
            expanded: this.expandedMenus.has('Account Management'),
            children: [
              { name: 'List Accounts', path: '/admin/accounts' },
              { name: 'Create Account', path: '/admin/accounts/create' },
              { name: 'Suspended Accounts', path: '/admin/accounts/suspended' }
            ]
          },
          {
            name: 'Security Center',
            expanded: this.expandedMenus.has('Security Center'),
            children: [
              { name: 'Overview', path: '/admin/security' },
              { name: 'Two Factor Auth', path: '/admin/security/2fa' },
              { name: 'Brute Force Protection', path: '/admin/security/brute-force' }
            ]
          }
        ]
      }
      
      if (route.includes('/reseller')) {
        return [
          { name: 'Dashboard', path: '/reseller/dashboard' },
          { name: 'Accounts', path: '/reseller/accounts' },
          { name: 'Packages', path: '/reseller/packages' },
          { name: 'Branding', path: '/reseller/branding' },
          { name: 'Tools', path: '/reseller/tools' }
        ]
      }
      
      return [
        { name: 'Dashboard', path: '/dashboard' },
        {
          name: 'Domain Management',
          expanded: this.expandedMenus.has('Domain Management'),
          children: [
            { name: 'Subdomains', path: '/domains/subdomains' },
            { name: 'Domain Pointers', path: '/domains/pointers' },
            { name: 'DNS Zone Editor', path: '/domains/dns' }
          ]
        },
        {
          name: 'File Management',
          expanded: this.expandedMenus.has('File Management'),
          children: [
            { name: 'File Manager', path: '/files/manager' },
            { name: 'FTP Accounts', path: '/files/ftp' },
            { name: 'Disk Usage', path: '/files/usage' }
          ]
        },
        {
          name: 'Email Management',
          expanded: this.expandedMenus.has('Email Management'),
          children: [
            { name: 'Email Accounts', path: '/email/accounts' },
            { name: 'Forwarders', path: '/email/forwarders' },
            { name: 'Webmail', path: '/email/webmail' }
          ]
        }
      ]
    }
  },
  methods: {
    toggleSubmenu(name) {
      if (this.expandedMenus.has(name)) {
        this.expandedMenus.delete(name)
      } else {
        this.expandedMenus.add(name)
      }
    }
  }
}
</script>
<template>
  <div class="h-full flex flex-col bg-gray-900">
    <!-- Logo section -->
    <div class="flex items-center h-16 px-4 bg-gray-800">
      <img class="h-8 w-auto" src="/assets/logos/adminisoftware-logo.svg" alt="AdminiSoftware" />
      <span class="ml-2 text-white font-semibold">{{ panelName }}</span>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 px-2 py-4 space-y-1 overflow-y-auto">
      <div v-for="section in navigation" :key="section.name" class="mb-6">
        <div class="px-3 py-2 text-xs font-semibold text-gray-400 uppercase tracking-wide">
          {{ section.name }}
        </div>
        
        <div v-for="item in section.items" :key="item.name">
          <!-- Simple navigation item -->
          <router-link
            v-if="!item.children"
            :to="item.path"
            class="group flex items-center px-2 py-2 text-sm font-medium rounded-md transition-colors"
            :class="isActive(item.path) ? 'bg-gray-800 text-white' : 'text-gray-300 hover:bg-gray-700 hover:text-white'"
          >
            <component :is="item.icon" class="mr-3 flex-shrink-0 h-5 w-5" />
            {{ item.name }}
          </router-link>

          <!-- Navigation item with children -->
          <div v-else>
            <button
              @click="toggleSubmenu(item.name)"
              class="group flex items-center w-full px-2 py-2 text-sm font-medium rounded-md transition-colors text-gray-300 hover:bg-gray-700 hover:text-white"
            >
              <component :is="item.icon" class="mr-3 flex-shrink-0 h-5 w-5" />
              <span class="flex-1 text-left">{{ item.name }}</span>
              <ChevronRightIcon
                class="ml-auto h-4 w-4 transition-transform"
                :class="openSubmenus.includes(item.name) ? 'rotate-90' : ''"
              />
            </button>
            
            <div
              v-if="openSubmenus.includes(item.name)"
              class="mt-1 ml-8 space-y-1"
            >
              <router-link
                v-for="child in item.children"
                :key="child.name"
                :to="child.path"
                class="group flex items-center px-2 py-2 text-sm font-medium rounded-md transition-colors"
                :class="isActive(child.path) ? 'bg-gray-800 text-white' : 'text-gray-400 hover:bg-gray-700 hover:text-gray-300'"
              >
                {{ child.name }}
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <!-- User info -->
    <div class="flex-shrink-0 px-4 py-3 bg-gray-800">
      <div class="flex items-center">
        <img class="h-8 w-8 rounded-full" :src="user.avatar || '/assets/images/default-avatar.png'" :alt="user.username" />
        <div class="ml-3">
          <p class="text-sm font-medium text-white">{{ user.username }}</p>
          <p class="text-xs text-gray-400">{{ user.role }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { ChevronRightIcon } from '@heroicons/vue/24/outline'
import {
  HomeIcon,
  UserGroupIcon,
  CogIcon,
  ShieldCheckIcon,
  ServerIcon,
  DatabaseIcon,
  EnvelopeIcon,
  LockClosedIcon,
  CloudArrowUpIcon,
  ChartBarIcon,
  GlobeAltIcon,
  DocumentIcon,
  WrenchIcon
} from '@heroicons/vue/24/outline'

export default {
  name: 'Sidebar',
  components: {
    ChevronRightIcon
  },
  props: {
    panel: {
      type: String,
      required: true // 'admin', 'reseller', 'user'
    }
  },
  setup(props) {
    const route = useRoute()
    const { user } = useAuth()
    const openSubmenus = ref([])

    const panelName = computed(() => {
      switch (props.panel) {
        case 'admin': return 'AdminiCore'
        case 'reseller': return 'AdminiReseller'
        case 'user': return 'AdminiPanel'
        default: return 'AdminiSoftware'
      }
    })

    const navigation = computed(() => {
      if (props.panel === 'admin') {
        return [
          {
            name: 'Dashboard',
            items: [
              { name: 'Overview', path: '/admin/dashboard', icon: HomeIcon }
            ]
          },
          {
            name: 'Account Management',
            items: [
              { name: 'List Accounts', path: '/admin/accounts', icon: UserGroupIcon },
              { name: 'Create Account', path: '/admin/accounts/create', icon: UserGroupIcon },
              { name: 'Suspended Accounts', path: '/admin/accounts/suspended', icon: UserGroupIcon }
            ]
          },
          {
            name: 'System',
            items: [
              { name: 'Server Configuration', path: '/admin/server', icon: ServerIcon },
              { name: 'Security Center', path: '/admin/security', icon: ShieldCheckIcon },
              { name: 'DNS Management', path: '/admin/dns', icon: GlobeAltIcon },
              { name: 'Email Management', path: '/admin/email', icon: EnvelopeIcon },
              { name: 'SSL Management', path: '/admin/ssl', icon: LockClosedIcon },
              { name: 'Backup Management', path: '/admin/backup', icon: CloudArrowUpIcon },
              { name: 'Monitoring', path: '/admin/monitoring', icon: ChartBarIcon }
            ]
          }
        ]
      } else if (props.panel === 'reseller') {
        return [
          {
            name: 'Dashboard',
            items: [
              { name: 'Overview', path: '/reseller/dashboard', icon: HomeIcon }
            ]
          },
          {
            name: 'Management',
            items: [
              { name: 'Accounts', path: '/reseller/accounts', icon: UserGroupIcon },
              { name: 'Packages', path: '/reseller/packages', icon: DatabaseIcon },
              { name: 'Statistics', path: '/reseller/stats', icon: ChartBarIcon }
            ]
          },
          {
            name: 'Customization',
            items: [
              { name: 'Branding', path: '/reseller/branding', icon: CogIcon }
            ]
          }
        ]
      } else {
        return [
          {
            name: 'Dashboard',
            items: [
              { name: 'Overview', path: '/user/dashboard', icon: HomeIcon }
            ]
          },
          {
            name: 'Website',
            items: [
              { name: 'File Manager', path: '/user/files', icon: DocumentIcon },
              { name: 'Domains', path: '/user/domains', icon: GlobeAltIcon },
              { name: 'SSL Certificates', path: '/user/ssl', icon: LockClosedIcon }
            ]
          },
          {
            name: 'Email & Database',
            items: [
              { name: 'Email Accounts', path: '/user/email', icon: EnvelopeIcon },
              { name: 'Databases', path: '/user/databases', icon: DatabaseIcon }
            ]
          },
          {
            name: 'Applications',
            items: [
              { name: 'WordPress', path: '/user/wordpress', icon: WrenchIcon },
              { name: 'App Installer', path: '/user/apps', icon: WrenchIcon }
            ]
          },
          {
            name: 'Tools',
            items: [
              { name: 'Backup', path: '/user/backup', icon: CloudArrowUpIcon },
              { name: 'Statistics', path: '/user/stats', icon: ChartBarIcon }
            ]
          }
        ]
      }
    })

    const isActive = (path) => {
      return route.path === path || route.path.startsWith(path + '/')
    }

    const toggleSubmenu = (name) => {
      const index = openSubmenus.value.indexOf(name)
      if (index > -1) {
        openSubmenus.value.splice(index, 1)
      } else {
        openSubmenus.value.push(name)
      }
    }

    return {
      user,
      panelName,
      navigation,
      openSubmenus,
      isActive,
      toggleSubmenu
    }
  }
}
</script>
