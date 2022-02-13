package opcodes

import (
	"fmt"

	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type sub struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (s *sub) _SUB(val uint8) {
	A := s.c.GetRegister(cpu.A)

	s.c.SET_NEG(true)
	s.c.SET_HALF_CARRY(*A&0x0F < val&0x0F)
	s.c.SET_CARRY(*A < val)

	*A -= val
	s.c.SET_ZERO(*A == 0x0)
}

func (s *sub) sub_r8_u8(val uint8) {
	// var carry uint8 = s.c.CarryVal()

	s._SUB(val)
}

func (s *sub) Exec(op byte) {
	if v, ok := s.regMap[op&0x0F]; ok {
		s.sub_r8_u8(*s.c.GetRegister(v))
	} else {
		switch op {
		case 0x96:
			// SUB A, (HL)
			s.sub_r8_u8(*s.c.GetMem(s.c.HL()))
		case 0xD6:
			// SUB A. u8
			arg, _ := s.c.Fetch()
			s._SUB(arg)
		default:
			panic(fmt.Sprintf("Failed to decode opcode in Sub 0x%02x", op))
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
