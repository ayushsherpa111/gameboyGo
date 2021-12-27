package opcodes

import (
	"fmt"

	"github.com/ayushsherpa111/gameboyEMU/cpu"
)

type push struct {
	c *cpu.CPU
}

func (p *push) push_r16(val uint16) {
	p.c.PushSP(val)
	fmt.Printf("Pushing 0x%02x\n", val)
}

func (p *push) Exec(op byte) {
	switch op {
	case 0xC5:
		// PUSH BC
		p.push_r16(p.c.BC())
	case 0xD5:
		// PUSH BC
		p.push_r16(p.c.DE())
	case 0xE5:
		// PUSH BC
		p.push_r16(p.c.HL())
	case 0xF5:
		// PUSH BC
		p.push_r16(p.c.AF())
	}
}

func NewPush(c *cpu.CPU) *push {
	return &push{c}
}
