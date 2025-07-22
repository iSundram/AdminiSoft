
import { createPinia } from 'pinia'
import { useAuthStore } from './auth'
import { useAdminStore } from './admin'
import { useResellerStore } from './reseller'
import { useUserStore } from './user'
import { useSystemStore } from './system'

const pinia = createPinia()

// Export stores for easy access
export {
  useAuthStore,
  useAdminStore,
  useResellerStore,
  useUserStore,
  useSystemStore
}

export default pinia
import { createPinia } from 'pinia'

const pinia = createPinia()

export default pinia
