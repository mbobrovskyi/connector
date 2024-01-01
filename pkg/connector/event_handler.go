package connector

import "github.com/mbobrovskyi/connector/pkg/connector/connection"

type EventHandler[T connection.Connection] interface {
	Handle(conn T, eventType int, data []byte) error
}
