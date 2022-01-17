package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type DI struct {
	c *cpu.CPU
}

func (d *DI) Exec(op byte) {
	d.c.NewIMEConf = cpu.NewImePayload(
		d.c.PC+1,
		false,
	)
}

func NewDI(c *cpu.CPU) *DI {
	return &DI{c}
}
