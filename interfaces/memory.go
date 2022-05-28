package interfaces

type Mem interface {
	MemRead(addr uint16) *uint8
	MemWrite(addr uint16, val uint8) error
	UnloadBootloader()
	TickAllComponents(uint64)
}
