
<template>
  <form @submit.prevent="submitForm" class="space-y-6">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Domain</label>
      <select
        v-model="form.domain"
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
      <label class="block text-sm font-medium text-gray-700 mb-2">SSL Type</label>
      <select
        v-model="form.type"
        required
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="letsencrypt">Let's Encrypt (Free)</option>
        <option value="self-signed">Self-Signed</option>
        <option value="upload">Upload Certificate</option>
      </select>
    </div>

    <div v-if="form.type === 'upload'">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Certificate</label>
        <textarea
          v-model="form.certificate"
          rows="5"
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        ></textarea>
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Private Key</label>
        <textarea
          v-model="form.private_key"
          rows="5"
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        ></textarea>
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Certificate Chain</label>
        <textarea
          v-model="form.chain"
          rows="5"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        ></textarea>
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
        {{ loading ? 'Installing...' : 'Install SSL' }}
      </button>
    </div>
  </form>
</template>

<script>
export default {
  name: 'SSLForm',
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
      form: {
        domain: '',
        type: 'letsencrypt',
        certificate: '',
        private_key: '',
        chain: ''
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
