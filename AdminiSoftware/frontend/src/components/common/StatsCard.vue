
<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center">
      <div class="flex-shrink-0">
        <div
          class="h-12 w-12 rounded-md flex items-center justify-center"
          :class="iconBackgroundClass"
        >
          <component
            :is="iconComponent"
            class="h-6 w-6"
            :class="iconClass"
          />
        </div>
      </div>
      <div class="ml-4 flex-1">
        <div class="flex items-baseline">
          <p class="text-2xl font-semibold text-gray-900">{{ value }}</p>
          <p
            v-if="change !== undefined"
            class="ml-2 flex items-baseline text-sm font-semibold"
            :class="changeClass"
          >
            <component
              :is="changeIcon"
              class="h-4 w-4 flex-shrink-0"
            />
            <span class="sr-only">{{ change > 0 ? 'Increased' : 'Decreased' }} by</span>
            {{ Math.abs(change) }}
          </p>
        </div>
        <p class="text-sm font-medium text-gray-500 truncate">{{ title }}</p>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import {
  UsersIcon,
  CheckCircleIcon,
  CpuChipIcon,
  CircleStackIcon,
  ArrowUpIcon,
  ArrowDownIcon
} from '@heroicons/vue/24/outline'

export default {
  name: 'StatsCard',
  props: {
    title: {
      type: String,
      required: true
    },
    value: {
      type: [String, Number],
      required: true
    },
    icon: {
      type: String,
      default: 'UsersIcon'
    },
    color: {
      type: String,
      default: 'blue',
      validator: (value) => ['blue', 'green', 'yellow', 'purple', 'red'].includes(value)
    },
    change: {
      type: Number,
      default: undefined
    }
  },
  setup(props) {
    const iconComponents = {
      UsersIcon,
      CheckCircleIcon,
      CpuChipIcon,
      CircleStackIcon
    }

    const iconComponent = computed(() => iconComponents[props.icon] || UsersIcon)

    const colorClasses = {
      blue: {
        background: 'bg-blue-50',
        icon: 'text-blue-600'
      },
      green: {
        background: 'bg-green-50',
        icon: 'text-green-600'
      },
      yellow: {
        background: 'bg-yellow-50',
        icon: 'text-yellow-600'
      },
      purple: {
        background: 'bg-purple-50',
        icon: 'text-purple-600'
      },
      red: {
        background: 'bg-red-50',
        icon: 'text-red-600'
      }
    }

    const iconBackgroundClass = computed(() => colorClasses[props.color].background)
    const iconClass = computed(() => colorClasses[props.color].icon)

    const changeIcon = computed(() => {
      if (props.change === undefined) return null
      return props.change > 0 ? ArrowUpIcon : ArrowDownIcon
    })

    const changeClass = computed(() => {
      if (props.change === undefined) return ''
      return props.change > 0 ? 'text-green-600' : 'text-red-600'
    })

    return {
      iconComponent,
      iconBackgroundClass,
      iconClass,
      changeIcon,
      changeClass
    }
  }
}
</script>
</template>
<template>
  <div class="bg-white overflow-hidden shadow rounded-lg">
    <div class="p-5">
      <div class="flex items-center">
        <div class="flex-shrink-0">
          <component :is="icon" class="h-6 w-6 text-gray-400" />
        </div>
        <div class="ml-5 w-0 flex-1">
          <dl>
            <dt class="text-sm font-medium text-gray-500 truncate">{{ title }}</dt>
            <dd>
              <div class="text-lg font-medium text-gray-900">{{ value }}</div>
            </dd>
          </dl>
        </div>
      </div>
    </div>
    <div v-if="change" class="bg-gray-50 px-5 py-3">
      <div class="text-sm">
        <span :class="[changeType === 'increase' ? 'text-green-600' : 'text-red-600', 'font-medium']">
          {{ changeType === 'increase' ? '+' : '-' }}{{ Math.abs(change) }}%
        </span>
        <span class="text-gray-500"> from last month</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'StatsCard',
  props: {
    title: {
      type: String,
      required: true
    },
    value: {
      type: [String, Number],
      required: true
    },
    icon: {
      type: Object,
      required: true
    },
    change: {
      type: Number,
      default: null
    },
    changeType: {
      type: String,
      default: 'increase',
      validator: value => ['increase', 'decrease'].includes(value)
    }
  }
}
</script>
<template>
  <div class="bg-white overflow-hidden shadow rounded-lg">
    <div class="p-5">
      <div class="flex items-center">
        <div class="flex-shrink-0">
          <component
            :is="icon"
            class="h-6 w-6"
            :class="iconColor"
          />
        </div>
        <div class="ml-5 w-0 flex-1">
          <dl>
            <dt class="text-sm font-medium text-gray-500 truncate">
              {{ title }}
            </dt>
            <dd>
              <div class="text-lg font-medium text-gray-900">
                {{ formattedValue }}
              </div>
            </dd>
          </dl>
        </div>
      </div>
    </div>
    <div v-if="change !== undefined" class="bg-gray-50 px-5 py-3">
      <div class="text-sm">
        <span
          class="font-medium"
          :class="changeColor"
        >
          {{ changeText }}
        </span>
        <span class="text-gray-500 ml-1">from last period</span>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'StatsCard',
  props: {
    title: {
      type: String,
      required: true
    },
    value: {
      type: [String, Number],
      required: true
    },
    icon: {
      type: Object,
      required: true
    },
    change: {
      type: Number,
      default: undefined
    },
    format: {
      type: String,
      default: 'number', // 'number', 'currency', 'percentage', 'bytes'
    },
    color: {
      type: String,
      default: 'blue', // 'blue', 'green', 'yellow', 'red', 'gray'
    }
  },
  setup(props) {
    const iconColor = computed(() => {
      const colors = {
        blue: 'text-blue-400',
        green: 'text-green-400',
        yellow: 'text-yellow-400',
        red: 'text-red-400',
        gray: 'text-gray-400'
      }
      return colors[props.color] || colors.blue
    })

    const formattedValue = computed(() => {
      switch (props.format) {
        case 'currency':
          return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD'
          }).format(props.value)
        case 'percentage':
          return `${props.value}%`
        case 'bytes':
          return formatBytes(props.value)
        case 'number':
        default:
          return new Intl.NumberFormat('en-US').format(props.value)
      }
    })

    const changeColor = computed(() => {
      if (props.change === undefined) return ''
      return props.change >= 0 ? 'text-green-600' : 'text-red-600'
    })

    const changeText = computed(() => {
      if (props.change === undefined) return ''
      const sign = props.change >= 0 ? '+' : ''
      return `${sign}${props.change}%`
    })

    const formatBytes = (bytes) => {
      if (bytes === 0) return '0 Bytes'
      const k = 1024
      const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    }

    return {
      iconColor,
      formattedValue,
      changeColor,
      changeText
    }
  }
}
</script>
