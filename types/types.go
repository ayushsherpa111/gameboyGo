package types

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ReadMemFunc func() *uint8

type WriteMemFunc func(uint8) error

type Events func()

type EventType uint8

const (
	EV_EI EventType = iota
	EV_TIMER
)

type KeyboardEvent struct {
	Key   sdl.Keycode
	State bool
}

func NewKeyboardEvent(keycode sdl.Keycode, state uint8) KeyboardEvent {
	var newState bool
	if state == 1 {
		newState = true
	}
	return KeyboardEvent{
		Key:   keycode,
		State: newState,
	}
}

var KeyMap map[sdl.Keycode]uint8 = map[sdl.Keycode]uint8{
	sdl.K_DOWN: 0x17,
	sdl.K_z:    0x27, // START BTN

	sdl.K_UP: 0x1B,
	sdl.K_x:  0x2B, // SELECT BTN

	sdl.K_LEFT: 0x1D,
	sdl.K_s:    0x2D, // B BTN

	sdl.K_RIGHT: 0x1E,
	sdl.K_a:     0x2E, // A BTN
}
