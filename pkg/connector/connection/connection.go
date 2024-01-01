package connection

import (
	"github.com/google/uuid"
)

type Connection interface {
	Equals(other Connection) bool
	UUID() uuid.UUID
	Metadata() map[string]any
	WithMetadata(key string, value any)
	GetMetadata(key string) any
	DeleteMetadata(key string)
	MessageChan() chan []byte
	CloseChan() chan struct{}
	Opened() bool
	Closed() bool
	Open()
	Close()
	SendEvent(eventType int, data any) error
}
