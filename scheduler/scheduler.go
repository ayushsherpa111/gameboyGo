package scheduler

import (
	"errors"

	"github.com/ayushsherpa111/gameboyEMU/cpu"
	"github.com/ayushsherpa111/gameboyEMU/types"
)

const MAX_EVENT_QUEUE = 256

type event struct {
	cycles uint64
	evFnc  types.Events
}

type scheduler struct {
	eventQueue []event
	refCPU     *cpu.CPU
}

func NewScheduler() scheduler {
	return scheduler{
		eventQueue: make([]event, 0, MAX_EVENT_QUEUE),
	}
}

func newEvent(evFunc types.Events, cycles uint64) event {
	return event{
		cycles,
		evFunc,
	}
}

func (s scheduler) ScheduleEvent(event types.Events, cycles uint64) error {
	if len(s.eventQueue) == cap(s.eventQueue) {
		return errors.New("Event queue full")
	}
	s.eventQueue = append(s.eventQueue, newEvent(event, cycles))
	return nil
}

func (s scheduler) Tick() {}
