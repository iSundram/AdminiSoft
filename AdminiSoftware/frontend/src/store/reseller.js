
import { defineStore } from 'pinia'
import resellerService from '@/services/reseller'

export const useResellerStore = defineStore('reseller', {
  state: () => ({
    accounts: [],
    packages: [],
    stats: {
      totalAccounts: 0,
      activeAccounts: 0,
      suspendedAccounts: 0,
      totalBandwidth: 0,
      usedBandwidth: 0
    },
    loading: false,
    error: null
  }),

  getters: {
    activeAccountsCount: (state) => state.accounts.filter(account => account.status === 'active').length,
    suspendedAccountsCount: (state) => state.accounts.filter(account => account.status === 'suspended').length,
    bandwidthUsagePercentage: (state) => {
      return state.stats.totalBandwidth > 0 
        ? (state.stats.usedBandwidth / state.stats.totalBandwidth) * 100 
        : 0
    }
  },

  actions: {
    async fetchAccounts() {
      this.loading = true
      try {
        const response = await resellerService.getAccounts()
        this.accounts = response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchPackages() {
      this.loading = true
      try {
        const response = await resellerService.getPackages()
        this.packages = response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async createAccount(accountData) {
      this.loading = true
      try {
        const response = await resellerService.createAccount(accountData)
        this.accounts.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async createPackage(packageData) {
      this.loading = true
      try {
        const response = await resellerService.createPackage(packageData)
        this.packages.push(response.data)
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
        const response = await resellerService.getStats()
        this.stats = response.data
      } catch (error) {
        this.error = error.message
        throw error
      }
    }
  }
})
