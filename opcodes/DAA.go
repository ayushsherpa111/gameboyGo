package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type daa struct {
	c *cpu.CPU
}

func (d *daa) _daa() {
	d.c.SET_HALF_CARRY(false)

	A := d.c.GetRegister(cpu.A)
	if !d.c.NegativeFlag() {
		if d.c.CarryFlag() || ((*A & 0xF0) > 0x90) {
			*A += 0x60 // skip A-F
			d.c.SET_CARRY(true)
		}
		if d.c.HalfCarryFlag() || (*A&0x0F) > 0x09 {
			*A += 0x06 // skip A-F
		}
	} else {
		if d.c.CarryFlag() {
			*A -= 0x60 // skip A-F
			d.c.SET_CARRY(true)
		}
		if d.c.HalfCarryFlag() {
			*A -= 0x06 // skip A-F
		}
	}

	d.c.SET_ZERO(*A == 0x0)
}

func (d *daa) Exec(op byte) {
	d._daa()
}

func NewDAA(c *cpu.CPU) *daa {
	return &daa{c}
}
