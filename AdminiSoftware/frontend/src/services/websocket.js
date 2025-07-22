
import { io } from 'socket.io-client'
import { useAuthStore } from '@/store/auth'

class WebSocketService {
  constructor() {
    this.socket = null
    this.connected = false
    this.listeners = new Map()
  }

  connect() {
    const auth = useAuthStore()
    
    if (!auth.token) {
      console.warn('Cannot connect WebSocket: No auth token')
      return
    }

    this.socket = io(process.env.VUE_APP_WS_URL || 'ws://localhost:5000', {
      auth: {
        token: auth.token
      },
      transports: ['websocket']
    })

    this.socket.on('connect', () => {
      this.connected = true
      console.log('WebSocket connected')
    })

    this.socket.on('disconnect', () => {
      this.connected = false
      console.log('WebSocket disconnected')
    })

    this.socket.on('error', (error) => {
      console.error('WebSocket error:', error)
    })

    // Handle real-time notifications
    this.socket.on('notification', (notification) => {
      this.emit('notification', notification)
    })

    // Handle system status updates
    this.socket.on('system_status', (status) => {
      this.emit('system_status', status)
    })

    // Handle resource usage updates
    this.socket.on('resource_usage', (usage) => {
      this.emit('resource_usage', usage)
    })
  }

  disconnect() {
    if (this.socket) {
      this.socket.disconnect()
      this.socket = null
      this.connected = false
    }
  }

  // Subscribe to events
  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set())
    }
    this.listeners.get(event).add(callback)
  }

  // Unsubscribe from events
  off(event, callback) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).delete(callback)
    }
  }

  // Emit to local listeners
  emit(event, data) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).forEach(callback => {
        try {
          callback(data)
        } catch (error) {
          console.error('Error in WebSocket listener:', error)
        }
      })
    }
  }

  // Send message to server
  send(event, data) {
    if (this.socket && this.connected) {
      this.socket.emit(event, data)
    } else {
      console.warn('WebSocket not connected')
    }
  }

  // Request real-time updates for specific resources
  subscribeToUpdates(resources) {
    this.send('subscribe', { resources })
  }

  unsubscribeFromUpdates(resources) {
    this.send('unsubscribe', { resources })
  }
}

export default new WebSocketService()
import { io } from 'socket.io-client'
import { useAuth } from '@/composables/useAuth'

class WebSocketService {
  constructor() {
    this.socket = null
    this.connected = false
    this.listeners = new Map()
  }

  connect() {
    if (this.socket?.connected) return

    const { token } = useAuth()
    
    this.socket = io(import.meta.env.VITE_WS_URL || 'ws://localhost:5000', {
      auth: {
        token: token.value
      },
      transports: ['websocket']
    })

    this.socket.on('connect', () => {
      this.connected = true
      console.log('WebSocket connected')
    })

    this.socket.on('disconnect', () => {
      this.connected = false
      console.log('WebSocket disconnected')
    })

    this.socket.on('error', (error) => {
      console.error('WebSocket error:', error)
    })

    // System events
    this.socket.on('system_alert', (data) => {
      this.emit('system_alert', data)
    })

    this.socket.on('backup_completed', (data) => {
      this.emit('backup_completed', data)
    })

    this.socket.on('service_status_change', (data) => {
      this.emit('service_status_change', data)
    })

    this.socket.on('resource_usage_update', (data) => {
      this.emit('resource_usage_update', data)
    })

    // Account events
    this.socket.on('account_created', (data) => {
      this.emit('account_created', data)
    })

    this.socket.on('account_suspended', (data) => {
      this.emit('account_suspended', data)
    })

    // Email events
    this.socket.on('email_queue_status', (data) => {
      this.emit('email_queue_status', data)
    })

    // Security events
    this.socket.on('security_alert', (data) => {
      this.emit('security_alert', data)
    })

    this.socket.on('login_attempt', (data) => {
      this.emit('login_attempt', data)
    })
  }

  disconnect() {
    if (this.socket) {
      this.socket.disconnect()
      this.socket = null
      this.connected = false
    }
  }

  emit(event, data) {
    const callbacks = this.listeners.get(event) || []
    callbacks.forEach(callback => callback(data))
  }

  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, [])
    }
    this.listeners.get(event).push(callback)

    // Return unsubscribe function
    return () => {
      const callbacks = this.listeners.get(event) || []
      const index = callbacks.indexOf(callback)
      if (index > -1) {
        callbacks.splice(index, 1)
      }
    }
  }

  off(event, callback) {
    const callbacks = this.listeners.get(event) || []
    const index = callbacks.indexOf(callback)
    if (index > -1) {
      callbacks.splice(index, 1)
    }
  }

  // Send messages to server
  send(event, data) {
    if (this.socket?.connected) {
      this.socket.emit(event, data)
    }
  }

  // Join room for real-time updates
  joinRoom(room) {
    this.send('join_room', { room })
  }

  // Leave room
  leaveRoom(room) {
    this.send('leave_room', { room })
  }

  // Get connection status
  isConnected() {
    return this.connected
  }
}

export default new WebSocketService()
