
<template>
  <teleport to="body">
    <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto">
      <div class="flex min-h-screen items-center justify-center px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <!-- Background overlay -->
        <transition
          enter-active-class="ease-out duration-300"
          enter-from-class="opacity-0"
          enter-to-class="opacity-100"
          leave-active-class="ease-in duration-200"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="isOpen" class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="close"></div>
        </transition>

        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

        <!-- Modal panel -->
        <transition
          enter-active-class="ease-out duration-300"
          enter-from-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          enter-to-class="opacity-100 translate-y-0 sm:scale-100"
          leave-active-class="ease-in duration-200"
          leave-from-class="opacity-100 translate-y-0 sm:scale-100"
          leave-to-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
        >
          <div
            v-if="isOpen"
            :class="[
              'inline-block align-bottom bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:p-6',
              sizeClass
            ]"
          >
            <!-- Header -->
            <div v-if="title || $slots.header" class="mb-4">
              <div class="flex items-center justify-between">
                <slot name="header">
                  <h3 class="text-lg font-medium text-gray-900">{{ title }}</h3>
                </slot>
                <button
                  v-if="showCloseButton"
                  @click="close"
                  class="text-gray-400 hover:text-gray-500 focus:outline-none focus:text-gray-500 transition ease-in-out duration-150"
                >
                  <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>

            <!-- Content -->
            <div class="mb-4">
              <slot></slot>
            </div>

            <!-- Footer -->
            <div v-if="$slots.footer" class="flex justify-end space-x-2">
              <slot name="footer"></slot>
            </div>
          </div>
        </transition>
      </div>
    </div>
  </teleport>
</template>

<script>
export default {
  name: 'Modal',
  props: {
    isOpen: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: ''
    },
    size: {
      type: String,
      default: 'md',
      validator: value => ['sm', 'md', 'lg', 'xl', 'full'].includes(value)
    },
    showCloseButton: {
      type: Boolean,
      default: true
    }
  },
  emits: ['close'],
  computed: {
    sizeClass() {
      const sizes = {
        sm: 'sm:max-w-sm sm:w-full',
        md: 'sm:max-w-lg sm:w-full',
        lg: 'sm:max-w-2xl sm:w-full',
        xl: 'sm:max-w-4xl sm:w-full',
        full: 'sm:max-w-6xl sm:w-full'
      }
      return sizes[this.size]
    }
  },
  methods: {
    close() {
      this.$emit('close')
    }
  },
  watch: {
    isOpen(newVal) {
      if (newVal) {
        document.body.style.overflow = 'hidden'
      } else {
        document.body.style.overflow = 'auto'
      }
    }
  },
  beforeUnmount() {
    document.body.style.overflow = 'auto'
  }
}
</script>
