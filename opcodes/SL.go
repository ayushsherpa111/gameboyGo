package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type sl struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (s *sl) _sla(reg *uint8) {
	s.c.SET_NEG(false)
	s.c.SET_HALF_CARRY(false)
	s.c.SET_CARRY(*reg&0x80 != 0x00)
	*reg <<= 1
	s.c.SET_ZERO(*reg == 0)
}

func (s *sl) Exec(op byte) {
	var reg *uint8
	if op&0x0F != 0x06 {
		reg = s.c.GetRegister(s.regMap[op&0x0F])
	} else {
		reg = s.c.GetMem(s.c.HL())
	}
	s._sla(reg)
}

func NewSL(c *cpu.CPU) *sl {
	return &sl{
		c: c,
		regMap: map[byte]uint8{
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
