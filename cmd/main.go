package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KrishnaGrg1/pulseway/internal/api"
	"github.com/KrishnaGrg1/pulseway/internal/config"
	"github.com/KrishnaGrg1/pulseway/internal/queue"
	"github.com/KrishnaGrg1/pulseway/internal/scheduler"
	"github.com/KrishnaGrg1/pulseway/internal/store"
	"github.com/KrishnaGrg1/pulseway/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DB_URL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Cannot ping database:", err)
	}
	log.Println("Connected to Neon database")

	q, err := queue.New(cfg.RABBITMQ_URL)
	if err != nil {
		log.Fatal("Cannot connect to RabbitMQ:", err)
	}
	defer q.Close()
	log.Println("Connected to RabbitMQ")

	s := store.New(pool)

	// Start scheduler
	ctx, cancel := context.WithCancel(context.Background())
	sched := scheduler.New(s, q)
	go sched.Start(ctx)

	// Start workers
	w := worker.New(s, q)
	go w.Start(ctx, cfg.WORKER_COUNT)

	router := api.NewRouter(s, cfg)

	srv := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: router,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down...")
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx)
	}()

	log.Printf("Pulseway running on port %s", cfg.PORT)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
