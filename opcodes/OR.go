package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type or struct {
	c      *cpu.CPU
	keyMap map[byte]uint8
}

func (o *or) or_r8_u8(val uint8) {
	o.c.SET_NEG(false)
	o.c.SET_HALF_CARRY(false)
	o.c.SET_CARRY(false)

	A := o.c.GetRegister(cpu.A)

	*A |= val
	o.c.SET_ZERO(*A == 0)
}

func (o *or) Exec(op byte) {
	if v, ok := o.keyMap[op&0x0F]; ok {
		o.or_r8_u8(*o.c.GetRegister(v))
	} else {
		arg, err := o.c.Fetch()
		if err != nil {
			return
		}
		switch op {
		case 0xB6:
			HL := o.c.HL()
			o.or_r8_u8(*o.c.GetMem(HL))
		case 0xF6:
			o.or_r8_u8(arg)
		default:
			panic("Failed to decode opcode for OR")
		}
	}
}

func NewOR(c *cpu.CPU) *or {
	return &or{
		c: c,
		keyMap: map[byte]uint8{
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
