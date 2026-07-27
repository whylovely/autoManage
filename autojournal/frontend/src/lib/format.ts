const categoryNames: Record<string, string> = {
  fuel: 'Топливо',
  service: 'Сервис',
  insurance: 'Страхование',
  wash: 'Мойка',
  parts: 'Запчасти',
}

const reminderTypeNames: Record<string, string> = {
  oil_change: 'Замена масла',
  tire_rotation: 'Перестановка шин',
  insurance: 'Страхование',
  custom: 'Другое',
}

const fuelTypeNames: Record<number, string> = {
  0: 'Не указан',
  1: 'Бензин',
  2: 'Дизель',
  3: 'Гибрид',
  4: 'Электро',
  5: 'Газ',
}

export function formatMoney(value: number): string {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    maximumFractionDigits: 0,
  }).format(value || 0)
}

export function formatNumber(value: number): string {
  return new Intl.NumberFormat('ru-RU').format(value || 0)
}

export function formatDate(value: unknown): string {
  if (!value) return '—'
  const date = new Date(value as string)
  if (Number.isNaN(date.getTime())) return '—'
  return new Intl.DateTimeFormat('ru-RU', { dateStyle: 'medium' }).format(date)
}

export function formatDateTime(value: unknown): string {
  if (!value) return '—'
  const date = new Date(value as string)
  if (Number.isNaN(date.getTime())) return '—'
  return new Intl.DateTimeFormat('ru-RU', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date)
}

export function formatMonth(value: string): string {
  const date = new Date(`${value}-01T00:00:00`)
  return new Intl.DateTimeFormat('ru-RU', { month: 'short', year: '2-digit' }).format(date)
}

export function categoryName(value: string): string {
  return categoryNames[value] ?? value
}

export function reminderTypeName(value: string): string {
  return reminderTypeNames[value] ?? value
}

export function fuelTypeName(value: number): string {
  return fuelTypeNames[value] ?? 'Не указан'
}

export function todayInputValue(): string {
  const now = new Date()
  const offset = now.getTimezoneOffset() * 60_000
  return new Date(now.getTime() - offset).toISOString().slice(0, 10)
}
