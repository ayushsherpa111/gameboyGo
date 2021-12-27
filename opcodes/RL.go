package opcodes

import (
	"fmt"

	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type rl struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

// RLCA: Rotate A and use higher order bit as Carry
// RLA: Rotate A and use carry flag

// TODO: fix RL. shifting to the right does nothing

// Before RL 0xc000100040 Carry: 0x01
// After RL 0xc000100040
func (r *rl) _rl(reg *uint8, carry uint8) {
	fmt.Printf("Before RL 0x%02x Carry: 0x%02x\n", reg, carry)
	r.c.SET_CARRY((*reg & 0x80) == 0x80)
	r.c.SET_NEG(false)
	r.c.SET_HALF_CARRY(false)

	*reg <<= 1
	*reg |= carry

	r.c.SET_ZERO(*reg == 0x0)
	fmt.Printf("After RL 0x%x\n", reg)
}

func (r *rl) _rlc(reg *uint8) {
	OF := *reg & 0x80
	r.c.SET_CARRY(OF > 0x0)
	*reg <<= 1
	*reg |= OF

	r.c.SET_ZERO(*reg == 0x0)
	r.c.SET_NEG(false)
	r.c.SET_HALF_CARRY(false)
}

func (r *rl) _rla(carry uint8) {
	A := r.c.GetRegister(cpu.A)
	r._rl(A, carry)

	r.c.SET_ZERO(false)
}

func (r *rl) _rlca() {
	A := r.c.GetRegister(cpu.A)
	r._rlc(A)

	r.c.SET_ZERO(false)
}

func (r *rl) Exec(op byte) {
	carry := r.c.CarryVal()
	var reg *uint8
	if op&0x0F != 0x06 {
		reg = r.c.GetRegister(r.regMap[op&0x0F])
	} else {
		reg = r.c.GetMem(r.c.HL())
	}

	if op&0xF0 == 0x00 {
		if op == 0x07 {
			r._rlca()
		} else {
			r._rlc(reg)
		}
	} else {
		if op == 0x17 {
			r._rla(carry)
		} else {
			r._rl(reg, carry)
		}
	}
}

func NewRl(c *cpu.CPU) *rl {
	return &rl{
		c,
		map[byte]uint8{
			0x0: cpu.B,
			0x1: cpu.C,
			0x2: cpu.D,
			0x3: cpu.E,
			0x4: cpu.H,
			0x5: cpu.L,
			0x7: cpu.A,
		},
	}
}
