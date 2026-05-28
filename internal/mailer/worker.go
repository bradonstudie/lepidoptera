package mailer

import (
	"context"
	"log"
	"time"

	db "github.com/bradonstudie/lepidoptera/internal/db/generated"
)

type Worker struct {
	mailer  Mailer
	queries *db.Queries
	secret  string
}

func NewWorker(mailer Mailer, queries *db.Queries, secret string) *Worker {
	return &Worker{mailer: mailer, queries: queries, secret: secret}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// run once immediately on startup
	w.run(ctx)

	for {
		select {
		case <-ticker.C:
			w.run(ctx)
		case <-ctx.Done():
			log.Println("worker: shutting down")
			return
		}
	}
}

func (w *Worker) run(ctx context.Context) {
	if err := w.sendPendingAnnouncements(ctx); err != nil {
		log.Printf("worker: announcements error: %v", err)
	}
	if err := w.sendPendingReminders(ctx); err != nil {
		log.Printf("worker: reminders error: %v", err)
	}

	if err := w.queries.DeleteExpiredSessions(ctx); err != nil {
		log.Printf("worker: session cleanup error: %v", err)
	}
	if err := w.queries.DeleteUsedLoginTokens(ctx); err != nil {
		log.Printf("worker: login token cleanup error: %v", err)
	}
}

const notificationBatchSize = 500

func (w *Worker) sendPendingAnnouncements(ctx context.Context) error {
	for {
		pending, err := w.queries.GetPendingAnnouncementNotifications(ctx, notificationBatchSize)
		if err != nil {
			return err
		}
		if len(pending) == 0 {
			return nil
		}

		log.Printf("worker: sending %d announcement(s)", len(pending))

		for _, n := range pending {
			if err := w.sendAnnouncement(ctx, n); err != nil {
				log.Printf("worker: failed announcement to %s: %v", n.SubscriberEmail, err)
				continue
			}
			if err := w.queries.MarkNotificationSent(ctx, n.ID); err != nil {
				log.Printf("worker: failed to mark notification sent %s: %v", n.ID, err)
			}
		}

		if len(pending) < notificationBatchSize {
			return nil
		}
	}
}

func (w *Worker) sendPendingReminders(ctx context.Context) error {
	for {
		pending, err := w.queries.GetPendingReminderNotifications(ctx, notificationBatchSize)
		if err != nil {
			return err
		}
		if len(pending) == 0 {
			return nil
		}

		log.Printf("worker: sending %d reminder(s)", len(pending))

		for _, n := range pending {
			if err := w.sendReminder(ctx, n); err != nil {
				log.Printf("worker: failed reminder to %s: %v", n.SubscriberEmail, err)
				continue
			}
			if err := w.queries.MarkNotificationSent(ctx, n.ID); err != nil {
				log.Printf("worker: failed to mark notification sent %s: %v", n.ID, err)
			}
		}

		if len(pending) < notificationBatchSize {
			return nil
		}
	}
}

func (w *Worker) sendAnnouncement(ctx context.Context, n db.GetPendingAnnouncementNotificationsRow) error {
	unsubToken := UnsubscribeToken(n.SubscriberID.UUID, w.secret)

	email := ShowAnnouncedEmail(ShowEmailData{
		Title:       n.Title,
		Venue:       n.VenueName.String,
		Date:        n.Date,
		Slug:        n.Slug,
		Description: n.Description.String,
		TicketURL:   n.TicketUrl.String,
	}, unsubToken)
	email.To = n.SubscriberEmail

	return w.mailer.Send(ctx, email)
}

func (w *Worker) sendReminder(ctx context.Context, n db.GetPendingReminderNotificationsRow) error {
	unsubToken := UnsubscribeToken(n.SubscriberID.UUID, w.secret)

	email := ReminderEmail(ShowEmailData{
		Title:     n.Title,
		Venue:     n.VenueName.String,
		Date:      n.Date,
		Slug:      n.Slug,
		TicketURL: n.TicketUrl.String,
	}, unsubToken)
	email.To = n.SubscriberEmail

	return w.mailer.Send(ctx, email)
}
