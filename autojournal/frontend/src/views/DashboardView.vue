<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import {
  GetDueReminders,
  GetExpenseStats,
  ListExpenseCategories,
  ListVehicleExpenses,
} from '../../wailsjs/go/handler/App'
import EmptyState from '../components/EmptyState.vue'
import ExpenseChart from '../components/ExpenseChart.vue'
import { categoryName, formatDate, formatMoney, formatMonth, formatNumber } from '../lib/format'
import type {
  DueReminderRecord,
  ExpenseCategoryRecord,
  ExpenseRecord,
  ExpenseStatsRecord,
  VehicleRecord,
} from '../types'

const props = defineProps<{ vehicle: VehicleRecord | null }>()
const emit = defineEmits<{ (event: 'notify', message: string, tone?: 'success' | 'error' | 'info'): void }>()

const loading = ref(false)
const stats = ref<ExpenseStatsRecord>({ totalAmount: 0, byCategory: [], byMonth: [] })
const expenses = ref<ExpenseRecord[]>([])
const dueReminders = ref<DueReminderRecord[]>([])
const categories = ref<ExpenseCategoryRecord[]>([])

const recentExpenses = computed(() => expenses.value.slice(0, 5))
const categoryById = computed(() => new Map(categories.value.map(item => [item.id, item.name])))

async function load() {
  if (!props.vehicle) {
    stats.value = { totalAmount: 0, byCategory: [], byMonth: [] }
    expenses.value = []
    dueReminders.value = []
    return
  }

  loading.value = true
  try {
    const [statsResult, expensesResult, dueResult, categoriesResult] = await Promise.all([
      GetExpenseStats(props.vehicle.id),
      ListVehicleExpenses(props.vehicle.id),
      GetDueReminders(),
      ListExpenseCategories(),
    ])
    stats.value = {
      ...statsResult,
      byCategory: statsResult.byCategory ?? [],
      byMonth: statsResult.byMonth ?? [],
    } as ExpenseStatsRecord
    expenses.value = (expensesResult ?? []) as ExpenseRecord[]
    dueReminders.value = (dueResult ?? []).filter(item => item.reminder.vehicleId === props.vehicle?.id) as DueReminderRecord[]
    categories.value = (categoriesResult ?? []) as ExpenseCategoryRecord[]
  } catch (error) {
    emit('notify', String(error), 'error')
  } finally {
    loading.value = false
  }
}

watch(() => props.vehicle?.id, load, { immediate: true })
</script>

<template>
  <section class="space-y-6">
    <header>
      <p class="text-sm font-bold uppercase tracking-[0.16em] text-brand-600">Обзор</p>
      <h1 class="mt-1 text-3xl font-black tracking-tight text-ink">Панель управления</h1>
      <p class="mt-2 text-sm text-slate-500">Расходы, пробег и ближайшие работы в одном месте.</p>
    </header>

    <EmptyState
      v-if="!vehicle"
      title="Сначала добавьте автомобиль"
      text="Откройте раздел «Автомобиль» и заполните профиль — после этого здесь появится сводка."
    />

    <template v-else>
      <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <article class="panel p-5">
          <p class="text-xs font-bold uppercase tracking-wider text-slate-400">Всего расходов</p>
          <p class="mt-3 text-2xl font-black">{{ formatMoney(stats.totalAmount) }}</p>
        </article>
        <article class="panel p-5">
          <p class="text-xs font-bold uppercase tracking-wider text-slate-400">Текущий пробег</p>
          <p class="mt-3 text-2xl font-black">{{ formatNumber(vehicle.odometer) }} км</p>
        </article>
        <article class="panel p-5">
          <p class="text-xs font-bold uppercase tracking-wider text-slate-400">Ближайшие работы</p>
          <p class="mt-3 text-2xl font-black">{{ dueReminders.length }}</p>
        </article>
        <article class="panel p-5">
          <p class="text-xs font-bold uppercase tracking-wider text-slate-400">Записей расходов</p>
          <p class="mt-3 text-2xl font-black">{{ expenses.length }}</p>
        </article>
      </div>

      <div class="grid gap-6 xl:grid-cols-2">
        <article class="panel p-5">
          <div class="mb-5">
            <h2 class="text-lg font-black">Расходы по месяцам</h2>
            <p class="text-sm text-slate-500">Динамика затрат на выбранный автомобиль</p>
          </div>
          <ExpenseChart
            v-if="stats.byMonth.length"
            :labels="stats.byMonth.map(item => formatMonth(item.month))"
            :values="stats.byMonth.map(item => item.totalAmount)"
          />
          <EmptyState v-else title="Пока нет данных" text="Добавьте первый расход, чтобы построить график." />
        </article>

        <article class="panel p-5">
          <div class="mb-5">
            <h2 class="text-lg font-black">По категориям</h2>
            <p class="text-sm text-slate-500">Как распределены все расходы</p>
          </div>
          <ExpenseChart
            v-if="stats.byCategory.length"
            type="doughnut"
            :labels="stats.byCategory.map(item => categoryName(item.categoryName))"
            :values="stats.byCategory.map(item => item.totalAmount)"
          />
          <EmptyState v-else title="Категории пусты" text="Диаграмма появится после добавления расходов." />
        </article>
      </div>

      <div class="grid gap-6 xl:grid-cols-[1.4fr_1fr]">
        <article class="panel overflow-hidden">
          <div class="border-b border-slate-100 px-5 py-4">
            <h2 class="text-lg font-black">Последние расходы</h2>
          </div>
          <div v-if="recentExpenses.length" class="overflow-x-auto">
            <table class="w-full text-left text-sm">
              <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-400">
                <tr><th class="px-5 py-3">Дата</th><th class="px-5 py-3">Категория</th><th class="px-5 py-3 text-right">Сумма</th></tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="expense in recentExpenses" :key="expense.id">
                  <td class="px-5 py-4 text-slate-500">{{ formatDate(expense.date) }}</td>
                  <td class="px-5 py-4 font-bold">{{ categoryName(categoryById.get(expense.categoryId) ?? 'Другое') }}</td>
                  <td class="px-5 py-4 text-right font-black">{{ formatMoney(expense.amount) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-else class="p-5"><EmptyState title="Расходов нет" /></div>
        </article>

        <article class="panel p-5">
          <h2 class="text-lg font-black">Ближайшие работы</h2>
          <div v-if="dueReminders.length" class="mt-4 space-y-3">
            <div v-for="item in dueReminders" :key="item.reminder.id" class="rounded-xl bg-amber-50 p-4">
              <p class="font-bold text-amber-950">{{ item.reminder.title }}</p>
              <p class="mt-1 text-xs text-amber-700">
                <span v-if="item.dueByDate">до {{ formatDate(item.reminder.nextDueDate) }}</span>
                <span v-if="item.dueByDate && item.dueByOdometer"> · </span>
                <span v-if="item.dueByOdometer">к {{ formatNumber(item.reminder.nextDueOdometer ?? 0) }} км</span>
              </p>
            </div>
          </div>
          <div v-else class="mt-4"><EmptyState title="Всё вовремя" text="Ближайших работ пока нет." /></div>
        </article>
      </div>
      <p v-if="loading" class="text-sm text-slate-400">Обновляем данные…</p>
    </template>
  </section>
</template>
