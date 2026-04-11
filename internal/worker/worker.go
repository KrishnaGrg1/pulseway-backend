package worker

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	db "github.com/KrishnaGrg1/pulseway/internal/db/sqlc"
	"github.com/KrishnaGrg1/pulseway/internal/queue"
	"github.com/KrishnaGrg1/pulseway/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	store      *store.Store
	queue      *queue.Queue
	httpClient *http.Client
}

func New(s *store.Store, q *queue.Queue) *Worker {
	return &Worker{
		store: s,
		queue: q,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (w *Worker) Start(ctx context.Context, workerCount int) {
	deliveries, err := w.queue.Consume()
	if err != nil {
		log.Fatal("Worker: failed to start consuming:", err)
	}

	log.Printf("Starting %d workers", workerCount)

	for i := 0; i < workerCount; i++ {
		go w.run(ctx, deliveries)
	}

	<-ctx.Done()
	log.Println("Workers stopped")
}

func (w *Worker) run(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-deliveries:
			if !ok {
				return
			}
			w.processJob(ctx, d)
		}
	}
}

func (w *Worker) processJob(ctx context.Context, d amqp.Delivery) {
	var job queue.CheckJob
	if err := json.Unmarshal(d.Body, &job); err != nil {
		log.Println("Worker: failed to parse job:", err)
		d.Nack(false, false) // reject, don't requeue
		return
	}

	log.Printf("Worker: checking monitor %d → %s", job.MonitorID, job.URL)

	// Record start time for latency
	start := time.Now()

	// Make HTTP request
	resp, err := w.httpClient.Get(job.URL)
	latency := int32(time.Since(start).Milliseconds())

	var status string
	var statusCode pgtype.Int4

	if err != nil {
		status = "down"
		statusCode = pgtype.Int4{Valid: false}
	} else {
		resp.Body.Close()
		statusCode = pgtype.Int4{Int32: int32(resp.StatusCode), Valid: true}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			status = "up"
		} else {
			status = "down"
		}
	}

	// Save check result
	_, err = w.store.Queries.CreateCheckResult(ctx, db.CreateCheckResultParams{
		MonitorID:  job.MonitorID,
		Status:     status,
		LatencyMs:  latency,
		StatusCode: statusCode,
	})
	if err != nil {
		log.Printf("Worker: failed to save result: %v", err)
		d.Nack(false, true) // requeue
		return
	}

	// Handle incident detection
	if status == "down" {
		w.handleDown(ctx, job.MonitorID)
	} else {
		w.handleUp(ctx, job.MonitorID)
	}

	log.Printf("Worker: monitor %d is %s (%dms)", job.MonitorID, status, latency)

	// Acknowledge — tell RabbitMQ job is done
	d.Ack(false)
}

func (w *Worker) handleDown(ctx context.Context, monitorID int64) {
	// Check if incident already exists
	_, err := w.store.Queries.GetActiveIncident(ctx, monitorID)
	if err == nil {
		return // incident already open
	}

	// Create new incident
	_, err = w.store.Queries.CreateIncident(ctx, monitorID)
	if err != nil {
		log.Printf("Worker: failed to create incident: %v", err)
		return
	}

	log.Printf("Worker: incident created for monitor %d", monitorID)
}

func (w *Worker) handleUp(ctx context.Context, monitorID int64) {
	// Resolve any open incident
	_, err := w.store.Queries.ResolveIncident(ctx, monitorID)
	if err != nil {
		return // no open incident, that's fine
	}

	log.Printf("Worker: incident resolved for monitor %d", monitorID)
}
