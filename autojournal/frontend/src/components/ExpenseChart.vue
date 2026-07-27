<script setup lang="ts">
import { Chart, registerables, type ChartType } from 'chart.js'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

Chart.register(...registerables)

const props = withDefaults(defineProps<{
  labels: string[]
  values: number[]
  type?: 'bar' | 'doughnut'
}>(), {
  type: 'bar',
})

const canvas = ref<HTMLCanvasElement | null>(null)
let chart: Chart | null = null

function renderChart() {
  if (!canvas.value) return
  chart?.destroy()

  const colors = ['#169b78', '#4f7cff', '#f59e0b', '#ef6a83', '#8b5cf6', '#06b6d4']
  chart = new Chart(canvas.value, {
    type: props.type as ChartType,
    data: {
      labels: props.labels,
      datasets: [{
        data: props.values,
        backgroundColor: props.type === 'doughnut' ? colors : '#169b78',
        borderRadius: props.type === 'bar' ? 8 : 0,
        borderWidth: 0,
      }],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: props.type === 'doughnut', position: 'bottom' },
      },
      scales: props.type === 'bar'
        ? {
            x: { grid: { display: false } },
            y: { beginAtZero: true, grid: { color: '#edf0f3' } },
          }
        : undefined,
    },
  })
}

onMounted(renderChart)
watch(() => [props.labels, props.values, props.type], renderChart, { deep: true })
onBeforeUnmount(() => chart?.destroy())
</script>

<template>
  <div class="h-64">
    <canvas ref="canvas" />
  </div>
</template>
