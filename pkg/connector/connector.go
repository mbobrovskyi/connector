package connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mbobrovskyi/connector/pkg/connector/connection"
	"github.com/mbobrovskyi/connector/pkg/connector/event"
	"github.com/mbobrovskyi/connector/pkg/logger"
	"sync"
	"time"
)

var AlreadyStartedError = errors.New("connector already started")

type Connector[T connection.Connection] interface {
	Start(ctx context.Context) error
	AddConnection(conn T)
	GetConnections() []T
}

type connectorImpl[T connection.Connection] struct {
	mtx sync.RWMutex

	log         logger.Logger
	connections []T

	eventHandler EventHandler[T]
	errorHandler ErrorHandler

	cleanInterval time.Duration

	isStarted bool
}

func (c *connectorImpl[T]) Start(ctx context.Context) error {
	if c.isStarted {
		return AlreadyStartedError
	}

	c.isStarted = true
	defer func() {
		c.isStarted = false
	}()

	for {
		select {
		case <-ctx.Done():
			c.closeAll()
			return nil
		case <-time.After(c.cleanInterval):
			c.clean()
		}
	}
}

func (c *connectorImpl[T]) closeAll() {
	c.log.Debug("Closing all connections...")

	c.mtx.Lock()
	defer c.mtx.Unlock()

	for _, conn := range c.connections {
		conn.Close()
	}
}

func (c *connectorImpl[T]) clean() {
	c.log.Debug("Cleaning closed connections...")

	c.mtx.Lock()
	defer c.mtx.Unlock()

	connections := make([]T, 0)

	for _, conn := range c.connections {
		if !conn.Closed() {
			connections = append(connections, conn)
		}
	}

	c.connections = connections
}

func (c *connectorImpl[T]) AddConnection(conn T) {
	c.log.Debug(fmt.Sprintf("Added connection %s.", conn.UUID().String()))

	conn.Open()
	c.addConnection(conn)
	go c.listen(conn)
}

func (c *connectorImpl[T]) addConnection(conn T) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.connections = append(c.connections, conn)
}

func (c *connectorImpl[T]) listen(conn T) {
	for {
		select {
		case <-conn.CloseChan():
			c.log.Debug(fmt.Sprintf("Connection %s closed.", conn.UUID()))
			return
		case msg := <-conn.MessageChan():
			c.onMessage(conn, msg)
		}
	}
}

func (c *connectorImpl[T]) onMessage(conn T, data []byte) {
	var rawEvent event.Event

	if err := json.Unmarshal(data, &rawEvent); err != nil {
		c.log.Debug(fmt.Sprintf("Error on parse raw event: %s", err.Error()))
		return
	}

	if err := c.eventHandler.Handle(conn, rawEvent.Type, rawEvent.Data); err != nil {
		if c.errorHandler != nil {
			c.errorHandler.Handle(err)
		} else {
			c.log.Error(err)
		}
	}
}

func (c *connectorImpl[T]) GetConnections() []T {
	return c.connections
}
