CREATE TABLE IF NOT EXISTS reminders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    vehicle_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    reminder_type TEXT NOT NULL,
    interval_km INTEGER,
    interval_days INTEGER,
    last_done_odometer INTEGER,
    last_done_date DATETIME,
    next_due_date DATETIME,
    next_due_odometer INTEGER,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id) ON DELETE CASCADE,

    CHECK (interval_km IS NULL OR interval_km > 0),
    CHECK (interval_days IS NULL OR interval_days > 0),
    CHECK (last_done_odometer IS NULL OR last_done_odometer >= 0),
    CHECK (is_active IN (0, 1))
);

CREATE INDEX IF NOT EXISTS idx_reminders_vehicle_active
ON reminders(vehicle_id, is_active);

CREATE INDEX IF NOT EXISTS idx_reminders_active_due_date
ON reminders(is_active, next_due_date);