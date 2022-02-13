package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type cpl struct {
	c *cpu.CPU
}

func (c *cpl) Exec(op byte) {
	A := c.c.GetRegister(cpu.A)
	*A = ^(*A)

	c.c.SET_HALF_CARRY(true)
	c.c.SET_NEG(true)
}

func NewCPL(c *cpu.CPU) *cpl {
	return &cpl{c}
}
