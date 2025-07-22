
<template>
  <footer class="bg-white border-t border-gray-200 px-4 py-3 sm:px-6">
    <div class="flex flex-col sm:flex-row justify-between items-center">
      <div class="text-sm text-gray-500">
        © {{ currentYear }} AdminiSoftware. All rights reserved.
      </div>
      <div class="flex items-center space-x-4 mt-2 sm:mt-0">
        <span class="text-xs text-gray-400">Version {{ version }}</span>
        <a href="/docs" class="text-xs text-gray-500 hover:text-gray-700">Documentation</a>
        <a href="/support" class="text-xs text-gray-500 hover:text-gray-700">Support</a>
      </div>
    </div>
  </footer>
</template>

<script>
export default {
  name: 'Footer',
  computed: {
    currentYear() {
      return new Date().getFullYear()
    },
    version() {
      return '1.0.0'
    }
  }
}
</script>
