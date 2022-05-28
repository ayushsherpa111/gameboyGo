package memory

func (m *memory) scheduleTimerEvents(timerVal uint8) {
	tac_addr := TAC - IO_START
	if (m.ioRegs[tac_addr] & TIMER_ENABLE) == TIMER_ENABLE {
		// schedule events
		switch m.ioRegs[tac_addr] & 0x3 {
		case 0x00:
			m.Scheduler.ScheduleEvent(m.SetIFTimer(), (0xFF-uint64(timerVal))*1024)
		case 0x01:
			m.Scheduler.ScheduleEvent(m.SetIFTimer(), (0xFF-uint64(timerVal))*16)
		case 0x02:
			m.Scheduler.ScheduleEvent(m.SetIFTimer(), (0xFF-uint64(timerVal))*64)
		case 0x03:
			m.Scheduler.ScheduleEvent(m.SetIFTimer(), (0xFF-uint64(timerVal))*256)
		}
	}
}
