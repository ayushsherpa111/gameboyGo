package types

type ReadMemFunc func() *uint8

type WriteMemFunc func(uint8) error

type Events func()

type EventType uint8

const (
	EV_EI EventType = iota
	EV_TIMER
)
