
<template>
  <div class="bg-white p-6 rounded-lg shadow">
    <h3 class="text-lg font-medium text-gray-900 mb-4">Resource Usage</h3>
    
    <div class="space-y-4">
      <div v-for="resource in resources" :key="resource.name" class="space-y-2">
        <div class="flex justify-between items-center">
          <span class="text-sm font-medium text-gray-700">{{ resource.name }}</span>
          <span class="text-sm text-gray-500">{{ resource.used }} / {{ resource.total }}</span>
        </div>
        
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div
            class="h-2 rounded-full transition-all duration-300"
            :class="getProgressBarColor(resource.percentage)"
            :style="{ width: `${resource.percentage}%` }"
          ></div>
        </div>
        
        <div class="text-xs text-gray-500">
          {{ resource.percentage.toFixed(1) }}% used
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ResourceUsage',
  props: {
    resources: {
      type: Array,
      default: () => [
        { name: 'CPU', used: '45%', total: '100%', percentage: 45 },
        { name: 'Memory', used: '2.1 GB', total: '4 GB', percentage: 52.5 },
        { name: 'Disk Space', used: '15 GB', total: '50 GB', percentage: 30 },
        { name: 'Bandwidth', used: '120 GB', total: '500 GB', percentage: 24 }
      ]
    }
  },
  methods: {
    getProgressBarColor(percentage) {
      if (percentage >= 90) return 'bg-red-500'
      if (percentage >= 75) return 'bg-yellow-500'
      if (percentage >= 50) return 'bg-blue-500'
      return 'bg-green-500'
    }
  }
}
</script>
