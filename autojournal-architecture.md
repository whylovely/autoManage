# Архитектура проекта: Авто-журнал (Go)

Локальное standalone-приложение для ведения учёта по автомобилю. Язык: Go. UI: Wails v2. БД: SQLite.

---

## 1. Рекомендации по стеку

### SQLite

- `modernc.org/sqlite` — чистый Go, без CGO. Рекомендуется для standalone без зависимости от компилятора.
- `github.com/mattn/go-sqlite3` — CGO, быстрее, требует GCC.
- `github.com/jmoiron/sqlx` — поверх database/sql, удобный маппинг строк в структуры.

### UI-фреймворк

**Рекомендация: Wails v2.**

- Нативный WebView (Chromium/WebKit).
- TypeScript-фронтенд, биндинги Go↔JS генерируются автоматически.
- Результат — один бинарник.
- Fyne проще в сборке, но выглядит устаревше.
- TUI (Bubbletea) — хорошо для CLI, но для журнала с таблицами и формами UI удобнее.

### Вспомогательные библиотеки

| Библиотека | Назначение |
|---|---|
| `github.com/golang-migrate/migrate` | Миграции БД |
| `github.com/robfig/cron/v3` | Планировщик напоминаний |
| `github.com/go-playground/validator/v10` | Валидация данных |
| `encoding/csv` + `encoding/json` (stdlib) | Экспорт данных |

---

## 2. Структура проекта

```
autojournal/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── vehicle.go
│   │   ├── expense.go
│   │   ├── reminder.go
│   │   └── backup.go
│   ├── storage/
│   │   ├── sqlite.go
│   │   ├── vehicle_repo.go
│   │   ├── expense_repo.go
│   │   └── reminder_repo.go
│   ├── service/
│   │   ├── vehicle_service.go
│   │   ├── expense_service.go
│   │   ├── reminder_service.go
│   │   └── backup_service.go
│   ├── handler/
│   │   └── app.go          # Wails-биндинги
│   └── scheduler/
│       └── scheduler.go
├── frontend/
│   └── src/                # Vue/TS фронтенд
├── migrations/
│   ├── 001_init.up.sql
│   ├── 001_init.down.sql
│   └── 002_seed_categories.up.sql
├── go.mod
├── go.sum
├── wails.json
└── Makefile
```

**Ключевой принцип**: `internal/domain` содержит только бизнес-структуры и интерфейсы репозиториев. `internal/storage` — реализации. `internal/service` — бизнес-логика. `internal/handler` — Wails-биндинги (экспортируемые методы, вызываемые из JS).

---

## 3. Схема базы данных

### Таблица `vehicles`

```sql
CREATE TABLE vehicles (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    vin           TEXT,
    make          TEXT NOT NULL,
    model         TEXT NOT NULL,
    year          INTEGER NOT NULL,
    color         TEXT,
    engine_volume REAL,
    fuel_type     TEXT,
    odometer      INTEGER NOT NULL DEFAULT 0,
    notes         TEXT,
    created_at    TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at    TEXT NOT NULL DEFAULT (datetime('now'))
);
```

### Таблица `expense_categories`

```sql
CREATE TABLE expense_categories (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    icon TEXT
);
```

### Таблица `expenses`

```sql
CREATE TABLE expenses (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    vehicle_id    INTEGER NOT NULL REFERENCES vehicles(id),
    category_id   INTEGER NOT NULL REFERENCES expense_categories(id),
    amount        REAL NOT NULL,
    currency      TEXT NOT NULL DEFAULT 'RUB',
    odometer_at   REAL,
    date          TEXT NOT NULL,
    description   TEXT,
    receipt_path  TEXT,
    created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);
```

### Таблица `reminders`

```sql
CREATE TABLE reminders (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    vehicle_id          INTEGER NOT NULL REFERENCES vehicles(id),
    title               TEXT NOT NULL,
    reminder_type       TEXT NOT NULL,   -- oil_change | tire_rotation | insurance | custom
    interval_km         INTEGER,
    interval_days       INTEGER,
    last_done_odometer  INTEGER,
    last_done_date      TEXT,
    next_due_date       TEXT,
    next_due_odometer   INTEGER,
    is_active           INTEGER NOT NULL DEFAULT 1,
    created_at          TEXT NOT NULL DEFAULT (datetime('now'))
);
```

### Таблица `backups`

```sql
CREATE TABLE backups (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    file_path  TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    note       TEXT
);
```

### Связи

- `vehicles` → `expenses` (1:N)
- `vehicles` → `reminders` (1:N)
- `expense_categories` → `expenses` (1:N)

### Архитектурные решения

- `reminder_type` — enum-строка: `oil_change`, `tire_rotation`, `insurance`, `custom`.
- Напоминание срабатывает по любому из двух условий: **пробег ИЛИ дата** — то, что наступит раньше.
- `receipt_path` хранит локальный путь к фото чека.
- `backups` — лог всех созданных резервных копий.

---

## 4. Пошаговый план реализации

### Фаза 1. Инициализация и инфраструктура

**1.1 Инициализация Wails-проекта**

```bash
wails init -n autojournal -t vue-ts
```

Проверить `go.mod`, настроить `wails.json` (appID, productName, icon).

**1.2 Настройка зависимостей Go**

```bash
go install modernc.org/sqlite
go install github.com/jmoiron/sqlx@latest
go install github.com/golang-migrate/migrate/v4@latest
go install github.com/robfig/cron/v3@latest
go install github.com/go-playground/validator/v10@latest
```

**1.3 Создание структуры директорий**

Создать `internal/{domain,storage,service,handler,scheduler}`, `migrations/`, `Makefile` с целями `build`, `dev`, `test`.

---

### Фаза 2. Domain layer и интерфейсы

**2.1 Определить бизнес-структуры**

`Vehicle`, `Expense`, `ExpenseCategory`, `Reminder`, `Backup` в пакете `domain`. Поля с `db`-тегами для sqlx и `validate`-тегами для валидатора.

**2.2 Описать интерфейсы репозиториев**

`VehicleRepo`, `ExpenseRepo`, `ReminderRepo` — интерфейсы с методами CRUD. Позволяют подменять реализацию в тестах.

---

### Фаза 3. Storage layer — SQLite

**3.1 Инициализация БД и миграции**

`sqlite.go`: открыть БД в `~/.config/autojournal/data.db`, включить WAL и FK, запустить `migrate.Up()` при старте.

```go
// Обязательные PRAGMA при открытии коннекта
PRAGMA journal_mode=WAL;
PRAGMA foreign_keys=ON;
```

**3.2 Написать SQL-миграции**

- `001_init.up.sql` — создание всех таблиц.
- `002_seed_categories.up.sql` — стандартные категории расходов (топливо, запчасти, мойка, страховка).

**3.3 Реализовать репозитории**

`vehicle_repo.go`, `expense_repo.go`, `reminder_repo.go`. Для SELECT использовать `sqlx.Select`/`sqlx.Get`. Пагинация через `LIMIT`/`OFFSET`.

---

### Фаза 4. Service layer — бизнес-логика

**4.1 VehicleService**

Создание/редактирование профиля авто, обновление одометра, валидация VIN (regexp 17 символов).

**4.2 ExpenseService**

Добавление расходов с валидацией, агрегация по категориям/периоду, расчёт средних затрат на 100 км.

**4.3 ReminderService**

Расчёт `next_due` по одометру и дате:

```
nextDate     = lastDoneDate + intervalDays
nextOdometer = lastDoneOdometer + intervalKm
```

`GetDueReminders()` возвращает просроченные и ближайшие (в течение 7 дней / 500 км). При обновлении одометра авто — пересчитывать `next_due_odometer` для всех активных напоминаний.

**4.4 BackupService**

Атомарная копия БД без блокировки основного файла:

```sql
VACUUM INTO 'autojournal_2026-06-28.db';
```

Экспорт в CSV/JSON через `encoding/csv`, `encoding/json`. Путь к файлу через `os.UserConfigDir()`.

---

### Фаза 5. Планировщик и уведомления

**5.1 Cron-задача проверки напоминаний**

`robfig/cron`: запускать `CheckReminders()` каждые 30 минут. При открытии приложения — немедленная проверка.

**5.2 Отправка событий на фронтенд**

```go
runtime.EventsEmit(ctx, "reminder:due", payload)
```

Фронтенд подписывается и показывает тост/баннер.

---

### Фаза 6. Wails handler и фронтенд

**6.1 Реализовать App struct (Wails биндинги)**

Экспортируемые методы: `GetVehicles`, `AddExpense`, `GetExpenseStats`, `GetDueReminders`, `CreateBackup`. Вызываются из JS через `window.go.main.App.*`.

**6.2 Разработать Vue/TS фронтенд**

Страницы:
- **Dashboard** — сводка + активные напоминания.
- **Expenses** — таблица расходов + форма добавления.
- **Reminders** — управление напоминаниями о ТО.
- **Vehicle Profile** — профиль авто, VIN, характеристики.
- **Backups** — создание резервных копий и экспорт.

Стилизация: Tailwind CSS.

**6.3 Графики расходов**

Chart.js или ApexCharts: расходы по месяцам, круговая диаграмма по категориям. Данные приходят из Go через Wails.

---

### Фаза 7. Тесты и сборка

**7.1 Unit-тесты сервисов**

Мокировать репозитории через интерфейсы. Обязательно покрыть `ReminderService.CalculateNextDue()` — критичная логика. Использовать `testify`.

**7.2 Интеграционные тесты хранилища**

`TestMain` с in-memory SQLite (`:memory:`). Прогнать все миграции, протестировать репозитории на реальной БД.

**7.3 Сборка финального бинарника**

```bash
# Windows
wails build -platform windows/amd64

# macOS (универсальный)
wails build -platform darwin/universal

# Linux
wails build -platform linux/amd64

# Windows с инсталлятором NSIS
wails build -platform windows/amd64 -nsis
```

Результат — в `build/bin/`.

---

## 5. Важные детали реализации

**Путь к данным**: использовать `os.UserConfigDir()` + `autojournal/`. Это правильные пути для каждой ОС:
- macOS: `~/Library/Application Support/autojournal/`
- Windows: `%AppData%\autojournal\`
- Linux: `~/.config/autojournal/`

Фото чеков хранить в подпапке `receipts/`.

**WAL-режим SQLite**: обязательно включить `PRAGMA journal_mode=WAL` и `PRAGMA foreign_keys=ON` при открытии коннекта — без этого FK не работают, и запись блокирует чтение.

**Расчёт следующего ТО**: сохранять оба значения (`next_due_date` и `next_due_odometer`) в таблицу; при обновлении одометра авто пересчитывать `next_due_odometer` для всех активных напоминаний.

**Backup без остановки**: `VACUUM INTO 'backup_path.db'` — атомарная копия без блокировки основного файла.
