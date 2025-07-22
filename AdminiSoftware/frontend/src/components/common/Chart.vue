
<template>
  <div class="bg-white p-6 rounded-lg shadow">
    <h3 class="text-lg font-medium text-gray-900 mb-4">{{ title }}</h3>
    <canvas :id="chartId" :width="width" :height="height"></canvas>
  </div>
</template>

<script>
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, BarElement } from 'chart.js'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, BarElement)

export default {
  name: 'Chart',
  props: {
    title: {
      type: String,
      required: true
    },
    type: {
      type: String,
      default: 'line',
      validator: value => ['line', 'bar', 'doughnut', 'pie'].includes(value)
    },
    data: {
      type: Object,
      required: true
    },
    options: {
      type: Object,
      default: () => ({})
    },
    width: {
      type: Number,
      default: 400
    },
    height: {
      type: Number,
      default: 200
    }
  },
  data() {
    return {
      chart: null,
      chartId: `chart-${Math.random().toString(36).substr(2, 9)}`
    }
  },
  mounted() {
    this.initChart()
  },
  beforeUnmount() {
    if (this.chart) {
      this.chart.destroy()
    }
  },
  watch: {
    data: {
      handler() {
        this.updateChart()
      },
      deep: true
    }
  },
  methods: {
    initChart() {
      const ctx = document.getElementById(this.chartId).getContext('2d')
      this.chart = new ChartJS(ctx, {
        type: this.type,
        data: this.data,
        options: {
          responsive: true,
          maintainAspectRatio: false,
          ...this.options
        }
      })
    },
    updateChart() {
      if (this.chart) {
        this.chart.data = this.data
        this.chart.update()
      }
    }
  }
}
</script>
<template>
  <div class="chart-container">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement
} from 'chart.js'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement
)

export default {
  name: 'Chart',
  props: {
    type: {
      type: String,
      default: 'line',
      validator: (value) => ['line', 'bar', 'doughnut', 'pie'].includes(value)
    },
    data: {
      type: Object,
      required: true
    },
    options: {
      type: Object,
      default: () => ({})
    },
    height: {
      type: Number,
      default: 400
    }
  },
  setup(props) {
    const chartCanvas = ref(null)
    let chartInstance = null

    const defaultOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'top'
        },
        title: {
          display: false
        }
      },
      scales: props.type === 'doughnut' || props.type === 'pie' ? undefined : {
        y: {
          beginAtZero: true
        }
      }
    }

    const createChart = () => {
      if (chartInstance) {
        chartInstance.destroy()
      }

      const ctx = chartCanvas.value?.getContext('2d')
      if (!ctx) return

      chartInstance = new ChartJS(ctx, {
        type: props.type,
        data: props.data,
        options: {
          ...defaultOptions,
          ...props.options
        }
      })
    }

    const updateChart = () => {
      if (chartInstance) {
        chartInstance.data = props.data
        chartInstance.options = {
          ...defaultOptions,
          ...props.options
        }
        chartInstance.update()
      }
    }

    watch(() => props.data, updateChart, { deep: true })
    watch(() => props.options, updateChart, { deep: true })

    onMounted(() => {
      createChart()
    })

    onBeforeUnmount(() => {
      if (chartInstance) {
        chartInstance.destroy()
      }
    })

    return {
      chartCanvas
    }
  }
}
</script>

<style scoped>
.chart-container {
  position: relative;
  height: 400px;
  width: 100%;
}
</style>
