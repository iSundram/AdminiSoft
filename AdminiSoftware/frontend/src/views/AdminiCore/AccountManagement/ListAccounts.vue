
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
<template>
  <div class="accounts-list">
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Account Management</h1>
        <p class="text-gray-600">Manage user accounts and hosting packages</p>
      </div>
      <button
        @click="$router.push('/admin/accounts/create')"
        class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
      >
        Create Account
      </button>
    </div>

    <!-- Search and Filters -->
    <div class="bg-white rounded-lg shadow p-6 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <SearchBar
          v-model="searchTerm"
          placeholder="Search accounts..."
          @search="handleSearch"
        />
        <select
          v-model="statusFilter"
          @change="applyFilters"
          class="border border-gray-300 rounded-lg px-3 py-2"
        >
          <option value="">All Status</option>
          <option value="active">Active</option>
          <option value="suspended">Suspended</option>
          <option value="pending">Pending</option>
        </select>
        <select
          v-model="packageFilter"
          @change="applyFilters"
          class="border border-gray-300 rounded-lg px-3 py-2"
        >
          <option value="">All Packages</option>
          <option v-for="pkg in packages" :key="pkg.id" :value="pkg.id">
            {{ pkg.name }}
          </option>
        </select>
        <button
          @click="resetFilters"
          class="bg-gray-500 text-white px-4 py-2 rounded-lg hover:bg-gray-600"
        >
          Reset Filters
        </button>
      </div>
    </div>

    <!-- Accounts Table -->
    <div class="bg-white rounded-lg shadow">
      <DataTable
        :data="accounts"
        :columns="columns"
        :loading="loading"
        @sort="handleSort"
        @action="handleAction"
      />
      
      <!-- Pagination -->
      <div class="px-6 py-4 border-t border-gray-200">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-700">
            Showing {{ ((currentPage - 1) * pageSize) + 1 }} to 
            {{ Math.min(currentPage * pageSize, totalAccounts) }} of 
            {{ totalAccounts }} accounts
          </div>
          <div class="flex space-x-2">
            <button
              @click="prevPage"
              :disabled="currentPage === 1"
              class="px-3 py-1 border border-gray-300 rounded disabled:opacity-50"
            >
              Previous
            </button>
            <span class="px-3 py-1">{{ currentPage }}</span>
            <button
              @click="nextPage"
              :disabled="currentPage >= totalPages"
              class="px-3 py-1 border border-gray-300 rounded disabled:opacity-50"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Action Modal -->
    <Modal v-if="showModal" @close="showModal = false">
      <div class="p-6">
        <h3 class="text-lg font-semibold mb-4">{{ modalTitle }}</h3>
        <p class="text-gray-600 mb-6">{{ modalMessage }}</p>
        <div class="flex justify-end space-x-3">
          <button
            @click="showModal = false"
            class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            @click="confirmAction"
            class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
          >
            Confirm
          </button>
        </div>
      </div>
    </Modal>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useAdminStore } from '@/store/admin'
import { useNotifications } from '@/composables/useNotifications'
import DataTable from '@/components/common/DataTable.vue'
import SearchBar from '@/components/common/SearchBar.vue'
import Modal from '@/components/common/Modal.vue'

export default {
  name: 'ListAccounts',
  components: {
    DataTable,
    SearchBar,
    Modal
  },
  setup() {
    const adminStore = useAdminStore()
    const notifications = useNotifications()

    const accounts = ref([])
    const packages = ref([])
    const loading = ref(false)
    const currentPage = ref(1)
    const pageSize = ref(20)
    const totalAccounts = ref(0)
    const searchTerm = ref('')
    const statusFilter = ref('')
    const packageFilter = ref('')
    const sortField = ref('created_at')
    const sortOrder = ref('desc')

    const showModal = ref(false)
    const modalTitle = ref('')
    const modalMessage = ref('')
    const pendingAction = ref(null)

    const columns = [
      { key: 'username', label: 'Username', sortable: true },
      { key: 'email', label: 'Email', sortable: true },
      { key: 'package.name', label: 'Package', sortable: false },
      { key: 'status', label: 'Status', sortable: true },
      { key: 'created_at', label: 'Created', sortable: true },
      { key: 'actions', label: 'Actions', sortable: false }
    ]

    const totalPages = computed(() => {
      return Math.ceil(totalAccounts.value / pageSize.value)
    })

    const loadAccounts = async () => {
      loading.value = true
      try {
        const params = {
          page: currentPage.value,
          limit: pageSize.value,
          search: searchTerm.value,
          status: statusFilter.value,
          package: packageFilter.value,
          sort: sortField.value,
          order: sortOrder.value
        }
        
        const response = await adminStore.getAccounts(params)
        accounts.value = response.accounts
        totalAccounts.value = response.total
      } catch (error) {
        notifications.error('Failed to load accounts')
      } finally {
        loading.value = false
      }
    }

    const loadPackages = async () => {
      try {
        packages.value = await adminStore.getPackages()
      } catch (error) {
        console.error('Failed to load packages:', error)
      }
    }

    const handleSearch = () => {
      currentPage.value = 1
      loadAccounts()
    }

    const applyFilters = () => {
      currentPage.value = 1
      loadAccounts()
    }

    const resetFilters = () => {
      searchTerm.value = ''
      statusFilter.value = ''
      packageFilter.value = ''
      currentPage.value = 1
      loadAccounts()
    }

    const handleSort = (field, order) => {
      sortField.value = field
      sortOrder.value = order
      loadAccounts()
    }

    const handleAction = (action, account) => {
      switch (action) {
        case 'view':
          $router.push(`/admin/accounts/${account.id}`)
          break
        case 'edit':
          $router.push(`/admin/accounts/${account.id}/edit`)
          break
        case 'suspend':
          showConfirmModal(
            'Suspend Account',
            `Are you sure you want to suspend ${account.username}?`,
            () => suspendAccount(account.id)
          )
          break
        case 'unsuspend':
          showConfirmModal(
            'Unsuspend Account',
            `Are you sure you want to unsuspend ${account.username}?`,
            () => unsuspendAccount(account.id)
          )
          break
        case 'delete':
          showConfirmModal(
            'Delete Account',
            `Are you sure you want to delete ${account.username}? This action cannot be undone.`,
            () => deleteAccount(account.id)
          )
          break
      }
    }

    const showConfirmModal = (title, message, action) => {
      modalTitle.value = title
      modalMessage.value = message
      pendingAction.value = action
      showModal.value = true
    }

    const confirmAction = async () => {
      if (pendingAction.value) {
        await pendingAction.value()
        pendingAction.value = null
      }
      showModal.value = false
    }

    const suspendAccount = async (accountId) => {
      try {
        await adminStore.suspendAccount(accountId)
        notifications.success('Account suspended successfully')
        loadAccounts()
      } catch (error) {
        notifications.error('Failed to suspend account')
      }
    }

    const unsuspendAccount = async (accountId) => {
      try {
        await adminStore.unsuspendAccount(accountId)
        notifications.success('Account unsuspended successfully')
        loadAccounts()
      } catch (error) {
        notifications.error('Failed to unsuspend account')
      }
    }

    const deleteAccount = async (accountId) => {
      try {
        await adminStore.deleteAccount(accountId)
        notifications.success('Account deleted successfully')
        loadAccounts()
      } catch (error) {
        notifications.error('Failed to delete account')
      }
    }

    const prevPage = () => {
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
      currentPage,
      pageSize,
      totalAccounts,
      totalPages,
      searchTerm,
      statusFilter,
      packageFilter,
      columns,
      showModal,
      modalTitle,
      modalMessage,
      handleSearch,
      applyFilters,
      resetFilters,
      handleSort,
      handleAction,
      confirmAction,
      prevPage,
      nextPage
    }
  }
}
</script>

<style scoped>
.accounts-list {
  @apply p-6;
}
</style>
