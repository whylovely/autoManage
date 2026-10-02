Таблица Машины
``` sql
TABLE vehicles (
    id            INTEGER PRIMARY KEY,
    vin           TEXT,
    make          TEXT,
    model         TEXT,
    year          INTEGER,
    color         TEXT,
    engine_volume INTEGER,
    fuel_type     INTEGER,
    odometer      INTEGER,
    notes         TEXT,
    created_at    DATETIME,
    updated_at    DATETIME
);
```
Таблица Категории
``` sql
TABLE expense_categories (
    id   INTEGER,
    name TEXT,
    icon TEXT
);
```
Таблица Расходы
``` sql
TABLE expenses (
    id          INTEGER PRIMARY KEY,
    vehicle_id  INTEGER,
    category_id INTEGER,
    amount      INTEGER,
    odometer_at INTEGER,
    date        DATETIME,
    description TEXT,
    created_at  DATETIME,
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
    FOREIGN KEY (category_id) REFERENCES expense_categories(id)
);
```
Таблица Бэкапы
``` sql
TABLE backups (
    id         INTEGER PRIMARY KEY,
    file_path  TEXT,
    note       TEXT,
    created_at DATETIME
);
```
Таблица Напоминания
``` sql
TABLE reminders (
    id                 INTEGER PRIMARY KEY,
    vehicle_id         INTEGER,
    title              TEXT,
    reminder_type      TEXT,
    interval_km        INTEGER,
    interval_days      INTEGER,
    last_done_odometer INTEGER,
    last_done_date     DATETIME,
    next_due_date      DATETIME,
    next_due_odometer  INTEGER,
    is_active          INTEGER,  
    created_at         DATETIME,
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
);
```