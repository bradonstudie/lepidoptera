package mailer

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Worker struct {
	mailer Mailer
	db     *pgxpool.Pool
	secret string
}

func NewWorker(m Mailer, db *pgxpool.Pool, secret string) *Worker {
	return &Worker{mailer: m, db: db, secret: secret}
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
}

func (w *Worker) sendPendingAnnouncements(ctx context.Context) error {
	// TODO: implement with sqlc generated queries
	// 1. fetch pending notifications of type 'published'
	// 2. for each, send email via w.mailer
	// 3. mark notification sent_at
	return nil
}

func (w *Worker) sendPendingReminders(ctx context.Context) error {
	// TODO: implement with sqlc generated queries
	// 1. fetch shows happening tomorrow with unsent reminder notifications
	// 2. for each subscriber, send reminder email
	// 3. mark notification sent_at
	return nil
}
