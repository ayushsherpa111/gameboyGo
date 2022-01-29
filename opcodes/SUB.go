package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type sub struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (s *sub) _SUB(val uint8, carry uint8) {
	A := s.c.GetRegister(cpu.A)

	s.c.SET_NEG(true)
	s.c.SET_HALF_CARRY(*A&0x0F < val&0x0F+carry)
	s.c.SET_HALF_CARRY(*A < val&0x0F+carry)

	*A -= (val + carry)
	s.c.SET_ZERO(*A == 0x0)
}

func (s *sub) sub_r8_u8(val uint8) {
	var carry uint8 = s.c.CarryVal()

	s._SUB(val, carry)
}

func (s *sub) Exec(op byte) {
	if v, ok := s.regMap[op&0x0F]; ok {
		s.sub_r8_u8(*s.c.GetRegister(v))
	} else {
		switch op {
		case 0x96:
			// SUB A, (HL)
			s.sub_r8_u8(*s.c.GetMem(s.c.HL()))
		default:
			panic("Failed to decode opcode in Sub")
		}
	}
}

func NewSub(s *cpu.CPU) *sub {
	return &sub{
		c: s,
		regMap: map[byte]uint8{
			0: cpu.B,
			1: cpu.C,
			2: cpu.D,
			3: cpu.E,
			4: cpu.H,
			5: cpu.L,
			7: cpu.A,
		},
	}
}
