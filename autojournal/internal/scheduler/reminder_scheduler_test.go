package scheduler

import (
	"autojournal/internal/domain"
	"context"
	"testing"
	"time"
)

func TestReminderScheduler_CheckNow_DeduplicatesUnchangedReminders(t *testing.T) {
	nextDate := time.Date(2026, time.January, 17, 12, 0, 0, 0, time.UTC)
	provider := &dueReminderProviderStub{reminders: []domain.DueReminder{
		{
			Reminder: domain.Reminder{
				ID:          1,
				VehicleID:   1,
				Title:       "Insurance",
				NextDueDate: &nextDate,
			},
			DueByDate: true,
		},
	}}

	events := make([]domain.DueReminder, 0)
	scheduler := NewReminderScheduler(provider, func(_ context.Context, event string, payload any) {
		if event != reminderDueEvent {
			t.Fatalf("event = %q, want %q", event, reminderDueEvent)
		}
		reminder, ok := payload.(domain.DueReminder)
		if !ok {
			t.Fatalf("payload has type %T, want domain.DueReminder", payload)
		}
		events = append(events, reminder)
	})
	scheduler.now = func() time.Time {
		return time.Date(2026, time.January, 10, 12, 0, 0, 0, time.UTC)
	}

	if err := scheduler.CheckNow(context.Background()); err != nil {
		t.Fatalf("first check: %v", err)
	}
	if err := scheduler.CheckNow(context.Background()); err != nil {
		t.Fatalf("second check: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("emitted %d events, want 1", len(events))
	}

	provider.reminders[0].DueByOdometer = true
	if err := scheduler.CheckNow(context.Background()); err != nil {
		t.Fatalf("check after reminder state changed: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("emitted %d events after state changed, want 2", len(events))
	}
}

type dueReminderProviderStub struct {
	reminders []domain.DueReminder
}

func (s *dueReminderProviderStub) GetDueReminders(context.Context, time.Time) ([]domain.DueReminder, error) {
	return s.reminders, nil
}
