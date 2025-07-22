
<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white shadow rounded-lg p-6">
      <h1 class="text-2xl font-bold text-gray-900">Control Panel</h1>
      <p class="text-gray-600">Welcome back, {{ user?.first_name }}!</p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <StatsCard 
        title="Disk Usage" 
        :value="stats.diskUsage" 
        :max="stats.diskQuota"
        icon="server"
        color="blue"
        type="progress"
      />
      <StatsCard 
        title="Bandwidth" 
        :value="stats.bandwidth" 
        :max="stats.bandwidthQuota"
        icon="chart-bar"
        color="green"
        type="progress"
      />
      <StatsCard 
        title="Email Accounts" 
        :value="stats.emailAccounts" 
        icon="mail"
        color="purple"
      />
      <StatsCard 
        title="Databases" 
        :value="stats.databases" 
        icon="database"
        color="orange"
      />
    </div>

    <!-- Quick Actions Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <div v-for="action in quickActions" :key="action.name"
           class="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow cursor-pointer"
           @click="$router.push(action.route)">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-blue-100 rounded-lg flex items-center justify-center">
              <component :is="action.icon" class="w-5 h-5 text-blue-600" />
            </div>
          </div>
          <div class="ml-4">
            <h3 class="text-sm font-medium text-gray-900">{{ action.name }}</h3>
            <p class="text-xs text-gray-500">{{ action.description }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">Recent Activity</h2>
      <div class="space-y-3">
        <div v-for="activity in recentActivity" :key="activity.id" 
             class="flex items-center p-3 border-l-4 border-blue-400 bg-blue-50">
          <div class="flex-1">
            <p class="text-sm text-gray-900">{{ activity.description }}</p>
            <p class="text-xs text-gray-500">{{ activity.timestamp }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'
import StatsCard from '@/components/common/StatsCard.vue'

const authStore = useAuthStore()
const user = computed(() => authStore.user)

const stats = ref({
  diskUsage: '2.3 GB',
  diskQuota: '10 GB',
  bandwidth: '15.2 GB',
  bandwidthQuota: '100 GB',
  emailAccounts: 8,
  databases: 3
})

const quickActions = ref([
  { name: 'File Manager', description: 'Manage your files', route: '/panel/files', icon: 'FolderIcon' },
  { name: 'Email Accounts', description: 'Manage email', route: '/panel/email', icon: 'MailIcon' },
  { name: 'Databases', description: 'MySQL & PostgreSQL', route: '/panel/databases', icon: 'DatabaseIcon' },
  { name: 'Domain Management', description: 'Subdomains & DNS', route: '/panel/domains', icon: 'GlobeIcon' },
  { name: 'SSL Certificates', description: 'Secure your sites', route: '/panel/ssl', icon: 'ShieldIcon' },
  { name: 'Backup & Restore', description: 'Protect your data', route: '/panel/backup', icon: 'BackupIcon' },
  { name: 'Statistics', description: 'Site analytics', route: '/panel/stats', icon: 'ChartIcon' },
  { name: 'WordPress', description: 'Manage WP sites', route: '/panel/wordpress', icon: 'WordPressIcon' }
])

const recentActivity = ref([
  { id: 1, description: 'SSL certificate renewed for example.com', timestamp: '2 hours ago' },
  { id: 2, description: 'New email account created: support@example.com', timestamp: '1 day ago' },
  { id: 3, description: 'Backup completed successfully', timestamp: '2 days ago' }
])

onMounted(() => {
  // Load dashboard data
})
</script>
</script>
