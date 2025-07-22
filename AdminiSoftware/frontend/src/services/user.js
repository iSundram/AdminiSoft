
import api from './api'

export default {
  // Domain management
  getDomains() {
    return api.get('/user/domains')
  },

  createSubdomain(subdomainData) {
    return api.post('/user/domains/subdomains', subdomainData)
  },

  createDomainPointer(pointerData) {
    return api.post('/user/domains/pointers', pointerData)
  },

  updateDNSRecord(recordData) {
    return api.put('/user/domains/dns', recordData)
  },

  // File management
  getFiles(path = '/') {
    return api.get(`/user/files?path=${encodeURIComponent(path)}`)
  },

  uploadFile(fileData) {
    return api.post('/user/files/upload', fileData)
  },

  deleteFile(filePath) {
    return api.delete(`/user/files?path=${encodeURIComponent(filePath)}`)
  },

  createDirectory(dirPath) {
    return api.post('/user/files/directory', { path: dirPath })
  },

  // FTP accounts
  getFTPAccounts() {
    return api.get('/user/ftp')
  },

  createFTPAccount(ftpData) {
    return api.post('/user/ftp', ftpData)
  },

  // Database management
  getDatabases() {
    return api.get('/user/databases')
  },

  createDatabase(dbData) {
    return api.post('/user/databases', dbData)
  },

  deleteDatabase(dbName) {
    return api.delete(`/user/databases/${dbName}`)
  },

  createDatabaseUser(userData) {
    return api.post('/user/databases/users', userData)
  },

  // Email management
  getEmails() {
    return api.get('/user/emails')
  },

  createEmail(emailData) {
    return api.post('/user/emails', emailData)
  },

  deleteEmail(email) {
    return api.delete(`/user/emails/${email}`)
  },

  createForwarder(forwarderData) {
    return api.post('/user/emails/forwarders', forwarderData)
  },

  createAutoresponder(autoresponderData) {
    return api.post('/user/emails/autoresponders', autoresponderData)
  },

  // SSL management
  getSSLCertificates() {
    return api.get('/user/ssl')
  },

  generateSSLCertificate(domain) {
    return api.post('/user/ssl/generate', { domain })
  },

  installSSLCertificate(sslData) {
    return api.post('/user/ssl/install', sslData)
  },

  // Applications
  getApplications() {
    return api.get('/user/apps')
  },

  installApplication(appData) {
    return api.post('/user/apps/install', appData)
  },

  // WordPress management
  getWordPressSites() {
    return api.get('/user/wordpress')
  },

  installWordPress(wpData) {
    return api.post('/user/wordpress/install', wpData)
  },

  // Statistics
  getStats() {
    return api.get('/user/stats')
  },

  getBandwidthStats() {
    return api.get('/user/stats/bandwidth')
  },

  getDiskUsageStats() {
    return api.get('/user/stats/disk-usage')
  },

  // Backup
  createBackup(backupData) {
    return api.post('/user/backup', backupData)
  },

  getBackups() {
    return api.get('/user/backup')
  },

  restoreBackup(backupId) {
    return api.post(`/user/backup/${backupId}/restore`)
  }
}
import api from './api'

export default {
  // Domain Management
  async getDomains() {
    return api.get('/api/user/domains')
  },

  async createDomain(data) {
    return api.post('/api/user/domains', data)
  },

  async deleteDomain(id) {
    return api.delete(`/api/user/domains/${id}`)
  },

  // File Management
  async getFiles(path = '/') {
    return api.get(`/api/user/files?path=${encodeURIComponent(path)}`)
  },

  async createFolder(data) {
    return api.post('/api/user/files/folder', data)
  },

  async uploadFile(formData) {
    return api.post('/api/user/files/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  async deleteFile(path) {
    return api.delete('/api/user/files', { data: { path } })
  },

  async downloadFile(path) {
    return api.get(`/api/user/files/download?path=${encodeURIComponent(path)}`, {
      responseType: 'blob'
    })
  },

  // Email Management
  async getEmailAccounts() {
    return api.get('/api/user/email/accounts')
  },

  async createEmailAccount(data) {
    return api.post('/api/user/email/accounts', data)
  },

  async deleteEmailAccount(id) {
    return api.delete(`/api/user/email/accounts/${id}`)
  },

  async getForwarders() {
    return api.get('/api/user/email/forwarders')
  },

  async createForwarder(data) {
    return api.post('/api/user/email/forwarders', data)
  },

  // Database Management
  async getDatabases() {
    return api.get('/api/user/databases')
  },

  async createDatabase(data) {
    return api.post('/api/user/databases', data)
  },

  async deleteDatabase(id) {
    return api.delete(`/api/user/databases/${id}`)
  },

  async getDatabaseUsers() {
    return api.get('/api/user/databases/users')
  },

  async createDatabaseUser(data) {
    return api.post('/api/user/databases/users', data)
  },

  // SSL Management
  async getSSLCertificates() {
    return api.get('/api/user/ssl/certificates')
  },

  async requestSSL(data) {
    return api.post('/api/user/ssl/request', data)
  },

  async installSSL(data) {
    return api.post('/api/user/ssl/install', data)
  },

  // Application Management
  async getApplications() {
    return api.get('/api/user/apps')
  },

  async installApplication(data) {
    return api.post('/api/user/apps/install', data)
  },

  async uninstallApplication(id) {
    return api.delete(`/api/user/apps/${id}`)
  },

  // WordPress Management
  async getWordPressInstalls() {
    return api.get('/api/user/wordpress')
  },

  async installWordPress(data) {
    return api.post('/api/user/wordpress/install', data)
  },

  async updateWordPress(id) {
    return api.post(`/api/user/wordpress/${id}/update`)
  },

  // Statistics
  async getBandwidthUsage() {
    return api.get('/api/user/stats/bandwidth')
  },

  async getDiskUsage() {
    return api.get('/api/user/stats/disk')
  },

  async getVisitorStats() {
    return api.get('/api/user/stats/visitors')
  },

  // Backup Management
  async getBackups() {
    return api.get('/api/user/backups')
  },

  async createBackup() {
    return api.post('/api/user/backups')
  },

  async downloadBackup(id) {
    return api.get(`/api/user/backups/${id}/download`, {
      responseType: 'blob'
    })
  },

  async restoreBackup(id) {
    return api.post(`/api/user/backups/${id}/restore`)
  },

  // Account Settings
  async getAccountInfo() {
    return api.get('/api/user/account')
  },

  async updateAccountInfo(data) {
    return api.put('/api/user/account', data)
  },

  async changePassword(data) {
    return api.post('/api/user/account/password', data)
  },

  // Two-Factor Authentication
  async enable2FA() {
    return api.post('/api/user/account/2fa/enable')
  },

  async disable2FA(data) {
    return api.post('/api/user/account/2fa/disable', data)
  },

  async verify2FA(data) {
    return api.post('/api/user/account/2fa/verify', data)
  }
}
</script>
