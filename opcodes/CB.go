package opcodes

import (
	"github.com/ayushsherpa111/gameboyEMU/cpu"
	"github.com/ayushsherpa111/gameboyEMU/instructions"
)

type cb struct {
	c      *cpu.CPU
	subMap map[byte]instructions.Instruction
}

func (c *cb) Exec(op byte) {
	nextOP := c.c.Fetch()
	c.subMap[nextOP].Exec(nextOP)
}

func NewCB(c *cpu.CPU) *cb {
	RL := NewRl(c)
	RR := NewRR(c)
	SL := NewSL(c)
	SR := NewSR(c)
	SWAP := NewSwap(c)
	BIT := NewBIT(c)
	subMap := map[byte]instructions.Instruction{}

	var i uint8
	for i = 0; i < 0xF; i++ {
		if i < 0x8 {
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
	}

	return &cb{
		c:      c,
		subMap: subMap,
	}
}
