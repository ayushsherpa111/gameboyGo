package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type cbrl struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (r *cbrl) Exec(op byte) {
	var reg *uint8
	if op&0x0F != 0x06 {
		reg = r.c.GetRegister(r.regMap[op&0x0F])
	} else {
		reg = r.c.GetMem(r.c.HL())
	}

	if op&0xF0 == 0x00 {
		_rlc(r.c, reg)
	} else if op&0xF0 == 0x10 {
		_rl(r.c, reg)
	}
}

func NewCBRl(c *cpu.CPU) *cbrl {
	return &cbrl{
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
