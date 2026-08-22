package hub

import (
	"sync"
)

type AnomalyDetectionHub struct {
	mu sync.RWMutex
	//sync.RWMutex allows multiple concurrent readers or a single writer at any given time.
	subscribers map[chan string]bool
	// map[T]bool is a standard golang idiom for building a Set of unique items.
	// this is because golang does not have a native, built-in Set data type
	//* so subscribers here is just a Set of channels
}

func NewAnomalyDetectionHub() *AnomalyDetectionHub {
	return &AnomalyDetectionHub{
		subscribers: make(map[chan string]bool),
		// in golang maps must be initialized before use, or writing to them causes a runtime panic.
		// Analogy: we are just creating a Set here called subscribers
	}
}

func (h *AnomalyDetectionHub) Subscribe() chan string {
	h.mu.Lock()
	// Lock() = no write and no read
	//* i.e blocks both read and write
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
	// RLock() =  read lock
	// multiple concurrent calls to Broadcast can read from h.subscribers simultaneously, 
	// but no writer can alter the subscribers list during execution.
	//* i.e blocks only write
	
	for ch := range h.subscribers {
		select {
		case ch <- msg:
		default:
			// channel buffer full, skip to avoid blocking the worker loop
		}
	}
}
