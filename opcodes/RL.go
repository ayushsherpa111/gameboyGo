package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type rl struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

// RLCA: Rotate A and use higher order bit as Carry
// RLA: Rotate A and use carry flag

func (r *rl) _rl(reg *uint8, carry uint8) {

	r.c.SET_NEG(false)
	r.c.SET_HALF_CARRY(false)
	r.c.SET_CARRY(carry != 0)

	*reg <<= 1
	*reg |= carry

	r.c.SET_ZERO(*reg == 0x00)
}

func (r *rl) _rlc(reg *uint8) {
	r._rl(reg, *reg>>7)
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
		r._rlc(reg)
	} else {
		r._rl(reg, carry)
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
