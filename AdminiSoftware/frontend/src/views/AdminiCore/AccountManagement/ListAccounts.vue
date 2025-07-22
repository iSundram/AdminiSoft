
<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900">Account Management</h1>
        <p class="mt-1 text-sm text-gray-600">
          Manage hosting accounts and their settings
        </p>
      </div>
      <router-link
        to="/admin/accounts/create"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
      >
        Create Account
      </router-link>
    </div>

    <!-- Filters and Search -->
    <div class="bg-white shadow rounded-lg p-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <SearchBar
            v-model="searchQuery"
            placeholder="Search accounts..."
            @search="loadAccounts"
          />
        </div>
        <div>
          <select
            v-model="statusFilter"
            @change="loadAccounts"
            class="block w-full pl-3 pr-10 py-2 text-base border border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
          >
            <option value="">All Status</option>
            <option value="active">Active</option>
            <option value="suspended">Suspended</option>
            <option value="pending">Pending</option>
          </select>
        </div>
        <div>
          <select
            v-model="packageFilter"
            @change="loadAccounts"
            class="block w-full pl-3 pr-10 py-2 text-base border border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
          >
            <option value="">All Packages</option>
            <option v-for="pkg in packages" :key="pkg.id" :value="pkg.id">
              {{ pkg.name }}
            </option>
          </select>
        </div>
        <div>
          <button
            @click="resetFilters"
            class="w-full inline-flex items-center justify-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Reset Filters
          </button>
        </div>
      </div>
    </div>

    <!-- Accounts Table -->
    <div class="bg-white shadow overflow-hidden sm:rounded-md">
      <DataTable
        :columns="columns"
        :data="accounts"
        :loading="loading"
        @action="handleAction"
      />
    </div>

    <!-- Pagination -->
    <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
      <div class="flex-1 flex justify-between sm:hidden">
        <button
          @click="previousPage"
          :disabled="currentPage === 1"
          class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
        >
          Previous
        </button>
        <button
          @click="nextPage"
          :disabled="currentPage === totalPages"
          class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
        >
          Next
        </button>
      </div>
      <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
        <div>
          <p class="text-sm text-gray-700">
            Showing {{ (currentPage - 1) * perPage + 1 }} to {{ Math.min(currentPage * perPage, totalItems) }} of {{ totalItems }} results
          </p>
        </div>
        <div>
          <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
            <button
              v-for="page in pageNumbers"
              :key="page"
              @click="goToPage(page)"
              :class="page === currentPage ? 'bg-blue-50 border-blue-500 text-blue-600' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'"
              class="relative inline-flex items-center px-4 py-2 border text-sm font-medium"
            >
              {{ page }}
            </button>
          </nav>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import adminService from '@/services/admin'
import SearchBar from '@/components/common/SearchBar.vue'
import DataTable from '@/components/common/DataTable.vue'
import { useNotifications } from '@/composables/useNotifications'

export default {
  name: 'ListAccounts',
  components: {
    SearchBar,
    DataTable
  },
  setup() {
    const router = useRouter()
    const { addNotification } = useNotifications()

    const accounts = ref([])
    const packages = ref([])
    const loading = ref(false)
    const searchQuery = ref('')
    const statusFilter = ref('')
    const packageFilter = ref('')
    const currentPage = ref(1)
    const perPage = ref(20)
    const totalItems = ref(0)

    const columns = [
      { key: 'username', title: 'Username', sortable: true },
      { key: 'email', title: 'Email', sortable: true },
      { key: 'domain', title: 'Primary Domain' },
      { key: 'package.name', title: 'Package' },
      { 
        key: 'status',
        title: 'Status',
        render: (value) => {
          const colors = {
            active: 'bg-green-100 text-green-800',
            suspended: 'bg-red-100 text-red-800',
            pending: 'bg-yellow-100 text-yellow-800'
          }
          return `<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full ${colors[value] || 'bg-gray-100 text-gray-800'}">${value}</span>`
        }
      },
      { key: 'created_at', title: 'Created', sortable: true },
      {
        key: 'actions',
        title: 'Actions',
        render: (value, row) => `
          <div class="flex space-x-2">
            <button class="text-blue-600 hover:text-blue-900 text-sm" data-action="edit" data-id="${row.id}">Edit</button>
            <button class="text-yellow-600 hover:text-yellow-900 text-sm" data-action="suspend" data-id="${row.id}">
              ${row.status === 'suspended' ? 'Unsuspend' : 'Suspend'}
            </button>
            <button class="text-red-600 hover:text-red-900 text-sm" data-action="delete" data-id="${row.id}">Delete</button>
          </div>
        `
      }
    ]

    const totalPages = computed(() => {
      return Math.ceil(totalItems.value / perPage.value)
    })

    const pageNumbers = computed(() => {
      const pages = []
      const start = Math.max(1, currentPage.value - 2)
      const end = Math.min(totalPages.value, currentPage.value + 2)
      
      for (let i = start; i <= end; i++) {
        pages.push(i)
      }
      return pages
    })

    const loadAccounts = async () => {
      loading.value = true
      try {
        const params = {
          page: currentPage.value,
          per_page: perPage.value,
          search: searchQuery.value,
          status: statusFilter.value,
          package_id: packageFilter.value
        }

        const response = await adminService.getAccounts(params)
        accounts.value = response.data.data
        totalItems.value = response.data.total
      } catch (error) {
        addNotification({
          type: 'error',
          title: 'Error',
          message: 'Failed to load accounts'
        })
      } finally {
        loading.value = false
      }
    }

    const loadPackages = async () => {
      try {
        const response = await adminService.getPackages()
        packages.value = response.data
      } catch (error) {
        console.error('Failed to load packages:', error)
      }
    }

    const handleAction = async (action, id) => {
      try {
        switch (action) {
          case 'edit':
            router.push(`/admin/accounts/${id}/edit`)
            break
          case 'suspend':
            const account = accounts.value.find(a => a.id === parseInt(id))
            if (account.status === 'suspended') {
              await adminService.unsuspendAccount(id)
              addNotification({
                type: 'success',
                title: 'Success',
                message: 'Account unsuspended successfully'
              })
            } else {
              await adminService.suspendAccount(id)
              addNotification({
                type: 'success',
                title: 'Success',
                message: 'Account suspended successfully'
              })
            }
            loadAccounts()
            break
          case 'delete':
            if (confirm('Are you sure you want to delete this account?')) {
              await adminService.deleteAccount(id)
              addNotification({
                type: 'success',
                title: 'Success',
                message: 'Account deleted successfully'
              })
              loadAccounts()
            }
            break
        }
      } catch (error) {
        addNotification({
          type: 'error',
          title: 'Error',
          message: error.response?.data?.message || 'Operation failed'
        })
      }
    }

    const resetFilters = () => {
      searchQuery.value = ''
      statusFilter.value = ''
      packageFilter.value = ''
      currentPage.value = 1
      loadAccounts()
    }

    const goToPage = (page) => {
      currentPage.value = page
      loadAccounts()
    }

    const previousPage = () => {
      if (currentPage.value > 1) {
        currentPage.value--
        loadAccounts()
      }
    }

    const nextPage = () => {
      if (currentPage.value < totalPages.value) {
        currentPage.value++
        loadAccounts()
      }
    }

    onMounted(() => {
      loadAccounts()
      loadPackages()
    })

    return {
      accounts,
      packages,
      loading,
      searchQuery,
      statusFilter,
      packageFilter,
      currentPage,
      totalItems,
      totalPages,
      pageNumbers,
      columns,
      loadAccounts,
      handleAction,
      resetFilters,
      goToPage,
      previousPage,
      nextPage
    }
  }
}
</script>
