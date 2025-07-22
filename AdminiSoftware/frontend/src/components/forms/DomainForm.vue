
<template>
  <form @submit.prevent="submitForm" class="space-y-6">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Domain Name</label>
      <input
        v-model="form.domain"
        type="text"
        required
        placeholder="example.com"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Type</label>
      <select
        v-model="form.type"
        required
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="addon">Addon Domain</option>
        <option value="subdomain">Subdomain</option>
        <option value="parked">Parked Domain</option>
      </select>
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Document Root</label>
      <input
        v-model="form.document_root"
        type="text"
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
        {{ loading ? 'Adding...' : 'Add Domain' }}
      </button>
    </div>
  </form>
</template>

<script>
export default {
  name: 'DomainForm',
  emits: ['submit', 'cancel'],
  data() {
    return {
      loading: false,
      form: {
        domain: '',
        type: 'addon',
        document_root: ''
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
