package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type call struct {
	c *cpu.CPU
}

func (c *call) callCond(cond bool, addr uint16) {
	if !cond {
		return
	}
	c.c.PushSP(c.c.PC)
	c.c.PC = addr
}

func (c *call) Exec(op byte) {
	ZF := c.c.ZeroFlag()
	CF := c.c.CarryFlag()
	addr := c.c.Fetch16()
	switch op {
	case 0xC4:
		// CALL NZ, u16
		c.callCond(!ZF, addr)
	case 0xCC:
		// CALL Z, u16
		c.callCond(ZF, addr)
	case 0xD4:
		// CALL NC, u16
		c.callCond(!CF, addr)
	case 0xDC:
		// CALL C, u16
		c.callCond(CF, addr)
	case 0xCD:
		// CALL u16
		c.callCond(true, addr)
	}
}

func NewCall(c *cpu.CPU) *call {
	return &call{c}
}
