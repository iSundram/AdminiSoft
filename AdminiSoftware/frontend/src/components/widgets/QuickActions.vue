
<template>
  <div class="bg-white p-6 rounded-lg shadow">
    <h3 class="text-lg font-medium text-gray-900 mb-4">Quick Actions</h3>
    
    <div class="grid grid-cols-2 gap-3">
      <button
        v-for="action in actions"
        :key="action.name"
        @click="handleAction(action)"
        class="flex items-center justify-center px-4 py-3 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
      >
        <component :is="action.icon" class="h-5 w-5 mr-2 text-gray-400" />
        {{ action.name }}
      </button>
    </div>
  </div>
</template>

<script>
import {
  PlusIcon,
  DocumentDuplicateIcon,
  ShieldCheckIcon,
  CogIcon,
  ChartBarIcon,
  UserIcon
} from '@heroicons/vue/24/outline'

export default {
  name: 'QuickActions',
  components: {
    PlusIcon,
    DocumentDuplicateIcon,
    ShieldCheckIcon,
    CogIcon,
    ChartBarIcon,
    UserIcon
  },
  computed: {
    actions() {
      const route = this.$route.path
      
      if (route.includes('/admin')) {
        return [
          { name: 'Create Account', icon: 'UserIcon', action: 'createAccount' },
          { name: 'System Stats', icon: 'ChartBarIcon', action: 'systemStats' },
          { name: 'Security Check', icon: 'ShieldCheckIcon', action: 'securityCheck' },
          { name: 'Server Config', icon: 'CogIcon', action: 'serverConfig' }
        ]
      }
      
      if (route.includes('/reseller')) {
        return [
          { name: 'New Account', icon: 'UserIcon', action: 'newAccount' },
          { name: 'Package Setup', icon: 'PlusIcon', action: 'packageSetup' },
          { name: 'Usage Report', icon: 'ChartBarIcon', action: 'usageReport' },
          { name: 'Branding', icon: 'CogIcon', action: 'branding' }
        ]
      }
      
      return [
        { name: 'File Manager', icon: 'DocumentDuplicateIcon', action: 'fileManager' },
        { name: 'Email Setup', icon: 'PlusIcon', action: 'emailSetup' },
        { name: 'Backup', icon: 'ShieldCheckIcon', action: 'backup' },
        { name: 'Statistics', icon: 'ChartBarIcon', action: 'statistics' }
      ]
    }
  },
  methods: {
    handleAction(action) {
      this.$emit('action', action.action)
      
      // Default actions based on action type
      switch (action.action) {
        case 'createAccount':
          this.$router.push('/admin/accounts/create')
          break
        case 'fileManager':
          this.$router.push('/files/manager')
          break
        case 'systemStats':
          this.$router.push('/admin/monitoring/stats')
          break
        default:
          console.log(`Action: ${action.action}`)
      }
    }
  }
}
</script>
