package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type rr struct {
	c *cpu.CPU
}

// RRCA: Rotate A and use higher order bit as Carry
// RRA: Rotate A and use carry flag

func _rr(c *cpu.CPU, reg *uint8) {
	carry := c.CarryVal()

	c.SET_NEG(false)
	c.SET_HALF_CARRY(false)
	c.SET_CARRY(*reg&0x01 == 0x01) // set carry flag if the least significant bit is 1

	*reg >>= 1
	*reg |= carry << 7

	c.SET_ZERO(*reg == 0x00)
}

func _rrc(c *cpu.CPU, reg *uint8) {
	c.SET_NEG(false)
	c.SET_HALF_CARRY(false)

	var lb uint8 = *reg & 0x01
	c.SET_CARRY(lb == 0x01)
	*reg >>= 1
	*reg |= lb << 7

	c.SET_ZERO(*reg == 0x0)
}

func (r *rr) _rra() {
	A := r.c.GetRegister(cpu.A)
	_rr(r.c, A)

	r.c.SET_ZERO(false)
}

func (r *rr) _rrca() {
	A := r.c.GetRegister(cpu.A)
	_rrc(r.c, A)
	r.c.SET_ZERO(false)
}

func (r *rr) Exec(op byte) {
	switch op {
	case 0x0F:
		r._rrca()
	case 0x1F:
		r._rra()
	}

}

func NewRR(c *cpu.CPU) *rr {
	return &rr{
		c,
	}
}
