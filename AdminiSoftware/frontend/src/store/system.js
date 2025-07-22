
import { defineStore } from 'pinia'

export const useSystemStore = defineStore('system', {
  state: () => ({
    settings: {
      siteName: 'AdminiSoftware',
      theme: 'default',
      maintenanceMode: false,
      registrationEnabled: false
    },
    serverInfo: {
      version: '1.0.0',
      uptime: 0,
      loadAverage: 0,
      memoryUsage: 0,
      diskUsage: 0
    },
    notifications: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchSettings() {
      this.loading = true
      try {
        // Mock data for now
        this.settings = {
          siteName: 'AdminiSoftware',
          theme: 'default',
          maintenanceMode: false,
          registrationEnabled: false
        }
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateSettings(newSettings) {
      this.loading = true
      try {
        this.settings = { ...this.settings, ...newSettings }
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchServerInfo() {
      try {
        // Mock data for now
        this.serverInfo = {
          version: '1.0.0',
          uptime: Math.floor(Math.random() * 1000000),
          loadAverage: Math.random() * 2,
          memoryUsage: Math.random() * 100,
          diskUsage: Math.random() * 100
        }
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    addNotification(notification) {
      this.notifications.unshift({
        id: Date.now(),
        timestamp: new Date(),
        read: false,
        ...notification
      })
    },

    markNotificationAsRead(id) {
      const notification = this.notifications.find(n => n.id === id)
      if (notification) {
        notification.read = true
      }
    }
  }
})
