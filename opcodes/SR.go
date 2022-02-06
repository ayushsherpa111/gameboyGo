package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type sr struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (s *sr) _sra(reg *uint8) {
	s.c.SET_NEG(false)
	s.c.SET_HALF_CARRY(false)
	s.c.SET_CARRY(*reg&0x01 == 0x01)

	hb := *reg & 0x80
	*reg >>= 1
	*reg |= hb

	s.c.SET_ZERO(*reg == 0)
}

func (s *sr) _srl(reg *uint8) {
	s.c.SET_CARRY(*reg&0x01 == 0x01)

	*reg >>= 1

	s.c.SET_NEG(false)
	s.c.SET_HALF_CARRY(false)
	s.c.SET_ZERO(*reg == 0)
}

func (s *sr) Exec(op byte) {
	var reg *uint8
	if op&0x0F != 0x06 {
		reg = s.c.GetRegister(s.regMap[op&0x0F])
	} else {
		reg = s.c.GetMem(s.c.HL())
	}
	if op&0xF0 == 0x20 {
		s._sra(reg)
	} else {
		s._srl(reg)
	}
}

func NewSR(c *cpu.CPU) *sl {
	return &sl{
		c: c,
		regMap: map[byte]uint8{
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
