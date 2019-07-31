package gopubsub

import "sync"

type Any interface{}

type Subscriber = func(Any)

type Hub struct {
	mutex       sync.Mutex
	subscribers []Subscriber
	channel     chan Any
}

func NewHub() *Hub {
	hub := Hub{
		subscribers: make([]Subscriber, 0),
		channel:     make(chan Any),
	}
	go hub.fanOut()
	return &hub
}

func (h *Hub) Subscribe(subscriber Subscriber) {
	h.mutex.Lock()

	h.subscribers = append(h.subscribers, subscriber)

	h.mutex.Unlock()
}

func (h *Hub) Publish(message Any) {
	// non blocking write
	select {
	case h.channel <- message:
	default:
	}
}

func (h *Hub) fanOut() {
	for {
		message := <-h.channel
		for _, subscriber := range h.subscribers {
			subscriber(message)
		}
	}
}
