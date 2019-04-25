package hub

import "sync"

// NewMemoryHub creates a new message hub which uses a
// simple in-memory transport.
func NewMemoryHub() Hub {
	return &memoryHub{
		subscribers: []Subscriber{},
	}
}

type memoryHub struct {
	subscribers []Subscriber

	l sync.Mutex
}

func (h *memoryHub) Notify(msg interface{}) {
	h.l.Lock()
	defer h.l.Unlock()

	for _, sub := range h.subscribers {
		go sub.Receive(msg)
	}
}

func (h *memoryHub) Subscribe(sub Subscriber) {
	h.l.Lock()
	defer h.l.Unlock()

	for _, s := range h.subscribers {
		if s == sub {
			return
		}
	}

	h.subscribers = append(h.subscribers, sub)
}

func (h *memoryHub) Unsubscribe(sub Subscriber) {
	h.l.Lock()
	defer h.l.Unlock()

	for i, s := range h.subscribers {
		if s == sub {
			h.subscribers = append(h.subscribers[:i], h.subscribers[i+1:]...)
			return
		}
	}
}
