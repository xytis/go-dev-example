package queue

type Message string

// ParseMessage is meant to simulate a domain translation
//
//	from transport type into domain type.
//	Obviously, it now is a dummy method.
func ParseMessage(message string) (Message, error) {
	return Message(message), nil
}
