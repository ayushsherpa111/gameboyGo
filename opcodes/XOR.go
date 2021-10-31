package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type xor struct {
	c      *cpu.CPU
	keyMap map[byte]uint8
}

func (x *xor) xor_r8_u8(val uint8) {
	x.c.SET_NEG(false)
	x.c.SET_HALF_CARRY(false)
	x.c.SET_CARRY(false)

	A := x.c.GetRegister(cpu.A)
	*A ^= val

	x.c.SET_ZERO(*A == 0)
}

func (x *xor) Exec(op byte) {
	if v, ok := x.keyMap[op&0x0F]; ok {
		x.xor_r8_u8(*x.c.GetRegister(v))
	} else {
		switch op {
		case 0xAE:
			HL := x.c.HL()
			x.xor_r8_u8(*x.c.GetMem(HL))
		case 0xEE:
			x.xor_r8_u8(x.c.Fetch())
		}
	}
}

func NewXOR(c *cpu.CPU) *xor {
	return &xor{
		c: c,
		keyMap: map[byte]uint8{
			0x08: cpu.B,
			0x09: cpu.C,
			0x0A: cpu.D,
			0x0B: cpu.E,
			0x0C: cpu.H,
			0x0D: cpu.L,
			0x0F: cpu.A,
		},
	}
}
