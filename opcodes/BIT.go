package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type bit struct {
	c      *cpu.CPU
	regMap map[byte]uint8
}

func (b *bit) _bit(mem *uint8, bitNum uint8) {
	b.c.SET_HALF_CARRY(true)
	b.c.SET_NEG(false)

	bitNum = 0x1 << bitNum
	b.c.SET_ZERO(*mem&bitNum == 0)
}

func (b *bit) Exec(op byte) {
	var tar *uint8
	var bitPos uint8

	if op&0xF0 == 0x40 {
		if v := op & 0x0F; v >= 0x00 && v <= 0x07 {
			bitPos = 0
		} else {
			bitPos = 1
		}
	} else if op&0xF0 == 0x50 {
		if v := op & 0x0F; v >= 0x00 && v <= 0x07 {
			bitPos = 2
		} else {
			bitPos = 3
		}
	} else if op&0xF0 == 0x60 {
		if v := op & 0x0F; v >= 0x00 && v <= 0x07 {
			bitPos = 4
		} else {
			bitPos = 5
		}
	} else {
		if v := op & 0x0F; v >= 0x00 && v <= 0x07 {
			bitPos = 6
		} else {
			bitPos = 7
		}
	}

	if v := op & 0x0F; v != 0x06 && v != 0x0E {
		tar = b.c.GetRegister(b.regMap[v])
	} else {
		tar = b.c.GetMem(b.c.HL())
	}

	b._bit(tar, bitPos)
}

func NewBIT(c *cpu.CPU) *bit {
	return &bit{
		c: c,
		regMap: map[byte]uint8{
			0x0: cpu.B,
			0x1: cpu.C,
			0x2: cpu.D,
			0x3: cpu.E,
			0x4: cpu.H,
			0x5: cpu.L,
			0x7: cpu.A,
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
