<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { GetDueReminders, GetVehicles } from '../wailsjs/go/handler/App'
import { EventsOn } from '../wailsjs/runtime/runtime'
import BackupsView from './views/BackupsView.vue'
import DashboardView from './views/DashboardView.vue'
import ExpensesView from './views/ExpensesView.vue'
import RemindersView from './views/RemindersView.vue'
import VehicleProfileView from './views/VehicleProfileView.vue'
import type { DueReminderRecord, Notice, PageName, VehicleRecord } from './types'

const navItems: Array<{ id: PageName; label: string; icon: string }> = [
  { id: 'dashboard', label: 'Обзор', icon: '⌂' },
  { id: 'expenses', label: 'Расходы', icon: '₽' },
  { id: 'reminders', label: 'Напоминания', icon: '◷' },
  { id: 'vehicle', label: 'Автомобиль', icon: '◇' },
  { id: 'backups', label: 'Копии', icon: '▣' },
]

const activePage = ref<PageName>('dashboard')
const vehicles = ref<VehicleRecord[]>([])
const selectedVehicleId = ref<number | null>(null)
const notices = ref<Notice[]>([])
let noticeID = 0

const selectedVehicle = computed(() => vehicles.value.find(item => item.id === selectedVehicleId.value) ?? null)
const pageTitle = computed(() => navItems.find(item => item.id === activePage.value)?.label ?? '')

function notify(message: string, tone: Notice['tone'] = 'info') {
  const id = ++noticeID
  notices.value.push({ id, message: message.replace(/^Error:\s*/, ''), tone })
  window.setTimeout(() => dismissNotice(id), 5000)
}

function dismissNotice(id: number) {
  notices.value = notices.value.filter(item => item.id !== id)
}

function showDueReminder(payload: DueReminderRecord) {
  const reasons = [
    payload.dueByDate ? 'по дате' : '',
    payload.dueByOdometer ? 'по пробегу' : '',
  ].filter(Boolean).join(' и ')
  notify(`«${payload.reminder.title}»: приближается срок ${reasons}`, 'info')
}

const unsubscribe = EventsOn('reminder:due', (payload: DueReminderRecord) => showDueReminder(payload))

async function loadVehicles(preferredId?: number) {
  try {
    vehicles.value = (await GetVehicles() ?? []) as VehicleRecord[]
    const requested = preferredId ?? selectedVehicleId.value
    selectedVehicleId.value = vehicles.value.some(item => item.id === requested)
      ? requested
      : vehicles.value[0]?.id ?? null
  } catch (error) {
    notify(String(error), 'error')
  }
}

async function loadInitialDueReminders() {
  try {
    const due = (await GetDueReminders() ?? []) as DueReminderRecord[]
    due.slice(0, 3).forEach(showDueReminder)
  } catch (error) {
    notify(String(error), 'error')
  }
}

onMounted(async () => {
  await loadVehicles()
  await loadInitialDueReminders()
})

onBeforeUnmount(unsubscribe)
</script>

<template>
  <div class="min-h-screen bg-canvas lg:grid lg:grid-cols-[250px_1fr]">
    <aside class="border-b border-slate-200 bg-ink px-4 py-4 text-white lg:fixed lg:inset-y-0 lg:w-[250px] lg:border-b-0">
      <div class="flex items-center gap-3 px-2 py-2">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-500 text-lg font-black">AJ</div>
        <div><p class="font-black tracking-tight">AutoJournal</p><p class="text-xs text-slate-400">Личный автожурнал</p></div>
      </div>

      <nav class="mt-5 grid grid-cols-3 gap-2 sm:grid-cols-5 lg:grid-cols-1" aria-label="Основная навигация">
        <button
          v-for="item in navItems"
          :key="item.id"
          class="flex items-center gap-3 rounded-xl px-3 py-2.5 text-left text-sm font-bold transition"
          :class="activePage === item.id ? 'bg-white text-ink' : 'text-slate-300 hover:bg-white/10 hover:text-white'"
          @click="activePage = item.id"
        >
          <span class="flex h-7 w-7 items-center justify-center rounded-lg" :class="activePage === item.id ? 'bg-brand-50 text-brand-700' : 'bg-white/5'">{{ item.icon }}</span>
          <span class="hidden sm:inline">{{ item.label }}</span>
        </button>
      </nav>

      <div class="mt-6 hidden rounded-xl bg-white/5 p-4 lg:block">
        <p class="text-xs font-bold uppercase tracking-wider text-slate-400">Активный автомобиль</p>
        <p class="mt-2 truncate text-sm font-bold">{{ selectedVehicle ? `${selectedVehicle.make} ${selectedVehicle.model}` : 'Не добавлен' }}</p>
        <p v-if="selectedVehicle" class="mt-1 text-xs text-slate-400">{{ selectedVehicle.year }} · {{ selectedVehicle.odometer.toLocaleString('ru-RU') }} км</p>
      </div>
    </aside>

    <div class="lg:col-start-2">
      <header class="sticky top-0 z-10 flex min-h-16 items-center justify-between gap-4 border-b border-slate-200/80 bg-white/90 px-5 backdrop-blur md:px-8">
        <p class="text-sm font-black text-slate-700">{{ pageTitle }}</p>
        <select v-if="vehicles.length" v-model.number="selectedVehicleId" class="field max-w-xs" aria-label="Выбранный автомобиль">
          <option v-for="vehicle in vehicles" :key="vehicle.id" :value="vehicle.id">{{ vehicle.make }} {{ vehicle.model }} · {{ vehicle.year }}</option>
        </select>
        <button v-else class="btn-primary" @click="activePage = 'vehicle'">Добавить автомобиль</button>
      </header>

      <main class="mx-auto max-w-[1500px] p-5 md:p-8">
        <DashboardView v-if="activePage === 'dashboard'" :vehicle="selectedVehicle" @notify="notify" />
        <ExpensesView v-else-if="activePage === 'expenses'" :vehicle="selectedVehicle" @notify="notify" />
        <RemindersView v-else-if="activePage === 'reminders'" :vehicle="selectedVehicle" @notify="notify" />
        <VehicleProfileView v-else-if="activePage === 'vehicle'" :vehicle="selectedVehicle" @notify="notify" @vehicles-changed="loadVehicles" />
        <BackupsView v-else :vehicle="selectedVehicle" @notify="notify" />
      </main>
    </div>

    <div class="fixed right-4 top-20 z-50 grid w-[min(380px,calc(100vw-2rem))] gap-3" aria-live="polite">
      <div
        v-for="notice in notices"
        :key="notice.id"
        class="rounded-2xl border p-4 shadow-panel"
        :class="{
          'border-emerald-200 bg-emerald-50 text-emerald-900': notice.tone === 'success',
          'border-rose-200 bg-rose-50 text-rose-900': notice.tone === 'error',
          'border-amber-200 bg-amber-50 text-amber-950': notice.tone === 'info',
        }"
      >
        <div class="flex items-start justify-between gap-3"><p class="text-sm font-bold">{{ notice.message }}</p><button class="text-lg leading-none opacity-60 hover:opacity-100" aria-label="Закрыть" @click="dismissNotice(notice.id)">×</button></div>
      </div>
    </div>
  </div>
</template>
