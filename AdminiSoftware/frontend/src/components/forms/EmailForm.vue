
<template>
  <form @submit.prevent="submitForm" class="space-y-6">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Email Address</label>
      <div class="flex">
        <input
          v-model="form.username"
          type="text"
          required
          class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <span class="px-3 py-2 bg-gray-50 border-t border-b border-r border-gray-300 rounded-r-md">
          @{{ selectedDomain }}
        </span>
      </div>
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Domain</label>
      <select
        v-model="selectedDomain"
        required
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="">Select Domain</option>
        <option v-for="domain in domains" :key="domain" :value="domain">
          {{ domain }}
        </option>
      </select>
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Password</label>
      <input
        v-model="form.password"
        type="password"
        required
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Quota (MB)</label>
      <input
        v-model="form.quota"
        type="number"
        min="0"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
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
        {{ loading ? 'Creating...' : 'Create Email Account' }}
      </button>
    </div>
  </form>
</template>

<script>
export default {
  name: 'EmailForm',
  props: {
    domains: {
      type: Array,
      default: () => []
    }
  },
  emits: ['submit', 'cancel'],
  data() {
    return {
      loading: false,
      selectedDomain: '',
      form: {
        username: '',
        password: '',
        quota: 100
      }
    }
  },
  methods: {
    async submitForm() {
      this.loading = true
      try {
        this.$emit('submit', {
          ...this.form,
          email: `${this.form.username}@${this.selectedDomain}`
        })
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
