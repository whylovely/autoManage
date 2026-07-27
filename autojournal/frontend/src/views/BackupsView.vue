<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { CreateBackup, ExportVehicleExpenses, ListBackups } from '../../wailsjs/go/handler/App'
import EmptyState from '../components/EmptyState.vue'
import { formatDateTime } from '../lib/format'
import type { BackupRecord, VehicleRecord } from '../types'

const props = defineProps<{ vehicle: VehicleRecord | null }>()
const emit = defineEmits<{ (event: 'notify', message: string, tone?: 'success' | 'error' | 'info'): void }>()
const backups = ref<BackupRecord[]>([])
const note = ref('')
const creating = ref(false)

async function load() {
  try {
    backups.value = (await ListBackups() ?? []) as BackupRecord[]
  } catch (error) {
    emit('notify', String(error), 'error')
  }
}

async function create() {
  creating.value = true
  try {
    await CreateBackup(note.value)
    note.value = ''
    await load()
    emit('notify', 'Резервная копия создана', 'success')
  } catch (error) {
    emit('notify', String(error), 'error')
  } finally {
    creating.value = false
  }
}

async function exportExpenses(format: 'csv' | 'json') {
  if (!props.vehicle) return
  try {
    const destination = await ExportVehicleExpenses(props.vehicle.id, format)
    if (destination) emit('notify', `Расходы экспортированы: ${destination}`, 'success')
  } catch (error) {
    emit('notify', String(error), 'error')
  }
}

onMounted(load)
</script>

<template>
  <section class="space-y-6">
    <header>
      <p class="text-sm font-bold uppercase tracking-[0.16em] text-brand-600">Безопасность данных</p>
      <h1 class="mt-1 text-3xl font-black tracking-tight">Резервные копии</h1>
      <p class="mt-2 text-sm text-slate-500">SQLite-снимок создаётся без остановки приложения.</p>
    </header>

    <form class="panel flex flex-col gap-4 p-5 sm:flex-row sm:items-end" @submit.prevent="create">
      <label class="flex-1"><span class="label">Комментарий к копии</span><input v-model="note" class="field" maxlength="200" placeholder="Например: перед большим ТО"></label>
      <button class="btn-primary whitespace-nowrap" :disabled="creating">{{ creating ? 'Создаём…' : 'Создать резервную копию' }}</button>
    </form>

    <article class="panel flex flex-col gap-4 p-5 sm:flex-row sm:items-center sm:justify-between">
      <div><h2 class="font-black">Экспорт расходов</h2><p class="mt-1 text-sm text-slate-500">Сохраните журнал выбранного автомобиля в переносимом формате.</p></div>
      <div class="flex gap-2"><button class="btn-secondary" :disabled="!vehicle" @click="exportExpenses('csv')">Экспорт CSV</button><button class="btn-secondary" :disabled="!vehicle" @click="exportExpenses('json')">Экспорт JSON</button></div>
    </article>

    <article class="panel overflow-hidden">
      <div class="border-b border-slate-100 px-5 py-4"><h2 class="text-lg font-black">История копий</h2></div>
      <ul v-if="backups.length" class="divide-y divide-slate-100">
        <li v-for="backup in backups" :key="backup.id" class="grid gap-2 px-5 py-4 md:grid-cols-[180px_1fr]">
          <time class="text-sm font-bold text-slate-700">{{ formatDateTime(backup.createdAt) }}</time>
          <div class="min-w-0">
            <p class="truncate text-sm text-slate-500" :title="backup.filePath">{{ backup.filePath }}</p>
            <p v-if="backup.note" class="mt-1 text-sm font-bold">{{ backup.note }}</p>
          </div>
        </li>
      </ul>
      <div v-else class="p-5"><EmptyState title="Копий ещё нет" text="Создайте первую копию базы данных кнопкой выше." /></div>
    </article>
  </section>
</template>
