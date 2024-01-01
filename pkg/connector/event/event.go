package event

import "encoding/json"

type Event struct {
	Type int             `json:"type"`
	Data json.RawMessage `json:"data"`
}

func New(eventType int, data json.RawMessage) Event {
	return Event{Type: eventType, Data: data}
}
