
<template>
  <form @submit.prevent="submitForm" class="space-y-6">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Database Type</label>
      <select
        v-model="form.type"
        required
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="mysql">MySQL</option>
        <option value="postgresql">PostgreSQL</option>
        <option value="mongodb">MongoDB</option>
      </select>
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Database Name</label>
      <input
        v-model="form.name"
        type="text"
        required
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Username</label>
      <input
        v-model="form.username"
        type="text"
        required
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
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
        {{ loading ? 'Creating...' : 'Create Database' }}
      </button>
    </div>
  </form>
</template>

<script>
export default {
  name: 'DatabaseForm',
  emits: ['submit', 'cancel'],
  data() {
    return {
      loading: false,
      form: {
        type: 'mysql',
        name: '',
        username: '',
        password: ''
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
