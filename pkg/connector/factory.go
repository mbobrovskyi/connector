package connector

import (
	"github.com/mbobrovskyi/connector/pkg/connector/connection"
	"github.com/mbobrovskyi/connector/pkg/logger"
	"time"
)

type Config struct {
	Logger        logger.Logger
	ErrorHandler  ErrorHandler
	CleanInterval time.Duration
}

func NewT[T connection.Connection](eventHandler EventHandler[T], configs ...Config) Connector[T] {
	conn := &connectorImpl[T]{
		log:           logger.NewNopLogger(),
		eventHandler:  eventHandler,
		cleanInterval: time.Minute,
	}

	for _, config := range configs {
		if config.CleanInterval > 0 {
			conn.cleanInterval = config.CleanInterval
		}

		if config.Logger != nil {
			conn.log = config.Logger
		}

		if config.ErrorHandler != nil {
			conn.errorHandler = config.ErrorHandler
		}
	}

	return conn
}

func New(eventHandler EventHandler[connection.Connection], configs ...Config) Connector[connection.Connection] {
	return NewT[connection.Connection](eventHandler, configs...)
}
