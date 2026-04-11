package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/redis/go-redis/v9"
)

const channel = "pulseway:checks"

// Client represents one open browser connection
type Client struct {
	id     string
	events chan []byte
}

// Hub manages all connected browser clients
type Hub struct {
	clients map[string]*Client
	mu      sync.RWMutex
	redis   *redis.Client
}

func NewHub(redisURL string) (*Hub, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Hub{
		clients: make(map[string]*Client),
		redis:   client,
	}, nil
}

// Subscribe listens to Redis and broadcasts to all connected browsers
func (h *Hub) Subscribe(ctx context.Context) {
	pubsub := h.redis.Subscribe(ctx, channel)
	defer pubsub.Close()

	log.Println("SSE hub subscribed to Redis")

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-pubsub.Channel():
			h.broadcast([]byte(msg.Payload))
		}
	}
}

// Publish sends a check result to Redis
func (h *Hub) Publish(ctx context.Context, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return h.redis.Publish(ctx, channel, data).Err()
}

// broadcast sends data to all connected browser clients
func (h *Hub) broadcast(data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		select {
		case client.events <- data:
		default:
			// client too slow, skip
		}
	}
}

// ServeHTTP handles the SSE connection from browser
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// SSE requires these headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create client
	clientID := r.RemoteAddr
	client := &Client{
		id:     clientID,
		events: make(chan []byte, 10),
	}

	// Register client
	h.mu.Lock()
	h.clients[clientID] = client
	h.mu.Unlock()

	log.Printf("SSE: client connected %s", clientID)

	// Unregister when browser disconnects
	defer func() {
		h.mu.Lock()
		delete(h.clients, clientID)
		h.mu.Unlock()
		log.Printf("SSE: client disconnected %s", clientID)
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send initial connected message
	fmt.Fprintf(w, "data: {\"status\":\"connected\"}\n\n")
	flusher.Flush()

	// Keep connection open, send events as they arrive
	for {
		select {
		case <-r.Context().Done():
			return
		case data := <-client.events:
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
