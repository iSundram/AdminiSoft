
import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

// Import views
import Login from '@/views/Login.vue'

// AdminiCore (Admin) views
const AdminDashboard = () => import('@/views/AdminiCore/Dashboard.vue')
const ListAccounts = () => import('@/views/AdminiCore/AccountManagement/ListAccounts.vue')
const CreateAccount = () => import('@/views/AdminiCore/AccountManagement/CreateAccount.vue')
const SecurityOverview = () => import('@/views/AdminiCore/SecurityCenter/Overview.vue')
const DNSZoneManager = () => import('@/views/AdminiCore/DNSManagement/DNSZoneManager.vue')
const BackupConfiguration = () => import('@/views/AdminiCore/BackupManagement/BackupConfiguration.vue')

// AdminiReseller views
const ResellerDashboard = () => import('@/views/AdminiReseller/Dashboard.vue')
const ResellerAccounts = () => import('@/views/AdminiReseller/AccountManagement/ListAccounts.vue')
const ResellerPackages = () => import('@/views/AdminiReseller/PackageManagement/ListPackages.vue')

// AdminiPanel (User) views
const UserDashboard = () => import('@/views/AdminiPanel/Dashboard.vue')
const FileManager = () => import('@/views/AdminiPanel/FileManagement/FileManager.vue')
const EmailAccounts = () => import('@/views/AdminiPanel/EmailManagement/EmailAccounts.vue')
const MySQLDatabases = () => import('@/views/AdminiPanel/DatabaseManagement/MySQLDatabases.vue')
const WordPressManager = () => import('@/views/AdminiPanel/ApplicationManagement/WordPressManager.vue')

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresGuest: true }
  },
  
  // AdminiCore Routes (Admin)
  {
    path: '/admin',
    name: 'Admin',
    redirect: '/admin/dashboard',
    meta: { requiresAuth: true, role: 'admin' },
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: AdminDashboard
      },
      {
        path: 'accounts',
        name: 'AdminAccounts',
        component: ListAccounts
      },
      {
        path: 'accounts/create',
        name: 'AdminCreateAccount',
        component: CreateAccount
      },
      {
        path: 'security',
        name: 'AdminSecurity',
        component: SecurityOverview
      },
      {
        path: 'dns',
        name: 'AdminDNS',
        component: DNSZoneManager
      },
      {
        path: 'backup',
        name: 'AdminBackup',
        component: BackupConfiguration
      }
    ]
  },
  
  // AdminiReseller Routes
  {
    path: '/reseller',
    name: 'Reseller',
    redirect: '/reseller/dashboard',
    meta: { requiresAuth: true, role: ['admin', 'reseller'] },
    children: [
      {
        path: 'dashboard',
        name: 'ResellerDashboard',
        component: ResellerDashboard
      },
      {
        path: 'accounts',
        name: 'ResellerAccounts',
        component: ResellerAccounts
      },
      {
        path: 'packages',
        name: 'ResellerPackages',
        component: ResellerPackages
      }
    ]
  },
  
  // AdminiPanel Routes (User)
  {
    path: '/',
    name: 'User',
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'UserDashboard',
        component: UserDashboard
      },
      {
        path: 'files',
        name: 'FileManager',
        component: FileManager
      },
      {
        path: 'email',
        name: 'EmailAccounts',
        component: EmailAccounts
      },
      {
        path: 'databases',
        name: 'MySQLDatabases',
        component: MySQLDatabases
      },
      {
        path: 'wordpress',
        name: 'WordPressManager',
        component: WordPressManager
      }
    ]
  },
  
  // Catch all route
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  const { isAuthenticated, user, checkAuth } = useAuth()
  
  // Check authentication status
  await checkAuth()
  
  if (to.meta.requiresAuth && !isAuthenticated.value) {
    next('/login')
    return
  }
  
  if (to.meta.requiresGuest && isAuthenticated.value) {
    next('/')
    return
  }
  
  if (to.meta.role) {
    const requiredRoles = Array.isArray(to.meta.role) ? to.meta.role : [to.meta.role]
    if (!requiredRoles.includes(user.value?.role)) {
      next('/')
      return
    }
  }
  
  next()
})

export default router
