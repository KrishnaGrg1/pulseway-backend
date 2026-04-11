package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/KrishnaGrg1/pulseway/internal/queue"
	"github.com/KrishnaGrg1/pulseway/internal/store"
)

type Scheduler struct {
	store *store.Store
	queue *queue.Queue
}

func New(s *store.Store, q *queue.Queue) *Scheduler {
	return &Scheduler{store: s, queue: q}
}

func (s *Scheduler) Start(ctx context.Context) {
	log.Println("Scheduler started")

	// Check every 5 seconds which monitors are due
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return
		case <-ticker.C:
			s.scheduleChecks(ctx)
		}
	}
}

func (s *Scheduler) scheduleChecks(ctx context.Context) {
	// Get all active monitors
	monitors, err := s.store.Queries.ListAllActiveMonitors(ctx)
	if err != nil {
		log.Println("Scheduler: failed to list monitors:", err)
		return
	}

	for _, m := range monitors {
		job := queue.CheckJob{
			MonitorID: m.ID,
			URL:       m.Url,
		}
		if err := s.queue.Publish(ctx, job); err != nil {
			log.Printf("Scheduler: failed to publish job for monitor %d: %v", m.ID, err)
		}
	}

	if len(monitors) > 0 {
		log.Printf("Scheduler: queued %d checks", len(monitors))
	}
}
