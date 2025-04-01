package queue

type SimpleQueue interface {
	Push(msg Message) error
	Pull() (Message, error)
}

type ReactiveQueue interface {
	Push(msg Message) error
	Subscribe() (<-chan Message, error)
}
