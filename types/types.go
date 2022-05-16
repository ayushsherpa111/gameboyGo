package types

type ReadMemFunc func() *uint8

type WriteMemFunc func(uint8) error

type Events func()
