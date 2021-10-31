package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type cp struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (c *cp) cp_r8_u8(val uint8) {
	A := c.c.GetRegister(cpu.A)

	c.c.SET_NEG(false)
	c.c.SET_ZERO(*A == 0)
	c.c.SET_HALF_CARRY(*A&0x0F < val&0x0F)
	c.c.SET_CARRY(*A < val)
}

func (c *cp) Exec(op byte) {
	if v, ok := c.regMap[op&0x0F]; ok {
		c.cp_r8_u8(*c.c.GetRegister(v))
	} else {
		switch op {
		case 0xB8:
			HL := c.c.HL()
			c.cp_r8_u8(*c.c.GetMem(HL))
		case 0xFE:
			c.cp_r8_u8(c.c.Fetch())
		default:
			panic("Invalid opcode for CP")
		}
	}
}

func NewCP(c *cpu.CPU) *cp {
	return &cp{
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
