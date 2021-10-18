package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type inc struct {
	c *cpu.CPU
}

func (i *inc) _INC(val *uint8) {
	i.c.SET_NEG(false)
	i.c.SET_HALF_CARRY(*val&0xF+0x1 > 0x0F)

	(*val)++

	i.c.SET_ZERO(*(val) == 0x00)
}

func (i *inc) inc_r8(reg uint8) {
	i._INC(i.c.GetRegister(reg))
}

func (i *inc) inc_r16(r1, r2 uint8, val uint16) {
	// set flags
	val++
	i.c.SetRegister(r1, uint8(val>>8))
	i.c.SetRegister(r2, uint8(val))
}

func (i *inc) inc_u16(addr uint16) {
	i._INC(i.c.GetMem(addr))
}

func (i *inc) Exec(op byte) {
	switch op {
	case 0x01:
		// INC BC
		i.inc_r16(cpu.B, cpu.C, i.c.BC())
	case 0x04:
		// INC B
		i.inc_r8(cpu.B)
	case 0x0C:
		// INC C
		i.inc_r8(cpu.C)
	case 0x13:
		// INC DE
		i.inc_r16(cpu.D, cpu.E, i.c.DE())
	case 0x14:
		// INC D
		i.inc_r8(cpu.D)
	case 0x1C:
		// INC Edependent
		i.inc_r8(cpu.E)
	case 0x23:
		// INC HL
		i.inc_r16(cpu.H, cpu.L, i.c.HL())
	case 0x24:
		// INC H
		i.inc_r8(cpu.H)
	case 0x2C:
		// INC L
		i.inc_r8(cpu.L)
	case 0x33:
		// INC SP
		i.c.SP++
	case 0x34:
		// INC (HL)
		i.inc_u16(i.c.HL())
	case 0x3C:
		// INC A
		i.inc_r8(cpu.A)
	default:
		panic("Invalid Opcode for INC")
	}
}

func NewInc(c *cpu.CPU) *inc {
	return &inc{c}
}
