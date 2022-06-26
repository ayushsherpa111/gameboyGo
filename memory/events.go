package memory

import "github.com/ayushsherpa111/gameboyEMU/types"

func (m *memory) scheduleTimerEvents(timerVal uint8, cycleCount uint64) {
	tac_addr := TAC - IO_START
	if (m.ioRegs[tac_addr] & TIMER_ENABLE) == TIMER_ENABLE {
		// schedule events
		var timerCycle uint64
		switch m.ioRegs[tac_addr] & 0x3 {
		case 0x00:
			timerCycle = (0xFF - uint64(timerVal)) * 1024
		case 0x01:
			timerCycle = (0xFF - uint64(timerVal)) * 16
		case 0x02:
			timerCycle = (0xFF - uint64(timerVal)) * 64
		case 0x03:
			timerCycle = (0xFF - uint64(timerVal)) * 256
		}
		m.Scheduler.ScheduleEvent(m.SetIFTimer(cycleCount), timerCycle, types.EV_TIMER)
	}
}

func (m *memory) getClockTiming() uint64 {
	switch m.ioRegs[TAC-IO_START] & 0x3 {
	case 0x00:
		return 1024
	case 0x01:
		return 16
	case 0x02:
		return 64
	default:
		return 256
	}
}
