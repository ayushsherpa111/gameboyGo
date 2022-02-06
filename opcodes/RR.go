package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type rr struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

// RRCA: Rotate A and use higher order bit as Carry
// RRA: Rotate A and use carry flag

func (r *rr) _rr(reg *uint8, carry uint8) {
	r.c.SET_NEG(false)
	r.c.SET_HALF_CARRY(false)
	r.c.SET_CARRY(*reg&0x01 == 0x01) // set carry flag if the least significant bit is 1

	*reg >>= 1
	*reg |= carry << 7

	r.c.SET_ZERO(*reg == 0x00)
}

func (r *rr) _rra(carry uint8) {
	A := r.c.GetRegister(cpu.A)
	r._rr(A, carry)

	r.c.SET_ZERO(false)
}

func (r *rr) _rrc(reg *uint8) {
	r.c.SET_NEG(false)
	r.c.SET_HALF_CARRY(false)

	var lb uint8 = *reg & 0x01
	r.c.SET_CARRY(lb == 0x01)
	*reg >>= 1
	*reg |= lb << 7

	r.c.SET_ZERO(*reg == 0x0)
}

func (r *rr) _rrca() {
	A := r.c.GetRegister(cpu.A)
	r._rrc(A)
	r.c.SET_ZERO(false)
}

func (r *rr) Exec(op byte) {
	carry := r.c.CarryVal()
	var reg *uint8
	if op&0x0F != 0x0E {
		reg = r.c.GetRegister(r.regMap[op&0x0F])
	} else {
		reg = r.c.GetMem(r.c.HL())
	}

	if op == 0x0f {
		r._rrca()
	} else if op == 0x1f {
		r._rra(carry)
	} else {
		if op&0xF0 == 0x00 {
			r._rrc(reg)
		} else {
			r._rr(reg, carry)
		}
	}

}

func NewRR(c *cpu.CPU) *rr {
	return &rr{
		c,
		map[byte]uint8{
			0x8: cpu.B,
			0x9: cpu.C,
			0xA: cpu.D,
			0xB: cpu.E,
			0xC: cpu.H,
			0xD: cpu.L,
			0xF: cpu.A,
		},
	}
}
