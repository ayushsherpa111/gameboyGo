package interfaces

import "github.com/ayushsherpa111/gameboyEMU/types"

type Scheduler interface {
	Tick()
	ScheduleEvent(types.Events, uint64) error
}
