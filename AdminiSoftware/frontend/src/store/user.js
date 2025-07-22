
import { defineStore } from 'pinia'
import userService from '@/services/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    domains: [],
    files: [],
    databases: [],
    emails: [],
    applications: [],
    backups: [],
    stats: {
      diskUsage: 0,
      diskQuota: 0,
      bandwidthUsage: 0,
      bandwidthQuota: 0,
      emailAccounts: 0,
      databases: 0
    },
    loading: false,
    error: null
  }),

  getters: {
    diskUsagePercentage: (state) => {
      return state.stats.diskQuota > 0 
        ? (state.stats.diskUsage / state.stats.diskQuota) * 100 
        : 0
    },
    bandwidthUsagePercentage: (state) => {
      return state.stats.bandwidthQuota > 0 
        ? (state.stats.bandwidthUsage / state.stats.bandwidthQuota) * 100 
        : 0
    }
  },

  actions: {
    async fetchDomains() {
      this.loading = true
      try {
        const response = await userService.getDomains()
        this.domains = response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchDatabases() {
      this.loading = true
      try {
        const response = await userService.getDatabases()
        this.databases = response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchEmails() {
      this.loading = true
      try {
        const response = await userService.getEmails()
        this.emails = response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async createDatabase(dbData) {
      this.loading = true
      try {
        const response = await userService.createDatabase(dbData)
        this.databases.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async createEmail(emailData) {
      this.loading = true
      try {
        const response = await userService.createEmail(emailData)
        this.emails.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchStats() {
      try {
        const response = await userService.getStats()
        this.stats = response.data
      } catch (error) {
        this.error = error.message
        throw error
      }
    }
  }
})
