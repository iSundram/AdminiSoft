
<template>
  <div class="bg-white p-6 rounded-lg shadow">
    <h3 class="text-lg font-medium text-gray-900 mb-4">Recent Activity</h3>
    
    <div class="flow-root">
      <ul class="-mb-8">
        <li v-for="(activity, index) in activities" :key="activity.id" class="relative pb-8">
          <div v-if="index !== activities.length - 1" class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"></div>
          
          <div class="relative flex space-x-3">
            <div>
              <span
                :class="[getActivityColor(activity.type), 'h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white']"
              >
                <component :is="getActivityIcon(activity.type)" class="h-5 w-5 text-white" />
              </span>
            </div>
            
            <div class="min-w-0 flex-1 pt-1.5 flex justify-between space-x-4">
              <div>
                <p class="text-sm text-gray-500">
                  {{ activity.description }}
                  <span class="font-medium text-gray-900">{{ activity.target }}</span>
                </p>
              </div>
              <div class="text-right text-sm whitespace-nowrap text-gray-500">
                {{ formatTime(activity.timestamp) }}
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import {
  UserIcon,
  ShieldCheckIcon,
  DocumentIcon,
  CogIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/solid'

export default {
  name: 'RecentActivity',
  components: {
    UserIcon,
    ShieldCheckIcon,
    DocumentIcon,
    CogIcon,
    ExclamationTriangleIcon
  },
  props: {
    activities: {
      type: Array,
      default: () => [
        {
          id: 1,
          type: 'user',
          description: 'New account created:',
          target: 'john@example.com',
          timestamp: new Date(Date.now() - 1000 * 60 * 5)
        },
        {
          id: 2,
          type: 'security',
          description: 'Security scan completed for',
          target: 'example.com',
          timestamp: new Date(Date.now() - 1000 * 60 * 15)
        },
        {
          id: 3,
          type: 'backup',
          description: 'Backup created for',
          target: 'testsite.com',
          timestamp: new Date(Date.now() - 1000 * 60 * 30)
        },
        {
          id: 4,
          type: 'config',
          description: 'Server configuration updated:',
          target: 'PHP version',
          timestamp: new Date(Date.now() - 1000 * 60 * 60)
        }
      ]
    }
  },
  methods: {
    getActivityColor(type) {
      const colors = {
        user: 'bg-blue-500',
        security: 'bg-green-500',
        backup: 'bg-yellow-500',
        config: 'bg-purple-500',
        error: 'bg-red-500'
      }
      return colors[type] || 'bg-gray-500'
    },
    getActivityIcon(type) {
      const icons = {
        user: 'UserIcon',
        security: 'ShieldCheckIcon',
        backup: 'DocumentIcon',
        config: 'CogIcon',
        error: 'ExclamationTriangleIcon'
      }
      return icons[type] || 'DocumentIcon'
    },
    formatTime(timestamp) {
      const now = new Date()
      const diff = now - timestamp
      const minutes = Math.floor(diff / (1000 * 60))
      
      if (minutes < 1) return 'Just now'
      if (minutes < 60) return `${minutes}m ago`
      
      const hours = Math.floor(minutes / 60)
      if (hours < 24) return `${hours}h ago`
      
      const days = Math.floor(hours / 24)
      return `${days}d ago`
    }
  }
}
</script>
