
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Toast from 'vue-toastification'
import LoadingPlugin from 'vue3-loading-overlay'

import App from './App.vue'
import router from './router'

// Import CSS
import './assets/css/main.css'
import 'vue-toastification/dist/index.css'
import 'vue3-loading-overlay/dist/vue3-loading-overlay.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(Toast, {
  transition: "Vue-Toastification__bounce",
  maxToasts: 20,
  newestOnTop: true
})
app.use(LoadingPlugin)

app.mount('#app')
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// Styles
import './assets/css/main.css'

// Create Vue app
const app = createApp(App)

// Use plugins
app.use(createPinia())
app.use(router)

// Global error handler
app.config.errorHandler = (error, instance, info) => {
  console.error('Global error:', error)
  console.error('Component:', instance)
  console.error('Info:', info)
}

// Mount app
app.mount('#app')
