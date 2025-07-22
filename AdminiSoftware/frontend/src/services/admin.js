
import api from './api'

export default {
  // Account Management
  async getAccounts() {
    return api.get('/api/admin/accounts')
  },

  async createAccount(data) {
    return api.post('/api/admin/accounts', data)
  },

  async updateAccount(id, data) {
    return api.put(`/api/admin/accounts/${id}`, data)
  },

  async deleteAccount(id) {
    return api.delete(`/api/admin/accounts/${id}`)
  },

  async suspendAccount(id) {
    return api.post(`/api/admin/accounts/${id}/suspend`)
  },

  async unsuspendAccount(id) {
    return api.post(`/api/admin/accounts/${id}/unsuspend`)
  },

  // Package Management
  async getPackages() {
    return api.get('/api/admin/packages')
  },

  async createPackage(data) {
    return api.post('/api/admin/packages', data)
  },

  async updatePackage(id, data) {
    return api.put(`/api/admin/packages/${id}`, data)
  },

  async deletePackage(id) {
    return api.delete(`/api/admin/packages/${id}`)
  },

  // DNS Management
  async getDNSZones() {
    return api.get('/api/admin/dns/zones')
  },

  async createDNSZone(data) {
    return api.post('/api/admin/dns/zones', data)
  },

  async updateDNSZone(id, data) {
    return api.put(`/api/admin/dns/zones/${id}`, data)
  },

  async deleteDNSZone(id) {
    return api.delete(`/api/admin/dns/zones/${id}`)
  },

  // SSL Management
  async getSSLCertificates() {
    return api.get('/api/admin/ssl/certificates')
  },

  async generateSSL(data) {
    return api.post('/api/admin/ssl/generate', data)
  },

  async installSSL(data) {
    return api.post('/api/admin/ssl/install', data)
  },

  // Backup Management
  async getBackups() {
    return api.get('/api/admin/backups')
  },

  async createBackup(data) {
    return api.post('/api/admin/backups', data)
  },

  async restoreBackup(id) {
    return api.post(`/api/admin/backups/${id}/restore`)
  },

  // System Statistics
  async getSystemStats() {
    return api.get('/api/admin/stats/system')
  },

  async getAccountStats() {
    return api.get('/api/admin/stats/accounts')
  },

  // Security
  async getSecuritySettings() {
    return api.get('/api/admin/security/settings')
  },

  async updateSecuritySettings(data) {
    return api.put('/api/admin/security/settings', data)
  },

  // Server Management
  async getServerInfo() {
    return api.get('/api/admin/server/info')
  },

  async restartService(service) {
    return api.post(`/api/admin/server/restart/${service}`)
  },

  // Email Management
  async getEmailSettings() {
    return api.get('/api/admin/email/settings')
  },

  async updateEmailSettings(data) {
    return api.put('/api/admin/email/settings', data)
  },

  // Monitoring
  async getSystemLoad() {
    return api.get('/api/admin/monitoring/load')
  },

  async getResourceUsage() {
    return api.get('/api/admin/monitoring/resources')
  },

  async getServiceStatus() {
    return api.get('/api/admin/monitoring/services')
  }
}
