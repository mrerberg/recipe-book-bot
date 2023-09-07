package events

import "context"

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(ctx context.Context, event Event) error
}

const (
	Unknown EventType = iota
	Message
	CallBack
)

type EventType int

type Event struct {
	Type     EventType
	Text     string
	Metadata interface{}
}
