
<template>
  <div class="bg-white p-6 rounded-lg shadow">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-medium text-gray-900">Notifications</h3>
      <button @click="markAllAsRead" class="text-sm text-indigo-600 hover:text-indigo-500">
        Mark all as read
      </button>
    </div>
    
    <div class="space-y-3">
      <div
        v-for="notification in notifications"
        :key="notification.id"
        :class="[
          'p-3 rounded-lg border transition-colors',
          notification.read ? 'bg-gray-50 border-gray-200' : 'bg-blue-50 border-blue-200'
        ]"
      >
        <div class="flex items-start">
          <div class="flex-shrink-0">
            <div
              :class="[
                'h-8 w-8 rounded-full flex items-center justify-center',
                getNotificationColor(notification.type)
              ]"
            >
              <component :is="getNotificationIcon(notification.type)" class="h-4 w-4 text-white" />
            </div>
          </div>
          
          <div class="ml-3 flex-1">
            <p class="text-sm font-medium text-gray-900">{{ notification.title }}</p>
            <p class="text-sm text-gray-500 mt-1">{{ notification.message }}</p>
            <p class="text-xs text-gray-400 mt-2">{{ formatTime(notification.timestamp) }}</p>
          </div>
          
          <div class="ml-4 flex-shrink-0">
            <button
              @click="markAsRead(notification.id)"
              v-if="!notification.read"
              class="text-blue-600 hover:text-blue-500 text-xs"
            >
              Mark as read
            </button>
          </div>
        </div>
      </div>
      
      <div v-if="notifications.length === 0" class="text-center py-4 text-gray-500">
        No notifications
      </div>
    </div>
  </div>
</template>

<script>
import {
  BellIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
  InformationCircleIcon
} from '@heroicons/vue/24/solid'

export default {
  name: 'Notifications',
  components: {
    BellIcon,
    ExclamationTriangleIcon,
    CheckCircleIcon,
    InformationCircleIcon
  },
  data() {
    return {
      notifications: [
        {
          id: 1,
          type: 'warning',
          title: 'High CPU Usage',
          message: 'Server CPU usage is above 85%',
          timestamp: new Date(Date.now() - 1000 * 60 * 10),
          read: false
        },
        {
          id: 2,
          type: 'success',
          title: 'Backup Completed',
          message: 'Daily backup completed successfully',
          timestamp: new Date(Date.now() - 1000 * 60 * 60 * 2),
          read: false
        },
        {
          id: 3,
          type: 'info',
          title: 'System Update Available',
          message: 'New system update is available for installation',
          timestamp: new Date(Date.now() - 1000 * 60 * 60 * 6),
          read: true
        }
      ]
    }
  },
  methods: {
    getNotificationColor(type) {
      const colors = {
        info: 'bg-blue-500',
        success: 'bg-green-500',
        warning: 'bg-yellow-500',
        error: 'bg-red-500'
      }
      return colors[type] || 'bg-gray-500'
    },
    getNotificationIcon(type) {
      const icons = {
        info: 'InformationCircleIcon',
        success: 'CheckCircleIcon',
        warning: 'ExclamationTriangleIcon',
        error: 'ExclamationTriangleIcon'
      }
      return icons[type] || 'BellIcon'
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
    },
    markAsRead(id) {
      const notification = this.notifications.find(n => n.id === id)
      if (notification) {
        notification.read = true
      }
    },
    markAllAsRead() {
      this.notifications.forEach(n => n.read = true)
    }
  }
}
</script>
