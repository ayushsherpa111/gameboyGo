package opcodes

import "github.com/ayushsherpa111/gameboyEMU/cpu"

type pop struct {
	c *cpu.CPU
}

func (p *pop) pop_r8(r1, r2 *uint8) {
	val := p.c.FetchSP()
	*r2 = uint8(val)
	*r2 = uint8(val >> 8)
}

func (p *pop) Exec(op byte) {
	B := p.c.GetRegister(cpu.B)
	C := p.c.GetRegister(cpu.C)

	D := p.c.GetRegister(cpu.D)
	E := p.c.GetRegister(cpu.E)

	H := p.c.GetRegister(cpu.H)
	L := p.c.GetRegister(cpu.L)

	A := p.c.GetRegister(cpu.A)
	F := p.c.GetRegister(cpu.F)

	switch op {
	case 0xC1:
		// POP BC
		p.pop_r8(B, C)
	case 0xD1:
		// POP DE
		p.pop_r8(D, E)
	case 0xE1:
		// POP HL
		p.pop_r8(H, L)
	case 0xF1:
		// POP AF
		p.pop_r8(A, F)
	}
}

func NewPOP(c *cpu.CPU) *pop {
	return &pop{c}
}
