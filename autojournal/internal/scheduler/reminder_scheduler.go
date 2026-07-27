package scheduler

import (
	"autojournal/internal/domain"
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

const reminderDueEvent = "reminder:due"

type DueReminderProvider interface {
	GetDueReminders(ctx context.Context, now time.Time) ([]domain.DueReminder, error)
}

type EventEmitter func(ctx context.Context, event string, payload any)

type ReminderScheduler struct {
	cron     *cron.Cron
	provider DueReminderProvider
	emit     EventEmitter
	now      func() time.Time

	checkMu sync.Mutex
	stateMu sync.Mutex
	startMu sync.Mutex
	started bool
	emitted map[string]struct{}
}

func NewReminderScheduler(provider DueReminderProvider, emit EventEmitter) *ReminderScheduler {
	return &ReminderScheduler{
		cron: cron.New(
			cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
		),
		provider: provider,
		emit:     emit,
		now:      time.Now,
		emitted:  make(map[string]struct{}),
	}
}

func (s *ReminderScheduler) Start(ctx context.Context) error {
	s.startMu.Lock()
	defer s.startMu.Unlock()

	if s.started {
		return nil
	}

	if err := s.CheckNow(ctx); err != nil {
		return fmt.Errorf("check reminders on scheduler start: %w", err)
	}

	if _, err := s.cron.AddFunc("@every 30m", func() {
		if err := s.CheckNow(ctx); err != nil {
			log.Printf("check due reminders: %v", err)
		}
	}); err != nil {
		return fmt.Errorf("schedule reminder check: %w", err)
	}

	s.cron.Start()
	s.started = true
	return nil
}

func (s *ReminderScheduler) Stop() {
	s.startMu.Lock()
	defer s.startMu.Unlock()

	if !s.started {
		return
	}

	s.cron.Stop()
	s.started = false
}

func (s *ReminderScheduler) CheckNow(ctx context.Context) error {
	s.checkMu.Lock()
	defer s.checkMu.Unlock()

	dueReminders, err := s.provider.GetDueReminders(ctx, s.now())
	if err != nil {
		return fmt.Errorf("get due reminders: %w", err)
	}

	newEvents := make([]domain.DueReminder, 0)
	currentKeys := make(map[string]struct{}, len(dueReminders))

	s.stateMu.Lock()
	for _, reminder := range dueReminders {
		key := dueReminderKey(reminder)
		currentKeys[key] = struct{}{}

		if _, alreadyEmitted := s.emitted[key]; !alreadyEmitted {
			newEvents = append(newEvents, reminder)
		}
	}
	s.emitted = currentKeys
	s.stateMu.Unlock()

	for _, reminder := range newEvents {
		s.emit(ctx, reminderDueEvent, reminder)
	}

	return nil
}

func dueReminderKey(reminder domain.DueReminder) string {
	nextDate := ""
	if reminder.Reminder.NextDueDate != nil {
		nextDate = reminder.Reminder.NextDueDate.UTC().Format(time.RFC3339Nano)
	}

	nextOdometer := ""
	if reminder.Reminder.NextDueOdometer != nil {
		nextOdometer = strconv.FormatInt(*reminder.Reminder.NextDueOdometer, 10)
	}

	return strconv.FormatInt(reminder.Reminder.ID, 10) +
		"|" + strconv.FormatBool(reminder.DueByDate) +
		"|" + strconv.FormatBool(reminder.DueByOdometer) +
		"|" + nextDate +
		"|" + nextOdometer
}
