package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type sbc struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (s *sbc) SUB_r8_u8(val uint8) {
	s.c.SET_NEG(true)

	var carry uint8 = s.c.CarryVal()
	A := s.c.GetRegister(cpu.A)

	s.c.SET_HALF_CARRY(*A&0x0F < (val&0x0F)+carry)
	s.c.SET_CARRY(uint16(*A) < uint16(carry)+uint16(val))

	*A -= (carry + val)

	s.c.SET_ZERO(*A == 0)
}

func (s *sbc) Exec(op byte) {
	if v, ok := s.regMap[op&0x0F]; ok {
		s.SUB_r8_u8(*s.c.GetRegister(v))
	} else {
		switch op {
		case 0x9E:
			// SBC A,(HL)
			HL := s.c.HL()
			s.SUB_r8_u8(*s.c.GetMem(HL))
		case 0xDE:
			//  SBC A, u8
			arg, err := s.c.Fetch()
			if err != nil {
				return
			}
			s.SUB_r8_u8(arg)
		default:
			panic("Invalid opcode for sbc")
		}
	}
}

func NewSBC(c *cpu.CPU) *sbc {
	return &sbc{
		c: c,
		regMap: map[byte]uint8{
			0x8: cpu.B,
			0x9: cpu.C,
			0xa: cpu.D,
			0xb: cpu.E,
			0xc: cpu.H,
			0xd: cpu.L,
			0xf: cpu.A,
		},
	}
}
