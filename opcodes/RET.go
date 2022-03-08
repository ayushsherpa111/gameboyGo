package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type ret struct {
	c *cpu.CPU
}

func (r *ret) POP_r16() {
	r.c.PC = r.c.FetchSP()
}

func (r *ret) RET_COND(cond bool) {
	if cond {
		r.POP_r16()
	}
}

func (r *ret) Exec(op byte) {
	switch op {
	case 0xC0:
		// RET NZ
		r.RET_COND(!r.c.ZeroFlag())
	case 0xD0:
		// RET NC
		r.RET_COND(!r.c.CarryFlag())
	case 0xC8:
		// RET Z
		r.RET_COND(r.c.ZeroFlag())
	case 0xD8:
		// RET C
		r.RET_COND(r.c.CarryFlag())
	case 0xC9:
		// RET
		r.RET_COND(true)
	case 0xD9:
		r.RET_COND(true)
		r.c.NewIMEConf = cpu.NewImePayload(r.c.PC+1, true)
	default:
		panic("Invalid opcode for RET")
	}
}

func NewRet(c *cpu.CPU) *ret {
	return &ret{c}
}
