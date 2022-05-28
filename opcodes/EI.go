package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type EI struct {
	c *cpu.CPU
}

func (d *EI) Exec(op byte) {
	d.c.ScheduleEI(1)
}

func NewEI(c *cpu.CPU) *EI {
	return &EI{c}
}
