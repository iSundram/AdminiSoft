
import { ref, reactive } from 'vue'
import { useToast } from 'vue-toastification'

export function useNotifications() {
  const toast = useToast()
  const notifications = ref([])

  const state = reactive({
    notifications: [],
    unreadCount: 0
  })

  function addNotification(notification) {
    const newNotification = {
      id: Date.now(),
      timestamp: new Date(),
      read: false,
      type: 'info',
      ...notification
    }

    notifications.value.unshift(newNotification)
    state.notifications.unshift(newNotification)
    updateUnreadCount()

    // Show toast notification
    showToast(newNotification)
  }

  function showToast(notification) {
    const options = {
      timeout: 5000,
      closeOnClick: true,
      pauseOnFocusLoss: true,
      pauseOnHover: true,
      draggable: true,
      draggablePercent: 0.6,
      showCloseButtonOnHover: false,
      hideProgressBar: false,
      closeButton: "button",
      icon: true,
      rtl: false
    }

    switch (notification.type) {
      case 'success':
        toast.success(notification.message, options)
        break
      case 'error':
        toast.error(notification.message, options)
        break
      case 'warning':
        toast.warning(notification.message, options)
        break
      case 'info':
      default:
        toast.info(notification.message, options)
        break
    }
  }

  function markAsRead(id) {
    const notification = notifications.value.find(n => n.id === id)
    if (notification && !notification.read) {
      notification.read = true
      updateUnreadCount()
    }

    const stateNotification = state.notifications.find(n => n.id === id)
    if (stateNotification && !stateNotification.read) {
      stateNotification.read = true
    }
  }

  function markAllAsRead() {
    notifications.value.forEach(n => n.read = true)
    state.notifications.forEach(n => n.read = true)
    updateUnreadCount()
  }

  function removeNotification(id) {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }

    const stateIndex = state.notifications.findIndex(n => n.id === id)
    if (stateIndex > -1) {
      state.notifications.splice(stateIndex, 1)
    }

    updateUnreadCount()
  }

  function clearAll() {
    notifications.value = []
    state.notifications = []
    state.unreadCount = 0
  }

  function updateUnreadCount() {
    state.unreadCount = notifications.value.filter(n => !n.read).length
  }

  // Convenience methods for different notification types
  function success(message, title = 'Success') {
    addNotification({ type: 'success', title, message })
  }

  function error(message, title = 'Error') {
    addNotification({ type: 'error', title, message })
  }

  function warning(message, title = 'Warning') {
    addNotification({ type: 'warning', title, message })
  }

  function info(message, title = 'Information') {
    addNotification({ type: 'info', title, message })
  }

  return {
    notifications,
    state,
    addNotification,
    markAsRead,
    markAllAsRead,
    removeNotification,
    clearAll,
    success,
    error,
    warning,
    info
  }
}
