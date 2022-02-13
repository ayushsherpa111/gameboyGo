package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type rl struct {
	c *cpu.CPU
}

// RLCA: Rotate A and use higher order bit as Carry
// RLA: Rotate A and use carry flag

func _rl(c *cpu.CPU, reg *uint8) {
	cv := c.CarryVal()

	var OF uint8 = (*reg & 0x80) >> 7
	c.SET_CARRY(OF > 0x00)
	c.SET_NEG(false)
	c.SET_HALF_CARRY(false)

	*reg <<= 1
	*reg |= cv

	c.SET_ZERO(*reg == 0x0)
}

func _rlc(c *cpu.CPU, reg *uint8) {
	c.SET_NEG(false)
	c.SET_HALF_CARRY(false)
	OF := (*reg & 0x80) >> 7
	c.SET_CARRY(OF > 0x0)

	*reg <<= 1
	*reg |= OF

	c.SET_ZERO(*reg == 0x0)
}

func (r *rl) _rla() {
	A := r.c.GetRegister(cpu.A)
	_rl(r.c, A)

	r.c.SET_ZERO(false)
}

func (r *rl) _rlca() {
	A := r.c.GetRegister(cpu.A)
	_rlc(r.c, A)

	r.c.SET_ZERO(false)
}

func (r *rl) Exec(op byte) {
	switch op {
	case 0x07:
		r._rlca()
	case 0x17:
		r._rla()
	}
}

func NewRl(c *cpu.CPU) *rl {
	return &rl{
		c,
	}
}
