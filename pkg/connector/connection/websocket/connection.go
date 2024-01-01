package websocket

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mbobrovskyi/connector/pkg/connector/connection"
	"github.com/mbobrovskyi/connector/pkg/connector/event"
)

var _ connection.Connection = (*Connection)(nil)

type Conn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteJSON(v interface{}) error
	Close() error
}

type Connection struct {
	conn Conn

	uuid     uuid.UUID
	metadata map[string]any

	messageChan chan []byte
	closeChan   chan struct{}

	opened bool
	closed bool
}

func (c *Connection) Equals(other connection.Connection) bool {
	if other == nil {
		return false
	}

	return c.UUID() == other.UUID()
}

func (c *Connection) UUID() uuid.UUID {
	return c.uuid
}

func (c *Connection) Metadata() map[string]any {
	return c.metadata
}

func (c *Connection) WithMetadata(key string, value any) {
	c.metadata[key] = value
}

func (c *Connection) GetMetadata(key string) any {
	return c.metadata[key]
}

func (c *Connection) DeleteMetadata(key string) {
	delete(c.metadata, key)
}

func (c *Connection) MessageChan() chan []byte {
	return c.messageChan
}

func (c *Connection) CloseChan() chan struct{} {
	return c.closeChan
}

func (c *Connection) Opened() bool {
	return c.opened
}

func (c *Connection) Closed() bool {
	return c.closed
}

func (c *Connection) Open() {
	go c.open()
}

func (c *Connection) open() {
	if c.opened || c.closed {
		return
	}

	c.opened = true

	defer func() {
		c.opened = false
	}()

	for {
		select {
		case <-c.closeChan:
			return
		default:
			_, msgData, err := c.conn.ReadMessage()
			if err != nil {
				c.Close()
				return
			}
			c.messageChan <- msgData
		}
	}
}

func (c *Connection) Close() {
	if !c.closed {
		_ = c.conn.Close()
		close(c.closeChan)
		close(c.messageChan)
	}

	c.closed = true
}

func (c *Connection) SendEvent(eventType int, data any) error {
	if c.closed {
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := c.conn.WriteJSON(event.New(eventType, jsonData)); err != nil {
		return err
	}

	return nil
}

func New(conn Conn) *Connection {
	return &Connection{
		conn:        conn,
		uuid:        uuid.New(),
		metadata:    make(map[string]any),
		messageChan: make(chan []byte),
		closeChan:   make(chan struct{}),
	}
}

func NewWithMetadata(conn Conn, metadata map[string]any) *Connection {
	return &Connection{
		conn:        conn,
		uuid:        uuid.New(),
		metadata:    metadata,
		messageChan: make(chan []byte),
		closeChan:   make(chan struct{}),
	}
}
