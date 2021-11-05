package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type daa struct {
	c *cpu.CPU
}

func (d *daa) Exec(op byte) {
	if d.c.NegativeFlag() {

	}
}

func NewDAA(c *cpu.CPU) *daa {
	return &daa{c}
}
