export type PageName = 'dashboard' | 'expenses' | 'reminders' | 'vehicle' | 'backups'

export interface VehicleRecord {
  id: number
  vin: string
  make: string
  model: string
  year: number
  color?: string
  engineVolume: number
  fuelType: number
  odometer: number
  notes: string
  createdAt: unknown
  updatedAt: unknown
}

export interface ExpenseRecord {
  id: number
  vehicleId: number
  categoryId: number
  amount: number
  odometerAt: number
  date: unknown
  description: string
  createdAt: unknown
}

export interface ExpenseCategoryRecord {
  id: number
  name: string
  icon: string
}

export interface ExpenseStatsRecord {
  totalAmount: number
  byCategory: Array<{
    categoryId: number
    categoryName: string
    totalAmount: number
  }>
  byMonth: Array<{
    month: string
    totalAmount: number
  }>
}

export interface ReminderRecord {
  id: number
  vehicleId: number
  title: string
  reminderType: string
  intervalKm: number | null
  intervalDays: number | null
  lastDoneOdometer: number | null
  lastDoneDate: unknown | null
  nextDueDate: unknown | null
  nextDueOdometer: number | null
  isActive: boolean
  createdAt: unknown
}

export interface DueReminderRecord {
  reminder: ReminderRecord
  dueByDate: boolean
  dueByOdometer: boolean
}

export interface BackupRecord {
  id: number
  filePath: string
  note: string
  createdAt: unknown
}

export type Notice = {
  id: number
  message: string
  tone: 'success' | 'error' | 'info'
}
