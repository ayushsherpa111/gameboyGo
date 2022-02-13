package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type cbrr struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (r *cbrr) Exec(op byte) {
	var reg *uint8
	if op&0x0F != 0x0E {
		reg = r.c.GetRegister(r.regMap[op&0x0F])
	} else {
		reg = r.c.GetMem(r.c.HL())
	}

	if op&0xF0 == 0x00 {
		_rrc(r.c, reg)
	} else if op&0xF0 == 0x10 {
		_rr(r.c, reg)
	}

}

func NewCBRR(c *cpu.CPU) *cbrr {
	return &cbrr{
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
