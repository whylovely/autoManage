<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { domain } from '../../wailsjs/go/models'
import { CreateVehicle, DeleteVehicle, UpdateVehicle } from '../../wailsjs/go/handler/App'
import { fuelTypeName } from '../lib/format'
import type { VehicleRecord } from '../types'

const props = defineProps<{ vehicle: VehicleRecord | null }>()
const emit = defineEmits<{
  (event: 'notify', message: string, tone?: 'success' | 'error' | 'info'): void
  (event: 'vehicles-changed', preferredId?: number): void
}>()

const saving = ref(false)
const form = reactive({
  vin: '', make: '', model: '', year: new Date().getFullYear(), color: '',
  engineVolume: 0, fuelType: 0, odometer: 0, notes: '',
})

function fillForm() {
  const vehicle = props.vehicle
  form.vin = vehicle?.vin ?? ''
  form.make = vehicle?.make ?? ''
  form.model = vehicle?.model ?? ''
  form.year = vehicle?.year ?? new Date().getFullYear()
  form.color = vehicle?.color ?? ''
  form.engineVolume = vehicle?.engineVolume ?? 0
  form.fuelType = vehicle?.fuelType ?? 0
  form.odometer = vehicle?.odometer ?? 0
  form.notes = vehicle?.notes ?? ''
}

async function submit() {
  saving.value = true
  try {
    const payload = new domain.Vehicle({
      id: props.vehicle?.id ?? 0,
      vin: form.vin.trim().toUpperCase(),
      make: form.make.trim(),
      model: form.model.trim(),
      year: Number(form.year),
      color: form.color.trim(),
      engineVolume: Math.round(Number(form.engineVolume)),
      fuelType: Number(form.fuelType),
      odometer: Math.round(Number(form.odometer)),
      notes: form.notes.trim(),
    })
    const saved = props.vehicle ? await UpdateVehicle(payload) : await CreateVehicle(payload)
    emit('vehicles-changed', saved.id)
    emit('notify', props.vehicle ? 'Профиль автомобиля обновлён' : 'Автомобиль добавлен', 'success')
  } catch (error) {
    emit('notify', String(error), 'error')
  } finally {
    saving.value = false
  }
}

async function remove() {
  if (!props.vehicle || !window.confirm('Удалить автомобиль вместе с его расходами и напоминаниями?')) return
  try {
    await DeleteVehicle(props.vehicle.id)
    emit('vehicles-changed')
    emit('notify', 'Автомобиль удалён', 'success')
  } catch (error) {
    emit('notify', String(error), 'error')
  }
}

watch(() => props.vehicle?.id, fillForm, { immediate: true })
</script>

<template>
  <section class="space-y-6">
    <header class="flex flex-wrap items-end justify-between gap-4">
      <div>
        <p class="text-sm font-bold uppercase tracking-[0.16em] text-brand-600">Гараж</p>
        <h1 class="mt-1 text-3xl font-black tracking-tight">{{ vehicle ? 'Профиль автомобиля' : 'Новый автомобиль' }}</h1>
        <p class="mt-2 text-sm text-slate-500">Характеристики и актуальный пробег.</p>
      </div>
      <button v-if="vehicle" class="btn-danger" @click="remove">Удалить автомобиль</button>
    </header>

    <form class="panel p-6" @submit.prevent="submit">
      <div class="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
        <label><span class="label">VIN</span><input v-model="form.vin" class="field uppercase" minlength="17" maxlength="17" required placeholder="1HGCM82633A004352"></label>
        <label><span class="label">Марка</span><input v-model="form.make" class="field" required placeholder="Toyota"></label>
        <label><span class="label">Модель</span><input v-model="form.model" class="field" required placeholder="Camry"></label>
        <label><span class="label">Год</span><input v-model.number="form.year" class="field" type="number" min="1886" :max="new Date().getFullYear() + 1" required></label>
        <label><span class="label">Цвет</span><input v-model="form.color" class="field" placeholder="Серебристый"></label>
        <label><span class="label">Объём двигателя, см³</span><input v-model.number="form.engineVolume" class="field" type="number" min="0" step="1"></label>
        <label><span class="label">Тип топлива</span><select v-model.number="form.fuelType" class="field"><option v-for="value in [0, 1, 2, 3, 4, 5]" :key="value" :value="value">{{ fuelTypeName(value) }}</option></select></label>
        <label><span class="label">Пробег, км</span><input v-model.number="form.odometer" class="field" type="number" min="0" required><small v-if="vehicle" class="mt-1 block text-xs text-slate-400">Пробег нельзя уменьшать.</small></label>
        <label class="md:col-span-2 xl:col-span-3"><span class="label">Заметки</span><textarea v-model="form.notes" class="field min-h-[110px] resize-y" maxlength="1000" placeholder="Комплектация, особенности обслуживания…"></textarea></label>
      </div>
      <div class="mt-6 flex justify-end"><button class="btn-primary min-w-40" :disabled="saving">{{ saving ? 'Сохраняем…' : vehicle ? 'Сохранить изменения' : 'Добавить автомобиль' }}</button></div>
    </form>
  </section>
</template>
