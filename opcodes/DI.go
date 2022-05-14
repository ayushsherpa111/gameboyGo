package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type DI struct {
	c *cpu.CPU
}

func (d *DI) Exec(op byte) {
	d.c.IME = false
}

func NewDI(c *cpu.CPU) *DI {
	return &DI{c}
}
