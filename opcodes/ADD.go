package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type add struct {
	c *cpu.CPU
}

func (a *add) _add(tar *uint8, src uint8) {
	a.c.SET_NEG(false)
	a.c.SET_HALF_CARRY((*tar&0x0F)+(src&0x0F) > 0x0F)
	a.c.SET_CARRY(uint16(*tar)+uint16(src) > 0xFF)

	*tar += src
	a.c.SET_ZERO(*tar == 0x00)
}

func (a *add) add_r8_r8(r1, r2 uint8) {
	a._add(a.c.GetRegister(r1), *a.c.GetRegister(r2))
}

func (a *add) add_r8_u8(r1, v2 uint8) {
	a._add(a.c.GetRegister(r1), v2)
}

func (a *add) add_u16_u16(r1, r2 uint8, src uint16) {
	a.c.SET_NEG(false)

	val1, val2 := a.c.GetRegister(r1), a.c.GetRegister(r2)
	r16 := uint16(*val1)<<8 | uint16(*val2)
	a.c.SET_HALF_CARRY((r16&0x0FFF)+(src&0x0FFF) > 0x0FFF)
	a.c.SET_CARRY(uint32(r16)+uint32(src) > 0xFFFF)

	sum := r16 + src

	a.c.SetRegister(r1, uint8(sum>>8))
	a.c.SetRegister(r2, uint8(sum))
}

func (a *add) add_SP_i8() {
	u8, _ := a.c.Fetch()
	a.c.SET_ZERO(false)
	a.c.SET_NEG(false)

	a.c.SET_HALF_CARRY(a.c.SP&0x0F+uint16(u8)&0x0F > 0x0F)
	a.c.SET_CARRY(a.c.SP&0xFF+uint16(u8) > 0xFF)

	a.c.SP = uint16(int16(a.c.SP) + int16(int8(u8)))
}

func (a *add) Exec(op byte) {
	switch op {
	case 0x09:
		// ADD HL, BC
		a.add_u16_u16(cpu.H, cpu.L, a.c.BC())
	case 0x19:
		// ADD HL, DE
		a.add_u16_u16(cpu.H, cpu.L, a.c.DE())
	case 0x29:
		// ADD HL, HL
		a.add_u16_u16(cpu.H, cpu.L, a.c.HL())
	case 0x39:
		// ADD HL, SP
		a.add_u16_u16(cpu.H, cpu.L, a.c.SP)
	case 0x80:
		// ADD A, B
		a.add_r8_r8(cpu.A, cpu.B)
	case 0x81:
		// ADD A, C
		a.add_r8_r8(cpu.A, cpu.C)
	case 0x82:
		// ADD A, D
		a.add_r8_r8(cpu.A, cpu.D)
	case 0x83:
		// ADD A, E
		a.add_r8_r8(cpu.A, cpu.E)
	case 0x84:
		// ADD A, H
		a.add_r8_r8(cpu.A, cpu.H)
	case 0x85:
		// ADD A, L
		a.add_r8_r8(cpu.A, cpu.L)
	case 0x86:
		// ADD A, (HL)
		a.add_r8_u8(cpu.A, *a.c.GetMem(a.c.HL()))
	case 0x87:
		// ADD A, A
		a.add_r8_r8(cpu.A, cpu.A)
	case 0xC6:
		// ADD A, u8
		arg, err := a.c.Fetch()
		if err != nil {
			return
		}
		a.add_r8_u8(cpu.A, arg)
	case 0xE8:
		// ADD SP, i8
		a.add_SP_i8()
	default:
		panic("not implemented")
	}
}
func NewADD(c *cpu.CPU) *add {
	return &add{
		c,
	}
}
