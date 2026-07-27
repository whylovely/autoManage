<script setup lang="ts">
withDefaults(defineProps<{
  title: string
  text?: string
}>(), {
  text: '',
})
</script>

<template>
  <div class="rounded-2xl border border-dashed border-slate-300 bg-slate-50 px-6 py-12 text-center">
    <div class="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-white text-xl shadow-sm">◇</div>
    <h3 class="font-bold text-slate-800">{{ title }}</h3>
    <p v-if="text" class="mx-auto mt-1 max-w-md text-sm text-slate-500">{{ text }}</p>
  </div>
</template>
