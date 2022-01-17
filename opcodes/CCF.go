package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type ccf struct {
	c *cpu.CPU
}

func (c *ccf) Exec(op byte) {
	c.c.SET_NEG(false)
	c.c.SET_HALF_CARRY(false)
	c.c.SET_CARRY(!c.c.CarryFlag())
}

func NewCCF(c *cpu.CPU) *ccf {
	return &ccf{c}
}
