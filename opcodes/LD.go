package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type ld struct {
	label string
	op    byte
	c     *cpu.CPU
}

func (i *ld) r16_u16(r1, r2 uint8, val uint16) {
	i.c.SetRegister(r1, cpu.Register(val&0x00FF))
	i.c.SetRegister(r2, cpu.Register((val&0xFF00)>>8))
}

func (i *ld) u16_u8(addr uint16, v uint8) {
	i.c.SetMem(addr, v)
}

func (i *ld) r8_u8(reg uint8, val uint8) {
	i.c.SetRegister(reg, cpu.Register(val))
}

func (i *ld) u16_SP(mem uint16) {
	// TODO verify if uint8(uint16) = 0x00FF
	i.c.SetMem(mem, uint8(i.c.SP))
	i.c.SetMem(mem+1, uint8(i.c.SP>>8))
}

func (i *ld) Exec() {
	switch i.op {
	case 0x01:
		// LD BC, u16
		i.r16_u16(cpu.B, cpu.C, i.c.Fetch16())
	case 0x02:
		// LD (BC), A
		i.u16_u8(i.c.BC(), uint8(i.c.GetRegister(cpu.A)))
	case 0x06:
		// LD B, u8
		i.r8_u8(cpu.B, i.c.Fetch())
	case 0x008:
		// LD (u16), SP
		i.u16_SP(i.c.Fetch16())
	default:
		panic("invalid opcode")
	}
}

func (i *ld) Label() string {
	return i.label
}

func NewLD(label string, op byte, cpu *cpu.CPU) *ld {
	return &ld{
		label,
		op,
		cpu,
	}
}
