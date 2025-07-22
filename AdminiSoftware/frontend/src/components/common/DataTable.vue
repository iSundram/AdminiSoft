
<template>
  <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
    <!-- Search and filters -->
    <div v-if="searchable || filterable" class="bg-gray-50 px-4 py-3 border-b border-gray-200">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
        <div v-if="searchable" class="relative">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <MagnifyingGlassIcon class="h-5 w-5 text-gray-400" />
          </div>
          <input
            v-model="searchQuery"
            type="search"
            placeholder="Search..."
            class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md text-sm placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
          />
        </div>
        <div v-if="filterable" class="mt-3 sm:mt-0 sm:ml-4">
          <select
            v-model="selectedFilter"
            class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All {{ filterLabel }}</option>
            <option v-for="filter in filters" :key="filter.value" :value="filter.value">
              {{ filter.label }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <!-- Table -->
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th
              v-for="column in columns"
              :key="column.key"
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
              @click="sort(column.key)"
            >
              <div class="flex items-center space-x-1">
                <span>{{ column.label }}</span>
                <div v-if="sortable && column.sortable !== false" class="flex flex-col">
                  <ChevronUpIcon 
                    :class="[
                      'h-3 w-3',
                      sortKey === column.key && sortDirection === 'asc' ? 'text-gray-900' : 'text-gray-400'
                    ]"
                  />
                  <ChevronDownIcon 
                    :class="[
                      'h-3 w-3 -mt-1',
                      sortKey === column.key && sortDirection === 'desc' ? 'text-gray-900' : 'text-gray-400'
                    ]"
                  />
                </div>
              </div>
            </th>
            <th v-if="actions.length > 0" scope="col" class="relative px-6 py-3">
              <span class="sr-only">Actions</span>
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-if="loading" class="animate-pulse">
            <td v-for="column in columns" :key="column.key" class="px-6 py-4">
              <div class="h-4 bg-gray-200 rounded"></div>
            </td>
            <td v-if="actions.length > 0" class="px-6 py-4">
              <div class="h-4 bg-gray-200 rounded w-16"></div>
            </td>
          </tr>
          <tr v-else-if="paginatedData.length === 0">
            <td :colspan="columns.length + (actions.length > 0 ? 1 : 0)" class="px-6 py-4 text-center text-gray-500">
              {{ emptyMessage }}
            </td>
          </tr>
          <tr v-else v-for="(item, index) in paginatedData" :key="index" class="hover:bg-gray-50">
            <td v-for="column in columns" :key="column.key" class="px-6 py-4 whitespace-nowrap">
              <slot :name="`cell-${column.key}`" :item="item" :value="item[column.key]">
                <div v-if="column.type === 'badge'" class="inline-flex px-2 py-1 text-xs font-medium rounded-full"
                     :class="getBadgeClass(item[column.key])">
                  {{ item[column.key] }}
                </div>
                <div v-else-if="column.type === 'date'" class="text-sm text-gray-900">
                  {{ formatDate(item[column.key]) }}
                </div>
                <div v-else-if="column.type === 'bytes'" class="text-sm text-gray-900">
                  {{ formatBytes(item[column.key]) }}
                </div>
                <div v-else class="text-sm text-gray-900">
                  {{ item[column.key] }}
                </div>
              </slot>
            </td>
            <td v-if="actions.length > 0" class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <div class="flex space-x-2">
                <button
                  v-for="action in actions"
                  :key="action.key"
                  @click="$emit('action', { action: action.key, item })"
                  :class="[
                    'inline-flex items-center px-3 py-1 border text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2',
                    action.variant === 'danger' ? 'border-red-300 text-red-700 bg-red-50 hover:bg-red-100 focus:ring-red-500' :
                    action.variant === 'warning' ? 'border-yellow-300 text-yellow-700 bg-yellow-50 hover:bg-yellow-100 focus:ring-yellow-500' :
                    'border-gray-300 text-gray-700 bg-white hover:bg-gray-50 focus:ring-blue-500'
                  ]"
                >
                  {{ action.label }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div v-if="paginated && totalPages > 1" class="bg-white px-4 py-3 border-t border-gray-200 sm:px-6">
      <div class="flex items-center justify-between">
        <div class="flex-1 flex justify-between sm:hidden">
          <button
            @click="previousPage"
            :disabled="currentPage === 1"
            class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <button
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="relative ml-3 inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
        <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-700">
              Showing {{ startItem }} to {{ endItem }} of {{ filteredData.length }} results
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
              <button
                @click="previousPage"
                :disabled="currentPage === 1"
                class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <ChevronLeftIcon class="h-5 w-5" />
              </button>
              <button
                v-for="page in visiblePages"
                :key="page"
                @click="currentPage = page"
                :class="[
                  'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                  page === currentPage
                    ? 'z-10 bg-blue-50 border-blue-500 text-blue-600'
                    : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
                ]"
              >
                {{ page }}
              </button>
              <button
                @click="nextPage"
                :disabled="currentPage === totalPages"
                class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <ChevronRightIcon class="h-5 w-5" />
              </button>
            </nav>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { MagnifyingGlassIcon, ChevronUpIcon, ChevronDownIcon, ChevronLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline'

export default {
  name: 'DataTable',
  components: {
    MagnifyingGlassIcon,
    ChevronUpIcon,
    ChevronDownIcon,
    ChevronLeftIcon,
    ChevronRightIcon
  },
  props: {
    data: {
      type: Array,
      default: () => []
    },
    columns: {
      type: Array,
      required: true
    },
    actions: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    },
    searchable: {
      type: Boolean,
      default: true
    },
    sortable: {
      type: Boolean,
      default: true
    },
    filterable: {
      type: Boolean,
      default: false
    },
    filters: {
      type: Array,
      default: () => []
    },
    filterLabel: {
      type: String,
      default: 'items'
    },
    paginated: {
      type: Boolean,
      default: true
    },
    perPage: {
      type: Number,
      default: 10
    },
    emptyMessage: {
      type: String,
      default: 'No data available'
    }
  },
  emits: ['action'],
  data() {
    return {
      searchQuery: '',
      selectedFilter: '',
      sortKey: '',
      sortDirection: 'asc',
      currentPage: 1
    }
  },
  computed: {
    filteredData() {
      let filtered = [...this.data]

      // Apply search
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase()
        filtered = filtered.filter(item => {
          return this.columns.some(column => {
            const value = item[column.key]
            return value && value.toString().toLowerCase().includes(query)
          })
        })
      }

      // Apply filter
      if (this.selectedFilter) {
        filtered = filtered.filter(item => {
          const filterColumn = this.filters.find(f => f.value === this.selectedFilter)
          return filterColumn && item[filterColumn.key] === this.selectedFilter
        })
      }

      // Apply sort
      if (this.sortKey) {
        filtered.sort((a, b) => {
          const aVal = a[this.sortKey]
          const bVal = b[this.sortKey]
          
          if (aVal < bVal) return this.sortDirection === 'asc' ? -1 : 1
          if (aVal > bVal) return this.sortDirection === 'asc' ? 1 : -1
          return 0
        })
      }

      return filtered
    },
    totalPages() {
      if (!this.paginated) return 1
      return Math.ceil(this.filteredData.length / this.perPage)
    },
    paginatedData() {
      if (!this.paginated) return this.filteredData
      
      const start = (this.currentPage - 1) * this.perPage
      const end = start + this.perPage
      return this.filteredData.slice(start, end)
    },
    startItem() {
      return (this.currentPage - 1) * this.perPage + 1
    },
    endItem() {
      return Math.min(this.currentPage * this.perPage, this.filteredData.length)
    },
    visiblePages() {
      const pages = []
      const total = this.totalPages
      const current = this.currentPage
      const delta = 2

      for (let i = Math.max(1, current - delta); i <= Math.min(total, current + delta); i++) {
        pages.push(i)
      }

      return pages
    }
  },
  methods: {
    sort(key) {
      if (this.sortKey === key) {
        this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc'
      } else {
        this.sortKey = key
        this.sortDirection = 'asc'
      }
    },
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++
      }
    },
    previousPage() {
      if (this.currentPage > 1) {
        this.currentPage--
      }
    },
    getBadgeClass(status) {
      const classes = {
        active: 'bg-green-100 text-green-800',
        inactive: 'bg-gray-100 text-gray-800',
        pending: 'bg-yellow-100 text-yellow-800',
        error: 'bg-red-100 text-red-800',
        success: 'bg-green-100 text-green-800',
        warning: 'bg-yellow-100 text-yellow-800',
        danger: 'bg-red-100 text-red-800'
      }
      return classes[status?.toLowerCase()] || 'bg-gray-100 text-gray-800'
    },
    formatDate(date) {
      if (!date) return ''
      return new Date(date).toLocaleDateString()
    },
    formatBytes(bytes) {
      if (!bytes) return '0 B'
      const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
      const i = Math.floor(Math.log(bytes) / Math.log(1024))
      return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i]
    }
  },
  watch: {
    searchQuery() {
      this.currentPage = 1
    },
    selectedFilter() {
      this.currentPage = 1
    }
  }
}
</script>
