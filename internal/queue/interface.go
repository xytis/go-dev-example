package queue

import "errors"

var EOQueue = errors.New("empty queue")

type SimpleQueue interface {
	Push(msg Message) error
	Pull() (Message, error)
}

type ReactiveQueue interface {
	Push(msg Message) error
	Subscribe() (<-chan Message, error)
}
