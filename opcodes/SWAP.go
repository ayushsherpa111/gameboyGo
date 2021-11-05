package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type swap struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (s *swap) _swp(val *uint8) {
	s.c.SET_NEG(false)
	s.c.SET_HALF_CARRY(false)
	s.c.SET_CARRY(false)

	ln := *val & 0x0F
	*val <<= 4
	*val |= ln

	s.c.SET_ZERO(*val == 0x0)
}

func (s *swap) Exec(op byte) {
	var tar *uint8

	if v := op & 0x0F; v != 0x06 {
		reg := s.c.GetRegister(s.regMap[v])
		tar = reg
	} else {
		tar = s.c.GetMem(s.c.HL())
	}

	s._swp(tar)
}

func NewSwap(c *cpu.CPU) *swap {
	return &swap{
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
