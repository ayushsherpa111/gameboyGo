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

func NewScheduler(c *cpu.CPU) scheduler {
	return scheduler{
		eventQueue: make([]event, 0, MAX_EVENT_QUEUE),
		refCPU:     c,
	}
}

func newEvent(evFunc types.Events, cycles uint64) event {
	return event{
		cycles,
		evFunc,
	}
}

func (s scheduler) ScheduleEvent(event types.Events, cycles uint64) error {
	if len(s.eventQueue) == MAX_EVENT_QUEUE {
		return errors.New("Event queue full")
	}
	evQ := newEvent(event, s.refCPU.CycleCount+cycles)

	if len(s.eventQueue) == 0 {
		s.eventQueue = append(s.eventQueue, evQ)
	} else {
		var i int
		for i = 0; i < len(s.eventQueue) && s.eventQueue[i].cycles > evQ.cycles; i++ {
		}

		s.eventQueue = append(s.eventQueue, evQ)
		copy(s.eventQueue[i+1:], s.eventQueue[i:])
		s.eventQueue[i] = evQ
	}
	return nil
}

func (s scheduler) Tick() {
	// loop through the queue and check if
	currCycle := s.refCPU.CycleCount
	for {
		if len(s.eventQueue) > 0 && s.eventQueue[0].cycles <= currCycle {
			s.eventQueue[0].evFnc()
			if len(s.eventQueue) >= 1 {
				s.eventQueue = append([]event(nil), s.eventQueue[1:]...)
			} else {
				s.eventQueue = make([]event, 0, MAX_EVENT_QUEUE)
			}
		} else {
			break
		}
	}
}
