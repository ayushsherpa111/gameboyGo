package interfaces

type GPU interface {
	UpdateGPU()
	Read_VRAM(addr uint16) *uint8
	Write_VRAM(uint16, uint8)
	Read_OAM(addr uint16) *uint8
	Write_OAM(uint16, uint8, bool)
	Read_Regs(uint16) *uint8
	Write_Regs(uint16, uint8) error
	PrintDetails()
	RefInterruptFlag(*uint8)
}
