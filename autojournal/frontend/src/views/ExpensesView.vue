<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { domain } from '../../wailsjs/go/models'
import {
  AddExpense,
  DeleteExpense,
  ListExpenseCategories,
  ListVehicleExpenses,
} from '../../wailsjs/go/handler/App'
import EmptyState from '../components/EmptyState.vue'
import { categoryName, formatDate, formatMoney, formatNumber, todayInputValue } from '../lib/format'
import type { ExpenseCategoryRecord, ExpenseRecord, VehicleRecord } from '../types'

const props = defineProps<{ vehicle: VehicleRecord | null }>()
const emit = defineEmits<{ (event: 'notify', message: string, tone?: 'success' | 'error' | 'info'): void }>()

const expenses = ref<ExpenseRecord[]>([])
const categories = ref<ExpenseCategoryRecord[]>([])
const saving = ref(false)
const form = reactive({
  categoryId: 0,
  amount: 0,
  odometerAt: 0,
  date: todayInputValue(),
  description: '',
})

const categoryById = computed(() => new Map(categories.value.map(item => [item.id, item.name])))

async function load() {
  if (!props.vehicle) {
    expenses.value = []
    return
  }
  try {
    const [expenseResult, categoryResult] = await Promise.all([
      ListVehicleExpenses(props.vehicle.id),
      ListExpenseCategories(),
    ])
    expenses.value = (expenseResult ?? []) as ExpenseRecord[]
    categories.value = (categoryResult ?? []) as ExpenseCategoryRecord[]
    if (!form.categoryId && categories.value.length) form.categoryId = categories.value[0].id
    form.odometerAt = props.vehicle.odometer
  } catch (error) {
    emit('notify', String(error), 'error')
  }
}

async function submit() {
  if (!props.vehicle) return
  saving.value = true
  try {
    await AddExpense(new domain.Expense({
      id: 0,
      vehicleId: props.vehicle.id,
      categoryId: Number(form.categoryId),
      amount: Math.round(Number(form.amount)),
      odometerAt: Math.round(Number(form.odometerAt)),
      date: new Date(`${form.date}T12:00:00`).toISOString(),
      description: form.description.trim(),
    }))
    form.amount = 0
    form.description = ''
    await load()
    emit('notify', 'Расход добавлен', 'success')
  } catch (error) {
    emit('notify', String(error), 'error')
  } finally {
    saving.value = false
  }
}

async function remove(id: number) {
  if (!window.confirm('Удалить эту запись расхода?')) return
  try {
    await DeleteExpense(id)
    await load()
    emit('notify', 'Расход удалён', 'success')
  } catch (error) {
    emit('notify', String(error), 'error')
  }
}

watch(() => props.vehicle?.id, load, { immediate: true })
</script>

<template>
  <section class="space-y-6">
    <header>
      <p class="text-sm font-bold uppercase tracking-[0.16em] text-brand-600">Журнал</p>
      <h1 class="mt-1 text-3xl font-black tracking-tight">Расходы</h1>
      <p class="mt-2 text-sm text-slate-500">Записывайте каждую трату и следите за историей обслуживания.</p>
    </header>

    <EmptyState v-if="!vehicle" title="Автомобиль не выбран" text="Добавьте автомобиль, чтобы вести журнал расходов." />

    <template v-else>
      <form class="panel grid gap-4 p-5 md:grid-cols-2 xl:grid-cols-5" @submit.prevent="submit">
        <label><span class="label">Категория</span><select v-model.number="form.categoryId" class="field" required><option v-for="item in categories" :key="item.id" :value="item.id">{{ categoryName(item.name) }}</option></select></label>
        <label><span class="label">Сумма, ₽</span><input v-model.number="form.amount" class="field" type="number" min="1" required></label>
        <label><span class="label">Пробег, км</span><input v-model.number="form.odometerAt" class="field" type="number" min="0" required></label>
        <label><span class="label">Дата</span><input v-model="form.date" class="field" type="date" required></label>
        <div class="flex items-end"><button class="btn-primary w-full" :disabled="saving">{{ saving ? 'Сохраняем…' : 'Добавить расход' }}</button></div>
        <label class="md:col-span-2 xl:col-span-5"><span class="label">Комментарий</span><input v-model="form.description" class="field" maxlength="250" placeholder="Например: замена масла и фильтра"></label>
      </form>

      <article class="panel overflow-hidden">
        <div class="border-b border-slate-100 px-5 py-4"><h2 class="text-lg font-black">История</h2></div>
        <div v-if="expenses.length" class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-400"><tr><th class="px-5 py-3">Дата</th><th class="px-5 py-3">Категория</th><th class="px-5 py-3">Описание</th><th class="px-5 py-3">Пробег</th><th class="px-5 py-3 text-right">Сумма</th><th class="px-5 py-3"></th></tr></thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="expense in expenses" :key="expense.id" class="hover:bg-slate-50/60">
                <td class="whitespace-nowrap px-5 py-4 text-slate-500">{{ formatDate(expense.date) }}</td>
                <td class="px-5 py-4 font-bold">{{ categoryName(categoryById.get(expense.categoryId) ?? 'Другое') }}</td>
                <td class="max-w-xs truncate px-5 py-4 text-slate-500">{{ expense.description || '—' }}</td>
                <td class="whitespace-nowrap px-5 py-4">{{ formatNumber(expense.odometerAt) }} км</td>
                <td class="whitespace-nowrap px-5 py-4 text-right font-black">{{ formatMoney(expense.amount) }}</td>
                <td class="px-5 py-4 text-right"><button class="text-xs font-bold text-rose-600 hover:text-rose-800" @click="remove(expense.id)">Удалить</button></td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="p-5"><EmptyState title="Записей пока нет" text="Добавьте первый расход в форме выше." /></div>
      </article>
    </template>
  </section>
</template>
