
import api from './api'

export default {
  // Account management
  getAccounts() {
    return api.get('/reseller/accounts')
  },

  createAccount(accountData) {
    return api.post('/reseller/accounts', accountData)
  },

  updateAccount(accountId, accountData) {
    return api.put(`/reseller/accounts/${accountId}`, accountData)
  },

  deleteAccount(accountId) {
    return api.delete(`/reseller/accounts/${accountId}`)
  },

  suspendAccount(accountId) {
    return api.post(`/reseller/accounts/${accountId}/suspend`)
  },

  unsuspendAccount(accountId) {
    return api.post(`/reseller/accounts/${accountId}/unsuspend`)
  },

  // Package management
  getPackages() {
    return api.get('/reseller/packages')
  },

  createPackage(packageData) {
    return api.post('/reseller/packages', packageData)
  },

  updatePackage(packageId, packageData) {
    return api.put(`/reseller/packages/${packageId}`, packageData)
  },

  deletePackage(packageId) {
    return api.delete(`/reseller/packages/${packageId}`)
  },

  // Statistics
  getStats() {
    return api.get('/reseller/stats')
  },

  getAccountStats(accountId) {
    return api.get(`/reseller/accounts/${accountId}/stats`)
  },

  // Branding
  getBranding() {
    return api.get('/reseller/branding')
  },

  updateBranding(brandingData) {
    return api.put('/reseller/branding', brandingData)
  },

  // Communication
  messageAllUsers(messageData) {
    return api.post('/reseller/message-all', messageData)
  },

  createAnnouncement(announcementData) {
    return api.post('/reseller/announcements', announcementData)
  }
}
