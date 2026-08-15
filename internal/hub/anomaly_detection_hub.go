package hub

import (
	"sync"
)

type AnomalyDetectionHub struct {
	mu          sync.RWMutex
	subscribers map[chan string]bool
}

// subscriber is a hashmap with key = channel with values of type string and
// value = bool (true indicates that the channel/subscriber is active)

func NewAnomalyDetectionHub() *AnomalyDetectionHub {
	return &AnomalyDetectionHub{
		subscribers: make(map[chan string]bool),
	}
}

func (h *AnomalyDetectionHub) Subscribe() chan string {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan string, 100) //! temp: buffer capacity is set to 100
	h.subscribers[ch] = true

	return ch
}

func (h *AnomalyDetectionHub) Unsubscribe(ch chan string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.subscribers[ch]; ok {
		delete(h.subscribers, ch)
		close(ch)
	}
}

func (h *AnomalyDetectionHub) Broadcast(msg string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for ch := range h.subscribers {
		select {
		case ch <- msg:
		default:
			// channel buffer full, skip to avoid blocking the worker loop
		}
	}
}
