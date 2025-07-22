
<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white shadow rounded-lg p-6">
      <h1 class="text-2xl font-bold text-gray-900">Reseller Dashboard</h1>
      <p class="text-gray-600">Manage your hosting accounts and resources</p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <StatsCard 
        title="Total Accounts" 
        :value="stats.totalAccounts" 
        icon="users"
        color="blue"
      />
      <StatsCard 
        title="Active Accounts" 
        :value="stats.activeAccounts" 
        icon="check-circle"
        color="green"
      />
      <StatsCard 
        title="Disk Usage" 
        :value="stats.diskUsage" 
        icon="server"
        color="purple"
      />
      <StatsCard 
        title="Bandwidth" 
        :value="stats.bandwidth" 
        icon="chart-bar"
        color="orange"
      />
    </div>

    <!-- Main Content -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Recent Accounts -->
      <div class="bg-white shadow rounded-lg p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Recent Accounts</h2>
        <div class="space-y-3">
          <div v-for="account in recentAccounts" :key="account.id" 
               class="flex items-center justify-between p-3 border rounded-lg">
            <div>
              <div class="font-medium">{{ account.username }}</div>
              <div class="text-sm text-gray-500">{{ account.domain }}</div>
            </div>
            <span class="px-2 py-1 text-xs rounded-full"
                  :class="account.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
              {{ account.status }}
            </span>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="bg-white shadow rounded-lg p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Quick Actions</h2>
        <div class="space-y-3">
          <router-link to="/reseller/accounts/create" 
                       class="block w-full bg-blue-600 text-white text-center py-2 px-4 rounded-md hover:bg-blue-700">
            Create New Account
          </router-link>
          <router-link to="/reseller/packages/create" 
                       class="block w-full bg-green-600 text-white text-center py-2 px-4 rounded-md hover:bg-green-700">
            Create Package
          </router-link>
          <router-link to="/reseller/accounts" 
                       class="block w-full bg-gray-600 text-white text-center py-2 px-4 rounded-md hover:bg-gray-700">
            Manage Accounts
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import StatsCard from '@/components/common/StatsCard.vue'

const stats = ref({
  totalAccounts: 42,
  activeAccounts: 38,
  diskUsage: '85.2 GB',
  bandwidth: '1.2 TB'
})

const recentAccounts = ref([
  { id: 1, username: 'john_doe', domain: 'johndoe.com', status: 'active' },
  { id: 2, username: 'jane_smith', domain: 'janesmith.org', status: 'active' },
  { id: 3, username: 'company_xyz', domain: 'company-xyz.net', status: 'suspended' },
])

onMounted(() => {
  // Load dashboard data
})
</script>
</script>
