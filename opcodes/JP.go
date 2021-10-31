package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type jp struct {
	c *cpu.CPU
}

func (j *jp) jmp_if(mem uint16, cond bool) {
	if cond {
		j.c.PC = mem
	}
}

func (j *jp) Exec(op byte) {
	switch op {
	case 0xC2:
		// JP NZ, u16
		j.jmp_if(j.c.Fetch16(), !j.c.ZeroFlag())
	case 0xC3:
		// JP u16
		j.jmp_if(j.c.Fetch16(), true)
	case 0xD2:
		// JP NC,u16
		j.jmp_if(j.c.Fetch16(), !j.c.CarryFlag())
	case 0xCA:
		// JP Z,u16
		j.jmp_if(j.c.Fetch16(), j.c.ZeroFlag())
	case 0xDA:
		// JP C,u16
		j.jmp_if(j.c.Fetch16(), j.c.CarryFlag())
	case 0xE9:
		// JP HL
		j.jmp_if(j.c.HL(), true)
	}
}

func NewJP(c *cpu.CPU) *jp {
	return &jp{c}
}
