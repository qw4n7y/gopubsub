package gopubsub

import "sync"

type Subscriber = func(message interface{})

type Hub struct {
	mutex       sync.Mutex
	subscribers []Subscriber
	channel     chan interface{}
}

func NewHub() *Hub {
	hub := Hub{
		subscribers: make([]Subscriber, 0),
		channel:     make(chan interface{}, 100),
	}
	go hub.fanOut()
	return &hub
}

func (h *Hub) Subscribe(subscriber Subscriber) {
	h.mutex.Lock()

	h.subscribers = append(h.subscribers, subscriber)

	h.mutex.Unlock()
}

func (h *Hub) Publish(message interface{}) {
	// non blocking write
	select {
	case h.channel <- message:
	default:
	}
}

func (h *Hub) fanOut() {
	for {
		message := <-h.channel
		go func() {
			for _, subscriber := range h.subscribers {
				subscriber(message)
			}
		}()
	}
}
