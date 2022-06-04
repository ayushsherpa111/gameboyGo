package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type halt struct {
	c *cpu.CPU
}

func (h *halt) Exec(opcode byte) {
	h.c.Halted = true
}

func NewHalt(c *cpu.CPU) *halt {
	return &halt{c}
}
