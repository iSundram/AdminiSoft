
<template>
  <div class="bg-white rounded-lg shadow p-6">
    <h3 class="text-lg font-medium text-gray-900 mb-4">System Status</h3>
    
    <div class="space-y-4">
      <div v-for="service in services" :key="service.name" class="flex items-center justify-between">
        <span class="text-sm font-medium text-gray-700">{{ service.name }}</span>
        <div class="flex items-center">
          <div 
            class="w-3 h-3 rounded-full mr-2"
            :class="{
              'bg-green-400': service.status === 'online',
              'bg-red-400': service.status === 'offline',
              'bg-yellow-400': service.status === 'warning'
            }"
          ></div>
          <span 
            class="text-sm capitalize"
            :class="{
              'text-green-600': service.status === 'online',
              'text-red-600': service.status === 'offline',
              'text-yellow-600': service.status === 'warning'
            }"
          >
            {{ service.status }}
          </span>
        </div>
      </div>
    </div>

    <div class="mt-6 pt-4 border-t border-gray-200">
      <div class="text-sm text-gray-500">
        Last updated: {{ lastUpdated }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const services = ref([
  { name: 'Apache HTTP Server', status: 'online' },
  { name: 'MySQL Database', status: 'online' },
  { name: 'DNS Server', status: 'online' },
  { name: 'Mail Server', status: 'warning' },
  { name: 'FTP Server', status: 'online' },
  { name: 'SSH Server', status: 'online' }
])

const lastUpdated = ref('')

onMounted(() => {
  updateTimestamp()
  setInterval(updateTimestamp, 60000) // Update every minute
})

function updateTimestamp() {
  lastUpdated.value = new Date().toLocaleTimeString()
}
</script>
</script>
