package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type jr struct {
	c *cpu.CPU
}

func (j *jr) jr_i8(jmp int8) {
	j.c.PC = uint16(int16(j.c.PC) + int16(jmp))
}

func (j *jr) jumpIf(cond bool, size int8) {
	if cond {
		j.c.PC = uint16(int16(j.c.PC) + int16(size))
	}
}

func (j *jr) Exec(op byte) {
	next, e := j.c.Fetch()

	if e != nil {
		return
	}

	switch op {
	case 0x18:
		// JR i8
		j.jr_i8(int8(next))
	case 0x28:
		// JR Z, i8
		j.jumpIf(j.c.ZeroFlag(), int8(next))
	case 0x38:
		// JR C, i8
		j.jumpIf(j.c.CarryFlag(), int8(next))
	case 0x20:
		// JR NZ, i8
		j.jumpIf(!j.c.ZeroFlag(), int8(next))
	case 0x30:
		// JR NC, i8
		j.jumpIf(!j.c.CarryFlag(), int8(next))
	default:
		panic("Invalid Opcode for JR")
	}
}

func NewJR(c *cpu.CPU) *jr {
	return &jr{c}
}
