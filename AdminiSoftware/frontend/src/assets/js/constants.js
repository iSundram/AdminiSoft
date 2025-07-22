
// API Configuration
export const API_BASE_URL = process.env.VUE_APP_API_URL || 'http://localhost:5000/api'
export const WS_URL = process.env.VUE_APP_WS_URL || 'ws://localhost:5000'

// Authentication
export const TOKEN_KEY = 'adminisoftware_token'
export const REFRESH_TOKEN_KEY = 'adminisoftware_refresh_token'
export const USER_KEY = 'adminisoftware_user'

// User Roles
export const USER_ROLES = {
  ADMIN: 'admin',
  RESELLER: 'reseller',
  USER: 'user'
}

// Account Status
export const ACCOUNT_STATUS = {
  ACTIVE: 'active',
  SUSPENDED: 'suspended',
  TERMINATED: 'terminated',
  PENDING: 'pending'
}

// File Types
export const FILE_TYPES = {
  IMAGE: ['jpg', 'jpeg', 'png', 'gif', 'svg', 'webp'],
  DOCUMENT: ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx'],
  ARCHIVE: ['zip', 'rar', '7z', 'tar', 'gz'],
  CODE: ['html', 'css', 'js', 'php', 'py', 'java', 'cpp', 'c'],
  TEXT: ['txt', 'md', 'csv', 'json', 'xml', 'yml', 'yaml']
}

// Database Types
export const DATABASE_TYPES = {
  MYSQL: 'mysql',
  POSTGRESQL: 'postgresql',
  MONGODB: 'mongodb',
  REDIS: 'redis'
}

// SSL Certificate Types
export const SSL_TYPES = {
  LETS_ENCRYPT: 'letsencrypt',
  SELF_SIGNED: 'self_signed',
  PURCHASED: 'purchased',
  IMPORTED: 'imported'
}

// Backup Types
export const BACKUP_TYPES = {
  FULL: 'full',
  PARTIAL: 'partial',
  DATABASE: 'database',
  FILES: 'files',
  EMAIL: 'email'
}

// Notification Types
export const NOTIFICATION_TYPES = {
  INFO: 'info',
  SUCCESS: 'success',
  WARNING: 'warning',
  ERROR: 'error'
}

// Service Status
export const SERVICE_STATUS = {
  RUNNING: 'running',
  STOPPED: 'stopped',
  FAILED: 'failed',
  RESTARTING: 'restarting'
}

// System Services
export const SYSTEM_SERVICES = [
  'apache2',
  'nginx',
  'mysql',
  'postgresql',
  'redis',
  'postfix',
  'dovecot',
  'named',
  'proftpd',
  'ssh'
]

// Resource Limits
export const RESOURCE_LIMITS = {
  MIN_DISK_SPACE: 100, // MB
  MIN_BANDWIDTH: 1, // GB
  MIN_EMAIL_ACCOUNTS: 1,
  MIN_DATABASES: 1,
  MIN_SUBDOMAINS: 1,
  MIN_FTP_ACCOUNTS: 1
}

// File Upload Limits
export const UPLOAD_LIMITS = {
  MAX_FILE_SIZE: 100 * 1024 * 1024, // 100MB
  MAX_FILES: 10,
  ALLOWED_TYPES: [
    'image/jpeg',
    'image/png',
    'image/gif',
    'application/pdf',
    'text/plain',
    'application/zip'
  ]
}

// Pagination
export const PAGINATION = {
  DEFAULT_PAGE_SIZE: 20,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100],
  MAX_PAGE_SIZE: 100
}

// Chart Colors
export const CHART_COLORS = {
  PRIMARY: '#3b82f6',
  SUCCESS: '#10b981',
  WARNING: '#f59e0b',
  ERROR: '#ef4444',
  INFO: '#6366f1',
  SECONDARY: '#6b7280'
}

// Time Intervals
export const TIME_INTERVALS = {
  MINUTE: 60 * 1000,
  HOUR: 60 * 60 * 1000,
  DAY: 24 * 60 * 60 * 1000,
  WEEK: 7 * 24 * 60 * 60 * 1000,
  MONTH: 30 * 24 * 60 * 60 * 1000
}

// Regular Expressions
export const REGEX = {
  EMAIL: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  DOMAIN: /^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])*$/,
  IP_ADDRESS: /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
  USERNAME: /^[a-zA-Z0-9_]{3,20}$/,
  PASSWORD: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/,
  SUBDOMAIN: /^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]$/
}

// Application Settings
export const APP_SETTINGS = {
  NAME: 'AdminiSoftware',
  VERSION: '1.0.0',
  COPYRIGHT: '© 2024 AdminiSoftware. All rights reserved.',
  SUPPORT_EMAIL: 'support@adminisoftware.com',
  DOCUMENTATION_URL: 'https://docs.adminisoftware.com',
  GITHUB_URL: 'https://github.com/adminisoftware/adminisoftware'
}

// Feature Flags
export const FEATURES = {
  TWO_FACTOR_AUTH: true,
  BRUTE_FORCE_PROTECTION: true,
  AUTOMATIC_BACKUPS: true,
  SSL_AUTOMATION: true,
  RESOURCE_MONITORING: true,
  EMAIL_NOTIFICATIONS: true,
  API_ACCESS: true,
  CLUSTERING: true,
  CLOUDLINUX_INTEGRATION: true,
  WORDPRESS_TOOLKIT: true
}

// Error Messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Network error. Please check your connection.',
  UNAUTHORIZED: 'You are not authorized to perform this action.',
  FORBIDDEN: 'Access denied.',
  NOT_FOUND: 'Resource not found.',
  SERVER_ERROR: 'Internal server error. Please try again later.',
  VALIDATION_ERROR: 'Please check your input and try again.',
  TIMEOUT_ERROR: 'Request timeout. Please try again.'
}

// Success Messages
export const SUCCESS_MESSAGES = {
  CREATED: 'Created successfully.',
  UPDATED: 'Updated successfully.',
  DELETED: 'Deleted successfully.',
  SAVED: 'Saved successfully.',
  SENT: 'Sent successfully.',
  INSTALLED: 'Installed successfully.',
  BACKUP_CREATED: 'Backup created successfully.',
  RESTORED: 'Restored successfully.'
}

export default {
  API_BASE_URL,
  WS_URL,
  TOKEN_KEY,
  REFRESH_TOKEN_KEY,
  USER_KEY,
  USER_ROLES,
  ACCOUNT_STATUS,
  FILE_TYPES,
  DATABASE_TYPES,
  SSL_TYPES,
  BACKUP_TYPES,
  NOTIFICATION_TYPES,
  SERVICE_STATUS,
  SYSTEM_SERVICES,
  RESOURCE_LIMITS,
  UPLOAD_LIMITS,
  PAGINATION,
  CHART_COLORS,
  TIME_INTERVALS,
  REGEX,
  APP_SETTINGS,
  FEATURES,
  ERROR_MESSAGES,
  SUCCESS_MESSAGES
}
