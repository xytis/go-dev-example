package queue

import (
	"errors"
	"sync"
)

// unbufferedQueue is meant to be used from different routines
//
//	messages are directly delivered while there are active subscriptions
//	message is dropped if no subscriptions are active
type unbufferedQueue struct {
	mutex sync.Mutex
	subs  []chan Message
}

func NewUnbufferedQueue() ReactiveQueue {
	return &unbufferedQueue{
		mutex: sync.Mutex{},
		subs:  []chan Message{},
	}
}

func (m *unbufferedQueue) Subscribe() (<-chan Message, error) {
	sub := make(chan Message)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.subs = append(m.subs, sub)
	return sub, nil
}

func (m *unbufferedQueue) Push(msg Message) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if len(m.subs) == 0 {
		return errors.New("dropping message, no active subscriptions")
	}
	for _, sub := range m.subs {
		sub <- msg
	}
	return nil
}

// arrayQueue is meant to be used from different routines
//
//	messages are buffered
//	only a single consumer can use this queue (interface limitation)
//	this implementation intentionally leaks memory in the long run
type arrayQueue struct {
	mutex sync.Mutex
	queue []Message
}

func NewArrayQueue() SimpleQueue {
	return &arrayQueue{
		mutex: sync.Mutex{},
		queue: make([]Message, 0),
	}
}

func (m *arrayQueue) Pull() (Message, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if len(m.queue) == 0 {
		// Note: This is used to avoid changing the code if Message type changes.
		var zero Message
		return zero, EOQueue
	}
	msg := m.queue[0]
	// Intentional memory leak.
	m.queue = m.queue[1:]
	return msg, nil
}

func (m *arrayQueue) Push(msg Message) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// Intentional memory leak.
	m.queue = append(m.queue, msg)
	return nil
}

// Update:
//  regarding the memory leak comments:
//  My initial source of this pattern was here: https://go.dev/blog/slices-intro
//  I've also found similar newer sources (https://go101.org/article/memory-leaking.html)
//
//  However.
//  I wrote few tests, and I was unable to confirm any of the mentioned cases in isolation.
//  Since I do not have any production code that uses this pattern, I am unable to neither confirm
//  neither debunk this case.
//
// P.S. Just in case I tested the same on go 1.18. No change to the observation.
