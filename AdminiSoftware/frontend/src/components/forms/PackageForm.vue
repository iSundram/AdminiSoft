
<template>
  <form @submit.prevent="submitForm" class="space-y-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Package Name</label>
        <input
          v-model="form.name"
          type="text"
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
        <input
          v-model="form.description"
          type="text"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Disk Quota (GB)</label>
        <input
          v-model="form.disk_quota"
          type="number"
          min="0"
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Bandwidth Quota (GB)</label>
        <input
          v-model="form.bandwidth_quota"
          type="number"
          min="0"
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Email Accounts</label>
        <input
          v-model="form.email_accounts"
          type="number"
          min="-1"
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <small class="text-gray-500">-1 for unlimited</small>
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Databases</label>
        <input
          v-model="form.databases"
          type="number"
          min="-1"
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <small class="text-gray-500">-1 for unlimited</small>
      </div>
    </div>

    <div class="flex justify-end space-x-4">
      <button
        type="button"
        @click="$emit('cancel')"
        class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
      >
        Cancel
      </button>
      <button
        type="submit"
        :disabled="loading"
        class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
      >
        {{ loading ? 'Creating...' : 'Create Package' }}
      </button>
    </div>
  </form>
</template>

<script>
export default {
  name: 'PackageForm',
  emits: ['submit', 'cancel'],
  data() {
    return {
      loading: false,
      form: {
        name: '',
        description: '',
        disk_quota: 1,
        bandwidth_quota: 10,
        email_accounts: 5,
        databases: 2
      }
    }
  },
  methods: {
    async submitForm() {
      this.loading = true
      try {
        this.$emit('submit', this.form)
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
