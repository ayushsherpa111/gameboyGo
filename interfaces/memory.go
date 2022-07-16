package interfaces

type Mem interface {
	MemRead(uint16, uint64) *uint8
	MemWrite(uint16, uint8, uint64) error
	UnloadBootloader()
	TickAllComponents(uint64)
	SetScheduler(Scheduler)
	HandleInput(uint8, bool)
}
