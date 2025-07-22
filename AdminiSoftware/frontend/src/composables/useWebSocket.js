
import { ref, onMounted, onUnmounted } from 'vue'
import webSocketService from '@/services/websocket'

export function useWebSocket() {
  const connected = ref(false)
  const error = ref(null)

  function connect() {
    try {
      webSocketService.connect()
      connected.value = webSocketService.connected
    } catch (err) {
      error.value = err.message
      console.error('WebSocket connection error:', err)
    }
  }

  function disconnect() {
    webSocketService.disconnect()
    connected.value = false
  }

  function subscribe(event, callback) {
    webSocketService.on(event, callback)
  }

  function unsubscribe(event, callback) {
    webSocketService.off(event, callback)
  }

  function send(event, data) {
    webSocketService.send(event, data)
  }

  function subscribeToUpdates(resources) {
    webSocketService.subscribeToUpdates(resources)
  }

  function unsubscribeFromUpdates(resources) {
    webSocketService.unsubscribeFromUpdates(resources)
  }

  // Auto-connect on mount and disconnect on unmount
  onMounted(() => {
    connect()
  })

  onUnmounted(() => {
    disconnect()
  })

  return {
    connected,
    error,
    connect,
    disconnect,
    subscribe,
    unsubscribe,
    send,
    subscribeToUpdates,
    unsubscribeFromUpdates
  }
}
