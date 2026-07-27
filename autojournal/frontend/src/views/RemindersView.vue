<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { domain } from '../../wailsjs/go/models'
import {
  CreateReminder,
  DeleteReminder,
  ListVehicleReminders,
  UpdateReminder,
} from '../../wailsjs/go/handler/App'
import EmptyState from '../components/EmptyState.vue'
import { formatDate, formatNumber, reminderTypeName, todayInputValue } from '../lib/format'
import type { ReminderRecord, VehicleRecord } from '../types'

const props = defineProps<{ vehicle: VehicleRecord | null }>()
const emit = defineEmits<{ (event: 'notify', message: string, tone?: 'success' | 'error' | 'info'): void }>()

const reminders = ref<ReminderRecord[]>([])
const saving = ref(false)
const editingId = ref<number | null>(null)
const form = reactive({
  title: '',
  reminderType: 'oil_change',
  intervalKm: 10_000 as number | null,
  intervalDays: null as number | null,
  lastDoneOdometer: 0 as number | null,
  lastDoneDate: todayInputValue(),
})

async function load() {
  if (!props.vehicle) {
    reminders.value = []
    return
  }
  try {
    reminders.value = (await ListVehicleReminders(props.vehicle.id) ?? []) as ReminderRecord[]
    form.lastDoneOdometer = props.vehicle.odometer
  } catch (error) {
    emit('notify', String(error), 'error')
  }
}

async function submit() {
  if (!props.vehicle) return
  const intervalKm = form.intervalKm ? Math.round(Number(form.intervalKm)) : null
  const intervalDays = form.intervalDays ? Math.round(Number(form.intervalDays)) : null
  const wasEditing = editingId.value !== null
  saving.value = true
  try {
    const current = reminders.value.find(item => item.id === editingId.value)
    const payload = new domain.Reminder({
      id: editingId.value ?? 0,
      vehicleId: props.vehicle.id,
      title: form.title.trim(),
      reminderType: form.reminderType,
      intervalKm,
      intervalDays,
      lastDoneOdometer: intervalKm === null ? null : Math.round(Number(form.lastDoneOdometer)),
      lastDoneDate: intervalDays === null ? null : new Date(`${form.lastDoneDate}T12:00:00`).toISOString(),
      nextDueDate: null,
      nextDueOdometer: null,
      isActive: current?.isActive ?? true,
      createdAt: current?.createdAt,
    })
    if (editingId.value) {
      await UpdateReminder(payload)
    } else {
      await CreateReminder(payload)
    }
    resetForm()
    await load()
    emit('notify', wasEditing ? 'Напоминание обновлено' : 'Напоминание создано', 'success')
  } catch (error) {
    emit('notify', String(error), 'error')
  } finally {
    saving.value = false
  }
}

function resetForm() {
  editingId.value = null
  form.title = ''
  form.reminderType = 'oil_change'
  form.intervalKm = 10_000
  form.intervalDays = null
  form.lastDoneOdometer = props.vehicle?.odometer ?? 0
  form.lastDoneDate = todayInputValue()
}

function edit(item: ReminderRecord) {
  editingId.value = item.id
  form.title = item.title
  form.reminderType = item.reminderType
  form.intervalKm = item.intervalKm
  form.intervalDays = item.intervalDays
  form.lastDoneOdometer = item.lastDoneOdometer
  form.lastDoneDate = item.lastDoneDate
    ? new Date(item.lastDoneDate as string).toISOString().slice(0, 10)
    : todayInputValue()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function update(item: ReminderRecord, changes: Partial<ReminderRecord>, message: string) {
  try {
    await UpdateReminder(new domain.Reminder({ ...item, ...changes }))
    await load()
    emit('notify', message, 'success')
  } catch (error) {
    emit('notify', String(error), 'error')
  }
}

async function markDone(item: ReminderRecord) {
  if (!props.vehicle) return
  await update(item, {
    lastDoneDate: item.intervalDays === null ? null : new Date().toISOString(),
    lastDoneOdometer: item.intervalKm === null ? null : props.vehicle.odometer,
    isActive: true,
  }, 'Срок напоминания пересчитан')
}

async function remove(id: number) {
  if (!window.confirm('Удалить напоминание?')) return
  try {
    await DeleteReminder(id)
    await load()
    emit('notify', 'Напоминание удалено', 'success')
  } catch (error) {
    emit('notify', String(error), 'error')
  }
}

watch(() => props.vehicle?.id, load, { immediate: true })
</script>

<template>
  <section class="space-y-6">
    <header>
      <p class="text-sm font-bold uppercase tracking-[0.16em] text-brand-600">Техническое обслуживание</p>
      <h1 class="mt-1 text-3xl font-black tracking-tight">Напоминания</h1>
      <p class="mt-2 text-sm text-slate-500">Приложение предупредит за 7 дней или за 500 км до срока.</p>
    </header>

    <EmptyState v-if="!vehicle" title="Автомобиль не выбран" text="Напоминания привязываются к конкретному автомобилю." />

    <template v-else>
      <form class="panel grid gap-4 p-5 md:grid-cols-2 xl:grid-cols-4" @submit.prevent="submit">
        <label class="md:col-span-2"><span class="label">Название</span><input v-model="form.title" class="field" required maxlength="120" placeholder="Замена моторного масла"></label>
        <label><span class="label">Тип</span><select v-model="form.reminderType" class="field"><option value="oil_change">Замена масла</option><option value="tire_rotation">Перестановка шин</option><option value="insurance">Страхование</option><option value="custom">Другое</option></select></label>
        <div class="flex items-end gap-2"><button class="btn-primary flex-1" :disabled="saving">{{ saving ? 'Сохраняем…' : editingId ? 'Сохранить' : 'Создать' }}</button><button v-if="editingId" type="button" class="btn-secondary" @click="resetForm">Отмена</button></div>
        <label><span class="label">Интервал, км</span><input v-model.number="form.intervalKm" class="field" type="number" min="1" placeholder="Не учитывать"></label>
        <label><span class="label">Последний пробег</span><input v-model.number="form.lastDoneOdometer" class="field" type="number" min="0" :disabled="!form.intervalKm"></label>
        <label><span class="label">Интервал, дней</span><input v-model.number="form.intervalDays" class="field" type="number" min="1" placeholder="Не учитывать"></label>
        <label><span class="label">Дата выполнения</span><input v-model="form.lastDoneDate" class="field" type="date" :disabled="!form.intervalDays"></label>
        <p class="text-xs text-slate-500 md:col-span-2 xl:col-span-4">Заполните хотя бы один интервал: по пробегу или по дням.</p>
      </form>

      <div v-if="reminders.length" class="grid gap-4 lg:grid-cols-2">
        <article v-for="item in reminders" :key="item.id" class="panel p-5" :class="{ 'opacity-60': !item.isActive }">
          <div class="flex items-start justify-between gap-4">
            <div>
              <span class="rounded-full bg-brand-50 px-2.5 py-1 text-xs font-bold text-brand-700">{{ reminderTypeName(item.reminderType) }}</span>
              <h2 class="mt-3 text-lg font-black">{{ item.title }}</h2>
            </div>
            <span class="rounded-full px-2.5 py-1 text-xs font-bold" :class="item.isActive ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-500'">{{ item.isActive ? 'Активно' : 'Выключено' }}</span>
          </div>
          <dl class="mt-5 grid grid-cols-2 gap-3 text-sm">
            <div class="rounded-xl bg-slate-50 p-3"><dt class="text-xs text-slate-400">Следующая дата</dt><dd class="mt-1 font-bold">{{ formatDate(item.nextDueDate) }}</dd></div>
            <div class="rounded-xl bg-slate-50 p-3"><dt class="text-xs text-slate-400">Следующий пробег</dt><dd class="mt-1 font-bold">{{ item.nextDueOdometer === null ? '—' : `${formatNumber(item.nextDueOdometer)} км` }}</dd></div>
          </dl>
          <div class="mt-5 flex flex-wrap gap-2">
            <button class="btn-primary" @click="markDone(item)">Выполнено сегодня</button>
            <button class="btn-secondary" @click="edit(item)">Изменить</button>
            <button class="btn-secondary" @click="update(item, { isActive: !item.isActive }, item.isActive ? 'Напоминание выключено' : 'Напоминание включено')">{{ item.isActive ? 'Выключить' : 'Включить' }}</button>
            <button class="btn-danger" @click="remove(item.id)">Удалить</button>
          </div>
        </article>
      </div>
      <EmptyState v-else title="Напоминаний пока нет" text="Создайте первое правило обслуживания в форме выше." />
    </template>
  </section>
</template>
