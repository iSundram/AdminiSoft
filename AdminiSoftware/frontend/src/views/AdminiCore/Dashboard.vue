
<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-semibold text-gray-900">AdminiCore Dashboard</h1>
      <p class="text-gray-600">System administration and server management</p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <StatsCard
        title="Total Accounts"
        :value="stats.totalAccounts"
        icon="UsersIcon"
        color="blue"
        :change="stats.accountsChange"
      />
      <StatsCard
        title="Active Accounts"
        :value="stats.activeAccounts"
        icon="CheckCircleIcon"
        color="green"
        :change="stats.activeChange"
      />
      <StatsCard
        title="Server Load"
        :value="stats.serverLoad + '%'"
        icon="CpuChipIcon"
        color="yellow"
        :change="stats.loadChange"
      />
      <StatsCard
        title="Disk Usage"
        :value="stats.diskUsage + '%'"
        icon="CircleStackIcon"
        color="purple"
        :change="stats.diskChange"
      />
    </div>

    <!-- Charts and Quick Actions -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- System Status -->
      <div class="lg:col-span-2">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">System Overview</h3>
          <div class="space-y-4">
            <!-- Server Status -->
            <div class="flex items-center justify-between py-3 border-b">
              <div class="flex items-center">
                <div class="h-3 w-3 bg-green-400 rounded-full mr-3"></div>
                <span class="text-sm font-medium">Web Server (Apache)</span>
              </div>
              <span class="text-sm text-gray-500">Running</span>
            </div>
            
            <div class="flex items-center justify-between py-3 border-b">
              <div class="flex items-center">
                <div class="h-3 w-3 bg-green-400 rounded-full mr-3"></div>
                <span class="text-sm font-medium">Database Server (MySQL)</span>
              </div>
              <span class="text-sm text-gray-500">Running</span>
            </div>
            
            <div class="flex items-center justify-between py-3 border-b">
              <div class="flex items-center">
                <div class="h-3 w-3 bg-green-400 rounded-full mr-3"></div>
                <span class="text-sm font-medium">Mail Server (Postfix)</span>
              </div>
              <span class="text-sm text-gray-500">Running</span>
            </div>
            
            <div class="flex items-center justify-between py-3">
              <div class="flex items-center">
                <div class="h-3 w-3 bg-yellow-400 rounded-full mr-3"></div>
                <span class="text-sm font-medium">DNS Server (BIND)</span>
              </div>
              <span class="text-sm text-gray-500">Warning</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div>
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Quick Actions</h3>
          <div class="space-y-3">
            <router-link
              to="/admin/accounts/create"
              class="block w-full text-left px-4 py-2 bg-blue-50 text-blue-700 rounded-md hover:bg-blue-100 transition-colors"
            >
              Create New Account
            </router-link>
            
            <router-link
              to="/admin/packages/create"
              class="block w-full text-left px-4 py-2 bg-green-50 text-green-700 rounded-md hover:bg-green-100 transition-colors"
            >
              Create Package
            </router-link>
            
            <router-link
              to="/admin/security"
              class="block w-full text-left px-4 py-2 bg-yellow-50 text-yellow-700 rounded-md hover:bg-yellow-100 transition-colors"
            >
              Security Center
            </router-link>
            
            <router-link
              to="/admin/backup"
              class="block w-full text-left px-4 py-2 bg-purple-50 text-purple-700 rounded-md hover:bg-purple-100 transition-colors"
            >
              Backup Configuration
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="bg-white rounded-lg shadow">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-medium text-gray-900">Recent Activity</h3>
      </div>
      <div class="divide-y divide-gray-200">
        <div
          v-for="activity in recentActivity"
          :key="activity.id"
          class="px-6 py-4 flex items-center justify-between"
        >
          <div class="flex items-center">
            <div
              class="h-2 w-2 rounded-full mr-3"
              :class="{
                'bg-green-400': activity.type === 'success',
                'bg-yellow-400': activity.type === 'warning',
                'bg-red-400': activity.type === 'error'
              }"
            ></div>
            <div>
              <p class="text-sm font-medium text-gray-900">{{ activity.message }}</p>
              <p class="text-sm text-gray-500">{{ activity.user }}</p>
            </div>
          </div>
          <span class="text-sm text-gray-500">{{ formatTime(activity.timestamp) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import StatsCard from '@/components/common/StatsCard.vue'

export default {
  name: 'AdminDashboard',
  components: {
    StatsCard
  },
  setup() {
    const stats = ref({
      totalAccounts: 0,
      activeAccounts: 0,
      serverLoad: 0,
      diskUsage: 0,
      accountsChange: 0,
      activeChange: 0,
      loadChange: 0,
      diskChange: 0
    })

    const recentActivity = ref([
      {
        id: 1,
        message: 'New account created: user123',
        user: 'admin@example.com',
        type: 'success',
        timestamp: new Date()
      },
      {
        id: 2,
        message: 'Package updated: Standard Hosting',
        user: 'admin@example.com',
        type: 'success',
        timestamp: new Date(Date.now() - 300000)
      },
      {
        id: 3,
        message: 'High disk usage warning',
        user: 'System',
        type: 'warning',
        timestamp: new Date(Date.now() - 600000)
      }
    ])

    const fetchStats = async () => {
      try {
        // Simulate API call
        stats.value = {
          totalAccounts: 150,
          activeAccounts: 142,
          serverLoad: 45,
          diskUsage: 67,
          accountsChange: 5,
          activeChange: 3,
          loadChange: -2,
          diskChange: 12
        }
      } catch (error) {
        console.error('Failed to fetch stats:', error)
      }
    }

    const formatTime = (timestamp) => {
      const now = new Date()
      const diff = now - timestamp
      const minutes = Math.floor(diff / 60000)
      
      if (minutes < 1) return 'Just now'
      if (minutes < 60) return `${minutes}m ago`
      
      const hours = Math.floor(minutes / 60)
      if (hours < 24) return `${hours}h ago`
      
      const days = Math.floor(hours / 24)
      return `${days}d ago`
    }

    onMounted(() => {
      fetchStats()
    })

    return {
      stats,
      recentActivity,
      formatTime
    }
  }
}
</script>
</template>
<template>
  <div class="admin-dashboard">
    <!-- Header -->
    <div class="dashboard-header">
      <h1 class="text-3xl font-bold text-gray-900">AdminiCore Dashboard</h1>
      <p class="text-gray-600">System Administrator Control Panel</p>
    </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <StatsCard
        title="Total Accounts"
        :value="stats.totalAccounts"
        icon="users"
        color="blue"
        :trend="stats.accountsTrend"
      />
      <StatsCard
        title="Active Domains"
        :value="stats.activeDomains"
        icon="globe"
        color="green"
        :trend="stats.domainsTrend"
      />
      <StatsCard
        title="Server Load"
        :value="stats.serverLoad + '%'"
        icon="cpu"
        color="orange"
        :trend="stats.loadTrend"
      />
      <StatsCard
        title="Total Storage"
        :value="formatBytes(stats.totalStorage)"
        icon="hard-drive"
        color="purple"
        :trend="stats.storageTrend"
      />
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Server Status -->
      <div class="lg:col-span-2">
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-xl font-semibold mb-4">Server Status</h2>
          <SystemStatus :status="systemStatus" />
        </div>
      </div>

      <!-- Quick Actions -->
      <div>
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-xl font-semibold mb-4">Quick Actions</h2>
          <QuickActions :actions="quickActions" @action="handleQuickAction" />
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="lg:col-span-2">
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-xl font-semibold mb-4">Recent Activity</h2>
          <RecentActivity :activities="recentActivities" />
        </div>
      </div>

      <!-- Resource Usage -->
      <div>
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-xl font-semibold mb-4">Resource Usage</h2>
          <ResourceUsage :usage="resourceUsage" />
        </div>
      </div>

      <!-- System Performance -->
      <div class="lg:col-span-3">
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-xl font-semibold mb-4">System Performance</h2>
          <Chart :data="performanceData" type="line" />
        </div>
      </div>
    </div>

    <!-- Security Alerts -->
    <div v-if="securityAlerts.length > 0" class="mt-6">
      <div class="bg-red-50 border border-red-200 rounded-lg p-4">
        <h3 class="text-lg font-semibold text-red-800 mb-2">Security Alerts</h3>
        <div class="space-y-2">
          <div
            v-for="alert in securityAlerts"
            :key="alert.id"
            class="flex items-center justify-between"
          >
            <span class="text-red-700">{{ alert.message }}</span>
            <button
              @click="dismissAlert(alert.id)"
              class="text-red-600 hover:text-red-800"
            >
              Dismiss
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useAdminStore } from '@/store/admin'
import StatsCard from '@/components/common/StatsCard.vue'
import SystemStatus from '@/components/widgets/SystemStatus.vue'
import QuickActions from '@/components/widgets/QuickActions.vue'
import RecentActivity from '@/components/widgets/RecentActivity.vue'
import ResourceUsage from '@/components/widgets/ResourceUsage.vue'
import Chart from '@/components/common/Chart.vue'

export default {
  name: 'AdminDashboard',
  components: {
    StatsCard,
    SystemStatus,
    QuickActions,
    RecentActivity,
    ResourceUsage,
    Chart
  },
  setup() {
    const adminStore = useAdminStore()
    
    const stats = ref({
      totalAccounts: 0,
      activeDomains: 0,
      serverLoad: 0,
      totalStorage: 0,
      accountsTrend: 0,
      domainsTrend: 0,
      loadTrend: 0,
      storageTrend: 0
    })

    const systemStatus = ref([])
    const quickActions = ref([
      { id: 'create-account', label: 'Create Account', icon: 'user-plus' },
      { id: 'backup-system', label: 'System Backup', icon: 'download' },
      { id: 'restart-services', label: 'Restart Services', icon: 'refresh' },
      { id: 'view-logs', label: 'View Logs', icon: 'file-text' }
    ])

    const recentActivities = ref([])
    const resourceUsage = ref({})
    const performanceData = ref({})
    const securityAlerts = ref([])

    const formatBytes = (bytes) => {
      if (bytes === 0) return '0 Bytes'
      const k = 1024
      const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    }

    const handleQuickAction = async (actionId) => {
      switch (actionId) {
        case 'create-account':
          await $router.push('/admin/accounts/create')
          break
        case 'backup-system':
          await adminStore.createSystemBackup()
          break
        case 'restart-services':
          await adminStore.restartServices()
          break
        case 'view-logs':
          await $router.push('/admin/logs')
          break
      }
    }

    const dismissAlert = async (alertId) => {
      await adminStore.dismissSecurityAlert(alertId)
      securityAlerts.value = securityAlerts.value.filter(alert => alert.id !== alertId)
    }

    const loadDashboardData = async () => {
      try {
        const data = await adminStore.getDashboardData()
        stats.value = data.stats
        systemStatus.value = data.systemStatus
        recentActivities.value = data.recentActivities
        resourceUsage.value = data.resourceUsage
        performanceData.value = data.performanceData
        securityAlerts.value = data.securityAlerts
      } catch (error) {
        console.error('Failed to load dashboard data:', error)
      }
    }

    onMounted(() => {
      loadDashboardData()
      // Refresh data every 30 seconds
      setInterval(loadDashboardData, 30000)
    })

    return {
      stats,
      systemStatus,
      quickActions,
      recentActivities,
      resourceUsage,
      performanceData,
      securityAlerts,
      formatBytes,
      handleQuickAction,
      dismissAlert
    }
  }
}
</script>

<style scoped>
.admin-dashboard {
  @apply p-6;
}

.dashboard-header {
  @apply mb-8;
}
</style>
