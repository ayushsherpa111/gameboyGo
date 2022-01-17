package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type and struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (a *and) AND_r8_u8(val uint8) {
	a.c.SET_NEG(false)
	a.c.SET_CARRY(false)
	a.c.SET_HALF_CARRY(true)

	A := a.c.GetRegister(cpu.A)
	*A &= val

	a.c.SET_ZERO(*A == 0)
}

func (a *and) Exec(op byte) {
	if v, ok := a.regMap[op&0x0F]; ok {
		a.AND_r8_u8(*a.c.GetRegister(v))
	} else {
		arg, err := a.c.Fetch()
		if err != nil {
			return
		}
		switch op {
		case 0xA6:
			// AND A, (HL)
			HL := a.c.HL()
			a.AND_r8_u8(*a.c.GetMem(HL))
		case 0xE6:
			// AND A, u8
			a.AND_r8_u8(arg)
		default:
			panic("Invalid opcode for AND")
		}
	}
}

func NewAND(c *cpu.CPU) *and {
	return &and{c: c,
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
