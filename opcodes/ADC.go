package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type ADC struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (a *ADC) add_r8_u8(u8 uint8) {
	var carry uint8 = a.c.CarryVal()
	A := a.c.GetRegister(cpu.A)

	a.c.SET_NEG(false)
	a.c.SET_HALF_CARRY((*A&0x0F)+(u8&0x0F)+carry > 0x0F)
	a.c.SET_CARRY(uint16(*A)+uint16(u8&0x0F)+uint16(carry) > 0xFF)

	*A += u8 + carry
	a.c.SET_ZERO(*A == 0x0)
}

func (a *ADC) Exec(op byte) {
	if v, ok := a.regMap[op&0x0F]; ok {
		a.add_r8_u8(*a.c.GetRegister(v))
	} else {
		arg, e := a.c.Fetch()
		if e != nil {
			return
		}
		switch op {
		case 0x8E:
			// ADC A, (HL)
			HL := a.c.HL()
			a.add_r8_u8(*a.c.GetMem(HL))
		case 0xCE:
			// ADC A, u8
			a.add_r8_u8(arg)
		default:
			panic("Failed to decode opcode for ADC")
		}
	}
}

func NewADC(c *cpu.CPU) *ADC {
	return &ADC{
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
