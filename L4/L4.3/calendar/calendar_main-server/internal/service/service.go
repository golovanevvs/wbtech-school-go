package service

import (
	"context"
	"sync"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgTelegram"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/transport/trhttp/handler/telegramHandler"
	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/zlog"
)

// CalendarService handles calendar events business logic
type CalendarService struct {
	repository *repository.Repository
	noticeSvc  *NoticeService
	lg         *zlog.Zerolog
	cfg        *Config

	// Channels for background workers
	reminderChan chan *model.Event
	logChan      chan LogEntry

	// Background workers
	wg sync.WaitGroup
}

// LogEntry represents a log entry for async logger
type LogEntry struct {
	Level   string
	Message string
	Data    map[string]interface{}
	Time    time.Time
}

// Service structure that combines all services
type Service struct {
	Calendar *CalendarService
	Notice   *NoticeService
}

// New creates a new Service structure
func New(
	cfg *Config,
	rp *repository.Repository,
	tgClient *pkgTelegram.Client,
	emailClient *pkgEmail.Client,
	rs *pkgRetry.Retry,
) *Service {
	lg := zlog.Logger.With().Str("layer", "service").Logger()

	noticeService := &NoticeService{
		lg: &lg,
		rs: rs,
		tg: tgClient,
		em: emailClient,
	}

	calendarService := &CalendarService{
		repository:   rp,
		noticeSvc:    noticeService,
		lg:           &lg,
		cfg:          cfg,
		reminderChan: make(chan *model.Event, 100),
		logChan:      make(chan LogEntry, 1000),
	}

	return &Service{
		Calendar: calendarService,
		Notice:   noticeService,
	}
}

// TelegramService returns the telegram service (for compatibility)
func (sv *Service) TelegramService() telegramHandler.ISvForTelegramHandler {
	// Return a simple handler that does nothing for now
	return &TelegramHandler{tg: nil}
}

// Start implements the telegramHandler.ISvForTelegramHandler interface
func (th *TelegramHandler) Start(ctx context.Context, username string, chatID int64, message string) error {
	// Simple implementation that does nothing for now
	return nil
}

// CalendarService returns the calendar service
func (sv *Service) CalendarService() CalendarServiceInterface {
	return sv.Calendar
}

// CalendarServiceInterface defines the interface for calendar service methods
type CalendarServiceInterface interface {
	GetMonthEvents(ctx context.Context, year, month int) ([]model.Event, error)
	CreateEvent(ctx context.Context, eventData *model.CreateEventRequest) (*model.Event, error)
	GetEvent(ctx context.Context, id int) (*model.Event, error)
	UpdateEvent(ctx context.Context, id int, eventData *model.CreateEventRequest) (*model.Event, error)
	DeleteEvent(ctx context.Context, id int) error
	GetDayEvents(ctx context.Context, date string) ([]model.Event, error)
	StartBackgroundWorkers(ctx context.Context)
	StopBackgroundWorkers()
}

// StartBackgroundWorkers starts all background workers for CalendarService
func (cs *CalendarService) StartBackgroundWorkers(ctx context.Context) {
	cs.startReminderWorker(ctx)
	cs.startCleanupWorker(ctx)
	cs.startAsyncLogger(ctx)
}

// StopBackgroundWorkers stops all background workers for CalendarService
func (cs *CalendarService) StopBackgroundWorkers() {
	close(cs.reminderChan)
	close(cs.logChan)
	cs.wg.Wait()
}

// StartBackgroundWorkers starts all background workers
func (sv *Service) StartBackgroundWorkers(ctx context.Context) {
	sv.Calendar.startReminderWorker(ctx)
	sv.Calendar.startCleanupWorker(ctx)
	sv.Calendar.startAsyncLogger(ctx)
}

// StopBackgroundWorkers stops all background workers
func (sv *Service) StopBackgroundWorkers() {
	close(sv.Calendar.reminderChan)
	close(sv.Calendar.logChan)
	sv.Calendar.wg.Wait()
}

// CalendarService methods

// GetMonthEvents returns events for a specific month
func (cs *CalendarService) GetMonthEvents(ctx context.Context, year, month int) ([]model.Event, error) {
	events, err := cs.repository.GetMonthEvents(year, month)
	if err != nil {
		cs.LogError("GetMonthEvents", err, map[string]interface{}{
			"year":  year,
			"month": month,
		})
		return nil, err
	}

	return events, nil
}

// CreateEvent creates a new event and schedules reminder if needed
func (cs *CalendarService) CreateEvent(ctx context.Context, eventData *model.CreateEventRequest) (*model.Event, error) {
	event, err := cs.repository.CreateEvent(eventData)
	if err != nil {
		cs.LogError("CreateEvent", err, map[string]interface{}{
			"title": eventData.Title,
		})
		return nil, err
	}

	// If event has reminder, add it to reminder channel
	if event.Reminder && event.ReminderTime != nil {
		select {
		case cs.reminderChan <- event:
		default:
			cs.LogWarn("CreateEvent", "Reminder channel full, skipping reminder", map[string]interface{}{
				"event_id": event.ID,
			})
		}
	}

	cs.LogInfo("CreateEvent", "Event created successfully", map[string]interface{}{
		"event_id":     event.ID,
		"title":        event.Title,
		"has_reminder": event.Reminder,
	})

	return event, nil
}

// GetEvent returns an event by ID
func (cs *CalendarService) GetEvent(ctx context.Context, id int) (*model.Event, error) {
	event, err := cs.repository.GetEvent(id)
	if err != nil {
		cs.LogError("GetEvent", err, map[string]interface{}{"event_id": id})
		return nil, err
	}

	return event, nil
}

// UpdateEvent updates an existing event
func (cs *CalendarService) UpdateEvent(ctx context.Context, id int, eventData *model.CreateEventRequest) (*model.Event, error) {
	event, err := cs.repository.UpdateEvent(id, eventData)
	if err != nil {
		cs.LogError("UpdateEvent", err, map[string]interface{}{"event_id": id})
		return nil, err
	}

	cs.LogInfo("UpdateEvent", "Event updated successfully", map[string]interface{}{
		"event_id": event.ID,
		"title":    event.Title,
	})

	return event, nil
}

// DeleteEvent deletes an event by ID
func (cs *CalendarService) DeleteEvent(ctx context.Context, id int) error {
	err := cs.repository.DeleteEvent(id)
	if err != nil {
		cs.LogError("DeleteEvent", err, map[string]interface{}{"event_id": id})
		return err
	}

	cs.LogInfo("DeleteEvent", "Event deleted successfully", map[string]interface{}{
		"event_id": id,
	})

	return nil
}

// GetDayEvents returns events for a specific day
func (cs *CalendarService) GetDayEvents(ctx context.Context, date string) ([]model.Event, error) {
	events, err := cs.repository.GetDayEvents(date)
	if err != nil {
		cs.LogError("GetDayEvents", err, map[string]interface{}{"date": date})
		return nil, err
	}

	return events, nil
}

// Logging methods for async logger
func (cs *CalendarService) LogInfo(method, message string, data map[string]interface{}) {
	cs.logChan <- LogEntry{
		Level:   "info",
		Message: message,
		Data:    data,
		Time:    time.Now(),
	}
}

func (cs *CalendarService) LogError(method string, err error, data map[string]interface{}) {
	data["error"] = err.Error()
	cs.logChan <- LogEntry{
		Level:   "error",
		Message: method + " failed",
		Data:    data,
		Time:    time.Now(),
	}
}

func (cs *CalendarService) LogWarn(method, message string, data map[string]interface{}) {
	cs.logChan <- LogEntry{
		Level:   "warn",
		Message: message,
		Data:    data,
		Time:    time.Now(),
	}
}

// Background worker methods

// startReminderWorker starts the reminder worker
func (cs *CalendarService) startReminderWorker(ctx context.Context) {
	cs.wg.Add(1)
	go func() {
		defer cs.wg.Done()
		cs.reminderWorker(ctx)
	}()
}

// startCleanupWorker starts the cleanup worker
func (cs *CalendarService) startCleanupWorker(ctx context.Context) {
	cs.wg.Add(1)
	go func() {
		defer cs.wg.Done()
		cs.cleanupWorker(ctx)
	}()
}

// startAsyncLogger starts the async logger
func (cs *CalendarService) startAsyncLogger(ctx context.Context) {
	cs.wg.Add(1)
	go func() {
		defer cs.wg.Done()
		cs.asyncLogger(ctx)
	}()
}

// reminderWorker processes reminder notifications
func (cs *CalendarService) reminderWorker(ctx context.Context) {
	lg := cs.lg.With().Str("worker", "reminder").Logger()
	lg.Info().Msg("Reminder worker started")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			lg.Info().Msg("Reminder worker stopped")
			return
		case event := <-cs.reminderChan:
			cs.processReminder(lg, event)
		case <-ticker.C:
			cs.checkDueReminders(lg)
		}
	}
}

// processReminder sends a reminder for an event
func (cs *CalendarService) processReminder(lg zerolog.Logger, event *model.Event) {
	// Format reminder message
	// message := fmt.Sprintf("Напоминание: %s\nВремя: %s",
	//     event.Title,
	//     event.Start.Format("2006-01-02 15:04"))

	// For now, we'll just log the reminder
	// In a real implementation, you might send to a specific user or channel
	lg.Info().Int("event_id", event.ID).Str("title", event.Title).Msg("Sending reminder")

	// Here you could integrate with Telegram/Email services
	// cs.noticeSvc.SendNotice(ctx, model.Notice{...})
}

// checkDueReminders checks for reminders that are due
func (cs *CalendarService) checkDueReminders(lg zerolog.Logger) {
	now := time.Now()

	// This is a simplified implementation
	// In a real application, you might query the database for events with due reminders
	lg.Debug().Time("check_time", now).Msg("Checking for due reminders")
}

// cleanupWorker cleans up old events
func (cs *CalendarService) cleanupWorker(ctx context.Context) {
	lg := cs.lg.With().Str("worker", "cleanup").Logger()
	lg.Info().Msg("Cleanup worker started")

	ticker := time.NewTicker(30 * time.Minute) // Run every 30 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			lg.Info().Msg("Cleanup worker stopped")
			return
		case <-ticker.C:
			cs.performCleanup(lg)
		}
	}
}

// performCleanup removes old events
func (cs *CalendarService) performCleanup(lg zerolog.Logger) {
	// Archive events older than 1 year
	cutoffDate := time.Now().AddDate(-1, 0, 0)

	lg.Info().Time("cutoff_date", cutoffDate).Msg("Starting cleanup of old events")

	// This would typically involve:
	// 1. Moving old events to an archive table
	// 2. Or deleting events that are no longer needed
	// For now, we'll just log the action

	lg.Info().Msg("Cleanup completed")
}

// asyncLogger processes log entries asynchronously
func (cs *CalendarService) asyncLogger(ctx context.Context) {
	lg := cs.lg.With().Str("worker", "logger").Logger()
	lg.Info().Msg("Async logger started")

	for {
		select {
		case <-ctx.Done():
			lg.Info().Msg("Async logger stopped")
			return
		case entry := <-cs.logChan:
			cs.writeLog(lg, entry)
		}
	}
}

// writeLog writes a log entry
func (cs *CalendarService) writeLog(lg zerolog.Logger, entry LogEntry) {
	switch entry.Level {
	case "info":
		lg.Info().Str("method", entry.Data["method"].(string)).Msg(entry.Message)
	case "error":
		lg.Error().Err(entry.Data["error"].(error)).Str("method", entry.Data["method"].(string)).Msg(entry.Message)
	case "warn":
		lg.Warn().Str("method", entry.Data["method"].(string)).Msg(entry.Message)
	default:
		lg.Info().Msg(entry.Message)
	}
}
