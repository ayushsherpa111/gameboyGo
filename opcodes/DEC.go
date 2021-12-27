package opcodes

import (
	"fmt"

	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type dec struct {
	c *cpu.CPU
}

func (d *dec) _dec(val *uint8) {
	fmt.Printf("Before DEC 0x%X", *val)
	d.c.SET_NEG(true)
	d.c.SET_HALF_CARRY((*val & 0x0F) == 0)

	*(val)--

	d.c.SET_ZERO(*val == 0x00)
	fmt.Printf("After DEC 0x%X", *val)
}

// dec_r8 decreases 1 from Register reg and sets flags
func (d *dec) dec_r8(reg uint8) {
	d._dec(d.c.GetRegister(reg))
}

func (d *dec) dec_u16(addr uint16) {
	d._dec(d.c.GetMem(addr))
}

func (d *dec) dec_SP() {
	d.c.SP--
}

func (d *dec) dec_r16(r1, r2 uint8, val uint16) {
	val--
	d.c.SetRegister(r1, uint8(val>>8))
	d.c.SetRegister(r2, uint8(val))
}

func (d *dec) Exec(op byte) {
	switch op {
	case 0x05:
		// DEC B
		d.dec_r8(cpu.B)
	case 0x0B:
		// DEC BC
		d.dec_r16(cpu.B, cpu.C, d.c.BC())
	case 0x0D:
		// DEC C
		d.dec_r8(cpu.C)
	case 0x15:
		// DEC D
		d.dec_r8(cpu.D)
	case 0x1B:
		// DEC DE
		d.dec_r16(cpu.D, cpu.E, d.c.DE())
	case 0x1D:
		// DEC E
		d.dec_r8(cpu.E)
	case 0x25:
		// DEC H
		d.dec_r8(cpu.H)
	case 0x2B:
		// DEC HL
		d.dec_r16(cpu.H, cpu.L, d.c.HL())
	case 0x2D:
		// DEC L
		d.dec_r8(cpu.L)
	case 0x35:
		// DEC (HL)
		d.dec_u16(d.c.HL())
	case 0x3B:
		// DEC SP
		d.dec_SP()
	case 0x3D:
		// DEC A
		d.dec_r8(cpu.A)
	default:
		panic("Invalid opcode for DEC")
	}

}

func NewDEC(c *cpu.CPU) *dec {
	return &dec{c}
}
