package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
	"github.com/ayushsherpa111/gameboyEMU/interfaces"
)

type cb struct {
	c      *cpu.CPU
	subMap map[byte]interfaces.Instruction
}

func (c *cb) Exec(op byte) {
	nextOP, _ := c.c.Fetch()
	// fmt.Printf("PC: 0x%x SUB CB: 0x%02x\n", c.c.PC, nextOP)
	c.subMap[nextOP].Exec(nextOP)
}

func NewCB(c *cpu.CPU) *cb {
	RL := NewCBRl(c)
	RR := NewCBRR(c)
	SL := NewSL(c)
	SR := NewSR(c)
	SWAP := NewSwap(c)
	BIT := NewBIT(c)
	RES := NewRes(c)
	SET := NewSET(c)
	subMap := map[byte]interfaces.Instruction{}

	var i uint8
	for i = 0; i <= 0xF; i++ {
		if i <= 0x7 {
			subMap[i] = RL
			subMap[i|0x10] = RL
			subMap[i|0x20] = SL
			subMap[i|0x30] = SWAP
		} else {
			subMap[i] = RR
			subMap[i|0x10] = RR

			subMap[i|0x20] = SR
			subMap[i|0x30] = SR
		}

		subMap[i|0x40] = BIT
		subMap[i|0x50] = BIT
		subMap[i|0x60] = BIT
		subMap[i|0x70] = BIT

		subMap[i|0x80] = RES
		subMap[i|0x90] = RES
		subMap[i|0xA0] = RES
		subMap[i|0xB0] = RES

		subMap[i|0xC0] = SET
		subMap[i|0xD0] = SET
		subMap[i|0xE0] = SET
		subMap[i|0xF0] = SET
	}

	return &cb{
		c:      c,
		subMap: subMap,
	}
}
