
<template>
  <div class="relative">
    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
      <MagnifyingGlassIcon class="h-5 w-5 text-gray-400" />
    </div>
    <input
      v-model="searchQuery"
      type="text"
      :placeholder="placeholder"
      class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
      @input="handleInput"
    />
  </div>
</template>

<script>
import { MagnifyingGlassIcon } from '@heroicons/vue/24/outline'

export default {
  name: 'SearchBar',
  components: {
    MagnifyingGlassIcon
  },
  props: {
    placeholder: {
      type: String,
      default: 'Search...'
    },
    modelValue: {
      type: String,
      default: ''
    }
  },
  emits: ['update:modelValue', 'search'],
  data() {
    return {
      searchQuery: this.modelValue
    }
  },
  watch: {
    modelValue(newValue) {
      this.searchQuery = newValue
    }
  },
  methods: {
    handleInput() {
      this.$emit('update:modelValue', this.searchQuery)
      this.$emit('search', this.searchQuery)
    }
  }
}
</script>
<template>
  <div class="relative">
    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
      <MagnifyingGlassIcon class="h-5 w-5 text-gray-400" />
    </div>
    <input
      v-model="searchQuery"
      type="text"
      :placeholder="placeholder"
      class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
      @input="handleSearch"
      @keyup.enter="$emit('search', searchQuery)"
    />
    <div v-if="searchQuery" class="absolute inset-y-0 right-0 pr-3 flex items-center">
      <button
        @click="clearSearch"
        class="text-gray-400 hover:text-gray-600"
      >
        <XMarkIcon class="h-5 w-5" />
      </button>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'
import { MagnifyingGlassIcon, XMarkIcon } from '@heroicons/vue/24/outline'

export default {
  name: 'SearchBar',
  components: {
    MagnifyingGlassIcon,
    XMarkIcon
  },
  props: {
    placeholder: {
      type: String,
      default: 'Search...'
    },
    modelValue: {
      type: String,
      default: ''
    },
    debounce: {
      type: Number,
      default: 300
    }
  },
  emits: ['update:modelValue', 'search', 'clear'],
  setup(props, { emit }) {
    const searchQuery = ref(props.modelValue)
    let debounceTimer = null

    const handleSearch = () => {
      emit('update:modelValue', searchQuery.value)
      
      if (debounceTimer) {
        clearTimeout(debounceTimer)
      }
      
      debounceTimer = setTimeout(() => {
        emit('search', searchQuery.value)
      }, props.debounce)
    }

    const clearSearch = () => {
      searchQuery.value = ''
      emit('update:modelValue', '')
      emit('search', '')
      emit('clear')
    }

    watch(() => props.modelValue, (newValue) => {
      searchQuery.value = newValue
    })

    return {
      searchQuery,
      handleSearch,
      clearSearch
    }
  }
}
</script>
